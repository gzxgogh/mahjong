package service

import (
	"fmt"
	"github.com/gzxgogh/ggin/utils"
	"mahjong/model"
	"mahjong/redis"
	"sort"
	"time"
)

// 获取从哪个用户开始抓牌
func GetStartGroup(diceNum, player int) int {
	if diceNum == 1 || diceNum == 5 || diceNum == 9 {
		return player
	} else if diceNum == 2 || diceNum == 6 || diceNum == 10 {
		if player == 4 {
			return 1
		} else {
			return player + 1
		}
	} else if diceNum == 3 || diceNum == 7 || diceNum == 11 {
		if player == 1 || player == 2 {
			return player + 2
		} else if player == 3 {
			return 1
		} else {
			return 2
		}
	} else if diceNum == 4 || diceNum == 8 || diceNum == 12 {
		if player == 1 {
			return 4
		} else if player == 2 {
			return 1
		} else if player == 3 {
			return 2
		} else if player == 4 {
			return 3
		}
	}
	return 0
}

// 获取金
func GetGoldCard(roomNum int) model.Card {
	value := redis.GetValue(fmt.Sprintf(`%d-gold`, roomNum))
	var gold model.Card
	utils.FromJSON(value, &gold)
	return gold
}

// 抓牌并且分牌
func GrabTheCard(roomNum, startGroupNum, startNum int, allCardsArr []model.Card) {
	var newCardsArr, surplusCardArr []model.Card
	var grabTheCardArr []int
	switch startGroupNum {
	case 1:
		grabTheCardArr = []int{1, 2, 3, 4}
	case 2:
		startNum = startNum + (28 * 1)
		grabTheCardArr = []int{2, 3, 4, 1}
	case 3:
		startNum = startNum + (28 * 2)
		grabTheCardArr = []int{3, 4, 1, 2}
	case 4:
		startNum = startNum + (28 * 3)
		grabTheCardArr = []int{4, 1, 2, 3}
	}
	//重新排序，用户从第0个牌开始抓取即可
	newCardsArr = append(newCardsArr, allCardsArr[startNum:]...)
	newCardsArr = append(newCardsArr, allCardsArr[:startNum]...)
	keyPlayerCards := make(map[int][]model.Card)
	curNum := 0

	for i := 0; i < 4; i++ { //每个人都能抓四次牌
		for _, v := range grabTheCardArr { //抓牌的用户顺序
			for j := 0; j < 4; j++ { //一次抓四张牌
				keyPlayerCards[v] = append(keyPlayerCards[v], newCardsArr[curNum])
				curNum++
			}
		}
	}

	//剩余的牌数,最后一张为金,庄家多模第一张门牌,扣减剩余排队
	surplusCardArr = append(surplusCardArr, newCardsArr[curNum:]...)
	length := len(surplusCardArr)
	gold := surplusCardArr[length-1]
	surplusCardArr = append(surplusCardArr[:length-1], surplusCardArr[length:]...)

	for playerNum, arr := range keyPlayerCards {
		kInfo := make(map[string][]int)
		for _, item := range arr {
			if item.String() == gold.String() {
				kInfo[model.CardTypeG] = append(kInfo[model.CardTypeG], 1)
			} else if item.Type != model.CardTypeZ && item.String() != gold.String() {
				kInfo[item.Type] = append(kInfo[item.Type], item.Value)
			} else {
				kInfo[model.CardTypeZ] = append(kInfo[model.CardTypeZ], 1)
			}
		}
		for _, v := range kInfo {
			sort.Ints(v)
		}
		fmt.Println("用户player", playerNum, "的手牌为:", utils.ToJSON(kInfo))
		//存入用户手牌
		redis.SetValue(fmt.Sprintf(`%d-player%d`, roomNum, playerNum), utils.ToJSON(kInfo), 1*time.Hour)
	}
	fmt.Println("当场的金为", gold.String())
	fmt.Println("剩余的牌堆为", utils.ToJSON(surplusCardArr))
	//存入分配玩后各个玩家手里的牌，和场上现有的牌
	redis.SetValue(fmt.Sprintf(`%d-gold`, roomNum), utils.ToJSON(gold), 1*time.Hour)
	redis.SetValue(fmt.Sprintf(`%d-surplusCard`, roomNum), utils.ToJSON(surplusCardArr), 1*time.Hour)
}

// 获取剩余牌堆的牌
func GetSurplusCard(roomNum int) []model.Card {
	value := redis.GetValue(fmt.Sprintf(`%d-surplusCard`, roomNum))
	var surplusCard []model.Card
	utils.FromJSON(value, &surplusCard)
	return surplusCard
}

// 获取用户手牌
func GetPlayerCardInfo(roomNum int, player string) map[string][]int {
	value := redis.GetValue(fmt.Sprintf(`%d-%s`, roomNum, player))
	cardInfo := make(map[string][]int)
	utils.FromJSON(value, &cardInfo)
	return cardInfo
}

// 获取下家用户
func GetNextPlayer(curPlayer string) string {
	player := ""
	switch curPlayer {
	case "player1":
		player = "player2"
	case "player2":
		player = "player3"
	case "player3":
		player = "player4"
	case "player4":
		player = "player1"
	}
	return player
}

// 抢金
func robGold(cardInfo map[string][]int) bool {
	cardInfo[model.CardTypeG] = append(cardInfo[model.CardTypeG], 1)
	goldNum := len(model.CardTypeG)
	pairNum := 0
	isXianJin := false
	pairCard := model.Card{}
	for typ, arr := range cardInfo {
		if typ == model.CardTypeG {
			continue
		}
		//1-9 每张牌的数量
		cardsNum := make([]int, 10)
		for _, card := range arr {
			cardsNum[card]++
		}
		isHu := computeCards(cardsNum, &isXianJin, &pairNum, &goldNum, &pairCard)
		if !isHu {
			return false
		}
	}
	return true
}

// 吃牌
func eatCard(curCard model.Card, cardInfo map[string][]int) (bool, [][]model.Card) {
	var lessTwoCard, lessOneCard, greaterOneCard, greaterTwoCard int
	flag := false
	var cardGroup []model.Card
	switch curCard.Value {
	case 1:
		for _, item := range cardInfo[curCard.Type] {
			fmt.Println("item", item)
			if item == 2 {
				greaterOneCard = item
			} else if item == 3 {
				greaterTwoCard = item
			}
		}
		if greaterOneCard != 0 && greaterTwoCard != 0 {
			flag = true
			cardGroup = append(cardGroup, curCard)
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterOneCard})
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterTwoCard})
		}
	case 2:
		for _, item := range cardInfo[curCard.Type] {
			if item == 1 {
				lessOneCard = item
			} else if item == 3 {
				greaterOneCard = item
			} else if item == 4 {
				greaterTwoCard = item
			}
		}
		if lessOneCard != 0 && greaterOneCard != 0 {
			flag = true
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessOneCard})
			cardGroup = append(cardGroup, curCard)
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterOneCard})
		}
		if greaterOneCard != 0 && greaterTwoCard != 0 {
			flag = true
			cardGroup = append(cardGroup, curCard)
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterOneCard})
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterTwoCard})
		}
	case 8:
		for _, item := range cardInfo[curCard.Type] {
			if item == 7 {
				lessOneCard = item
			} else if item == 9 {
				greaterOneCard = item
			} else if item == 6 {
				lessTwoCard = item
			}
		}
		if lessOneCard != 0 && greaterOneCard != 0 {
			flag = true
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessOneCard})
			cardGroup = append(cardGroup, curCard)
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterOneCard})
		}
		if lessOneCard != 0 && lessTwoCard != 0 {
			flag = true
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessTwoCard})
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessOneCard})
			cardGroup = append(cardGroup, curCard)
		}
	case 9:
		for _, item := range cardInfo[curCard.Type] {
			if item == 7 {
				lessTwoCard = item
			} else if item == 8 {
				lessOneCard = item
			}
		}
		if lessOneCard != 0 && lessTwoCard != 0 {
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessTwoCard})
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessOneCard})
			cardGroup = append(cardGroup, curCard)
		}
	default:
		for _, item := range cardInfo[curCard.Type] {
			if item == curCard.Value-2 {
				lessTwoCard = item
			} else if item == curCard.Value-1 {
				lessOneCard = item
			} else if item == curCard.Value+1 {
				greaterOneCard = item
			} else if item == curCard.Value+2 {
				greaterTwoCard = item
			}
		}
		if lessTwoCard != 0 && lessOneCard != 0 {
			flag = true
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessTwoCard})
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessOneCard})
			cardGroup = append(cardGroup, curCard)
		}
		if lessOneCard != 0 && greaterOneCard != 0 {
			flag = true
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: lessOneCard})
			cardGroup = append(cardGroup, curCard)
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterOneCard})
		}
		if greaterOneCard != 0 && greaterTwoCard != 0 {
			flag = true
			cardGroup = append(cardGroup, curCard)
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterOneCard})
			cardGroup = append(cardGroup, model.Card{Type: curCard.Type, Value: greaterTwoCard})
		}
	}
	total := len(cardGroup) / 3

	var finalArr [][]model.Card
	for i := 0; i < total; i++ {
		var newArr []model.Card
		for j := 0; j < 3; j++ {
			if len(newArr) < 3 {
				fmt.Println("(i*3)+j", (i*3)+j)
				newArr = append(newArr, cardGroup[(i*3)+j])
			}
		}
		finalArr = append(finalArr, newArr)
	}
	return flag, finalArr
}

// 碰牌
func touchCard(curCard model.Card, cardInfo map[string][]int) bool {
	total := 0
	for _, item := range cardInfo[curCard.Type] {
		if curCard.Value == item {
			total++
		}
	}
	if total >= 2 {
		return true
	}
	return false
}

// 明杠
func rightBarCard(curCard model.Card, cardInfo map[string][]int) (bool, [][]model.Card) {
	total := 0
	for _, item := range cardInfo[curCard.Type] {
		if curCard.Value == item {
			total++
		}
	}
	if total == 3 {
		cardGroup := []model.Card{
			{Type: curCard.Type, Value: curCard.Value},
			{Type: curCard.Type, Value: curCard.Value},
			{Type: curCard.Type, Value: curCard.Value},
			{Type: curCard.Type, Value: curCard.Value},
		}
		var finalCards [][]model.Card
		finalCards = append(finalCards, cardGroup)
		return true, finalCards
	}
	return false, nil
}

// 暗杠
func darkBarCard(cardInfo map[string][]int) (bool, [][]model.Card) {
	var finalArr [][]model.Card
	for typ, arr := range cardInfo {
		var cardGroup []model.Card
		for _, item := range arr {
			if len(cardGroup) == 0 {
				cardGroup = append(cardGroup, model.Card{
					Type:  typ,
					Value: item,
				})
			} else if len(cardGroup) == 4 {
				finalArr = append(finalArr, cardGroup)
			} else {
				if item != cardGroup[0].Value {
					cardGroup = []model.Card{}
				}
			}
		}
	}

	return false, finalArr
}

// 用户副手牌，即吃牌碰牌等。
func assistantCards(roomNum int, player, groupType string, cardGroup []model.Card) {
	key := fmt.Sprintf(`%d-assistantCards`, roomNum)
	value := redis.GetValue(key)
	info := make(map[string]map[string][]model.Card)
	if value == "" {
		obj := make(map[string][]model.Card)
		obj[groupType] = cardGroup
		info[player] = obj
		redis.SetValue(key, utils.ToJSON(info), time.Hour)
	} else {
		utils.FromJSON(value, &info)
		if info[player] == nil {
			obj := make(map[string][]model.Card)
			obj[groupType] = cardGroup
			info[player] = obj
		} else {
			info[player][groupType] = append(info[player][groupType], cardGroup...)
		}
	}
	redis.SetValue(key, utils.ToJSON(info), time.Hour)
}

// 胡牌
func huCard(curCard model.Card, cardInfo map[string][]int) bool {
	goldNum := len(model.CardTypeG)
	pairNum := 0
	isXianJin := false
	pairCard := model.Card{}
	for typ, arr := range cardInfo {
		if typ == model.CardTypeG {
			continue
		}
		newArr := make([]int, len(arr))
		copy(newArr, arr)
		//1-9 每张牌的数量
		if typ == curCard.Type {
			newArr = append(newArr, curCard.Value)
			sort.Ints(newArr)
		}
		cardsNum := make([]int, 10)
		for _, card := range newArr {
			cardsNum[card]++
		}
		isHu := computeCards(cardsNum, &isXianJin, &pairNum, &goldNum, &pairCard)
		if !isHu {
			return false
		}
	}
	return true
}

// 自摸
func ziMoCard(cardInfo map[string][]int, curCard model.Card) (bool, string) {
	goldNum := len(cardInfo[model.CardTypeG])
	pairNum := 0
	isXianJin := false
	pairCard := model.Card{}
	for typ, arr := range cardInfo {
		if typ == model.CardTypeG {
			continue
		}
		//1-9 每张牌的数量
		cardsNum := make([]int, 10)
		for _, card := range arr {
			cardsNum[card]++
		}
		isHu := computeCards(cardsNum, &isXianJin, &pairNum, &goldNum, &pairCard)
		if isHu {
			if pairCard.Value != 0 {
				pairCard.Type = typ
			}
		} else {
			if goldNum == 3 {
				return true, model.WinThreeGold
			}
			return false, ""
		}
	}
	if goldNum == 3 {
		return true, model.WinGoldDragon
	} else if goldNum == 2 {
		return true, model.WinGoldSparrow
	}

	if isXianJin && curCard.String() == pairCard.String() {
		return true, model.WinIdleGold
	}
	return true, model.WinZiMo
}

/*
	对子只能存在一对,如果存在金，则可以用金做抵扣成为刻子
	从第一张牌开始计算，假如一个牌有4张，在整个牌里面他只能做刻字和一个顺子；除开 333344445555 这种特殊情况，但是拆分出来也是判断可以胡的。
	所以减去三张牌，SplitCards，这个时候它的第一张牌就只有一张，自然而然的就走找顺子的道路上了。
	但是减去三张发现后面也没有办法胡，看代码继续走下面，再减去2张试试呢。比如 22223344 这种牌
	一张牌它就只能去找后面的顺子，没有就把该张牌记录，往computeCards传，把剩余的牌和金的数量进行组合看能否组成新牌型。
	这里还有一个问题，就是有重复计算的部分
	比如 33334567 的牌，减去三个 3 剩下 34567，减去345剩67 则可以用金做抵扣，如果没有则不能糊；
	在回来减去两个 3 剩下 334567 ，在减去345剩下367不能胡；
	在回来到下面减一个345 剩33367，减去333 剩下67 ，这里和第一次其实是一样的算法，只是顺序不同。
*/

// 根据剩余的牌，和金来重新组合牌
func computeCards(cardsNum []int, isXianJin *bool, pairNum, goldNum *int, pairCard *model.Card) bool {
	cnt := 0
	for _, num := range cardsNum {
		if num > 0 {
			break
		}
		cnt++
	}
	//判断没有牌为可以胡牌
	if len(cardsNum) == cnt {
		return true
	}
	for i := 0; i < len(cardsNum); i++ {
		switch cardsNum[i] {
		case 4:
			isHu := computeThree(cardsNum, i, isXianJin, pairNum, goldNum, pairCard)
			if isHu {
				return isHu
			}
		case 3:
			isHu := computeThree(cardsNum, i, isXianJin, pairNum, goldNum, pairCard)
			if isHu {
				return isHu
			}
		case 2:
			isHu := computeTwo(cardsNum, i, isXianJin, pairNum, goldNum, pairCard)
			if isHu {
				return isHu
			}
		case 1:
			isHu := computeOne(cardsNum, i, isXianJin, pairNum, goldNum, pairCard)
			if isHu {
				return isHu
			}
		}
	}
	return false
}

func computeThree(cardsNum []int, i int, isXianJin *bool, pairNum, goldNum *int, pairCard *model.Card) bool {
	//先判断111刻子能否成功
	cardsNum[i] -= 3
	if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
		return true
	}
	cardsNum[i] += 3

	//如果不成功，则判断11,做对子
	if *pairNum == 0 {
		*pairNum++
		cardsNum[i] -= 2
		if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
			return true
		}
		*pairNum--
		cardsNum[i] += 2
	}

	//如果不成功则判断123，做顺子
	if i+2 < len(cardsNum) && cardsNum[i+1] > 0 && cardsNum[i+2] > 0 {
		cardsNum[i]--
		cardsNum[i+1]--
		cardsNum[i+2]--
		if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
			return true
		}
		cardsNum[i]++
		cardsNum[i+1]++
		cardsNum[i+2]++
	}

	//最后都不成功则开始使用金
	if *pairNum == 1 && *goldNum > 0 {
		cardsNum[i] -= 2
		*goldNum--
		if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
			return true
		}
		cardsNum[i] += 2
		*goldNum++
	}

	return false
}

func computeTwo(cardsNum []int, i int, isXianJin *bool, pairNum, goldNum *int, pairCard *model.Card) bool {

	//先做对子
	if *pairNum == 0 {
		*pairNum++
		cardsNum[i] -= 2
		if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
			return true
		}
		*pairNum--
		cardsNum[i] += 2
	}

	//如果不成功则判断123，做顺子
	if i+2 < len(cardsNum) && cardsNum[i+1] > 0 && cardsNum[i+2] > 0 {
		cardsNum[i]--
		cardsNum[i+1]--
		cardsNum[i+2]--
		if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
			return true
		}
		cardsNum[i]++
		cardsNum[i+1]++
		cardsNum[i+2]++
	}

	//最后都不成功则开始使用金
	if *pairNum == 1 && *goldNum > 0 {
		cardsNum[i] -= 2
		*goldNum--
		if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
			return true
		}
		cardsNum[i] += 2
		*goldNum++
	}

	return false
}

func computeOne(cardsNum []int, i int, isXianJin *bool, pairNum, goldNum *int, pairCard *model.Card) bool {
	//先做顺子
	if i+2 < len(cardsNum) && cardsNum[i+1] > 0 && cardsNum[i+2] > 0 {
		cardsNum[i]--
		cardsNum[i+1]--
		cardsNum[i+2]--
		if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
			return true
		}
		cardsNum[i]++
		cardsNum[i+1]++
		cardsNum[i+2]++
	}

	//不成功，用金做顺子
	if *goldNum > 0 {
		if i+2 < len(cardsNum) && cardsNum[i+1] > 0 && cardsNum[i+2] == 0 {
			cardsNum[i]--
			cardsNum[i+1]--
			*goldNum--
			if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
				return true
			}
			cardsNum[i]++
			cardsNum[i+1]++
			*goldNum++
		} else if i+2 < len(cardsNum) && cardsNum[i+1] == 0 && cardsNum[i+2] > 0 {
			cardsNum[i]--
			cardsNum[i+2]--
			*goldNum--
			if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
				return true
			}
			cardsNum[i]++
			cardsNum[i+2]++
			*goldNum++
		} else if i+2 < len(cardsNum) && cardsNum[i+1] == 0 && cardsNum[i+2] == 0 {
			if *pairNum == 0 && *goldNum > 0 {
				cardsNum[i]--
				*goldNum--
				*pairNum++
				*pairCard = model.Card{
					Value: i,
				}
				*isXianJin = true
				if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
					return true
				}
				*isXianJin = false
				cardsNum[i]++
				*goldNum++
				*pairNum--
			} else if *pairNum == 1 && *goldNum > 1 {
				cardsNum[i]--
				*goldNum -= 2
				if computeCards(cardsNum, isXianJin, pairNum, goldNum, pairCard) {
					return true
				}
				cardsNum[i]++
				*goldNum += 2
			} else {
				return false
			}
		}
	}

	return false
}
