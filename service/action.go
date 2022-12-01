package service

import (
	"fmt"
	"mahjong/model"
	"mahjong/redis"
	"mahjong/utils"
	"sort"
	"strconv"
	"time"
)

//摇骰子
func Dice(player string) int64 {
	a := utils.GetRandomWithAll(1, 6)
	b := utils.GetRandomWithAll(1, 6)
	sum := a + b
	return sum
}

//洗牌分牌
func ShuffleCards(roomNum, diceNum int, player string) {
	var totalCardsArr, finalCardsArr []model.Card
	typeArr := []string{model.CardType_W, model.CardType_T, model.CardType_S}
	for _, item := range typeArr {
		for i := 0; i < 4; i++ {
			for value := 1; value <= 9; value++ {
				totalCardsArr = append(totalCardsArr, model.Card{
					Type:  item,
					Value: value,
				})
			}
		}
	}
	for i := 0; i < 4; i++ {
		totalCardsArr = append(totalCardsArr, model.Card{
			Type:  model.CardType_Z,
			Value: 1,
		})
	}
	for i := 112; i > 0; i-- {
		randomNum := utils.GetRandomWithAll(0, i-1)
		finalCardsArr = append(finalCardsArr, totalCardsArr[randomNum])
		totalCardsArr = append(totalCardsArr[:randomNum], totalCardsArr[(randomNum+1):]...)
	}

	/*var groupA, groupB, groupC, groupD []model.Card
	for i, item := range finalCardsArr {
		if i < 28 {
			groupA = append(groupA, item)
		} else if i >= 28 && i < 56 {
			groupB = append(groupB, item)
		} else if i >= 56 && i < 84 {
			groupC = append(groupC, item)
		} else if i >= 84 && i < 112 {
			groupD = append(groupD, item)
		}
	}
	fmt.Println("用户1面前的牌堆", groupA)
	fmt.Println("用户2面前的牌堆", groupB)
	fmt.Println("用户3面前的牌堆", groupC)
	fmt.Println("用户4面前的牌堆", groupD)
	*/

	//确定是那个用户摇的骰子，并且更具点数开始抓牌
	var startGroupNum, startNum int
	switch player {
	case "player1":
		startGroupNum = GetStartGroup(diceNum, 1)
	case "player2":
		startGroupNum = GetStartGroup(diceNum, 2)
	case "player3":
		startGroupNum = GetStartGroup(diceNum, 3)
	case "player4":
		startGroupNum = GetStartGroup(diceNum, 4)
	}
	startNum = diceNum * 2

	fmt.Println("从", startGroupNum, "个用户牌堆的第", startNum+1, "开始抓牌")
	GrabTheCard(roomNum, startGroupNum, startNum, finalCardsArr)
}

//获取从哪个用户开始抓牌
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

//抓牌并且分牌
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
	gold := surplusCardArr[len(surplusCardArr)-1].String()
	keyPlayerCards[startGroupNum] = append(keyPlayerCards[startGroupNum], surplusCardArr[0])
	surplusCardArr = append(surplusCardArr[:0], surplusCardArr[(0+1):]...)

	for playerNum, arr := range keyPlayerCards {
		kInfo := make(map[string][]int)
		for _, item := range arr {
			if item.String() == gold {
				kInfo[model.CardType_G] = append(kInfo[model.CardType_G], 1)
			} else if item.Type != model.CardType_Z && item.String() != gold {
				kInfo[item.Type] = append(kInfo[item.Type], item.Value)
			} else {
				kInfo[model.CardType_Z] = append(kInfo[model.CardType_Z], 1)
			}
		}
		for _, v := range kInfo {
			sort.Ints(v)
		}
		fmt.Println("用户player", playerNum, "的手牌为:", utils.ToJSON(kInfo))
		//存入用户手牌
		//redis.SetValue(fmt.Sprintf(`%d-player%d`, roomNum, k), utils.ToJSON(kInfo), 1*time.Hour)
	}

	//存入分配玩后各个玩家手里的牌，和场上现有的牌
	//redis.SetValue(fmt.Sprintf(`%d-glod`, roomNum), gold, 1*time.Hour)
	//redis.SetValue(fmt.Sprintf(`%d-surplusCard`, roomNum), utils.ToJSON(surplusCard), 1*time.Hour)
}

//获取剩余牌堆的牌
func GetSurplusCard(roomNum int) []model.Card {
	value := redis.GetValue(fmt.Sprintf(`%d-surplusCard`, roomNum))
	var surplusCard []model.Card
	utils.FromJSON(value, &surplusCard)
	return surplusCard
}

//获取用户手牌
func GetPlayerCardInfo(roomNum int, player string) map[string][]int {
	value := redis.GetValue(fmt.Sprintf(`%d-%s`, roomNum, player))
	cardInfo := make(map[string][]int)
	utils.FromJSON(value, &cardInfo)
	return cardInfo
}

//抢金
func robGold(goldCardType string, goldCardNum int, cardInfo map[string][]int) bool {

	return false
}

//吃牌,只有下家可以吃牌
func EatCard(roomNum int, curPlayer string, curCard model.Card) map[string]interface{} {
	//确定下家是谁
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
	cardInfo := GetPlayerCardInfo(roomNum, player)

	var lessTwoCard, lessOneCard, greaterOneCard, greaterTwoCard int
	result := make(map[string]interface{})
	var cardGroup []model.Card
	switch curCard.Value {
	case 1:
		for _, item := range cardInfo[curCard.Type] {
			if item == 2 {
				greaterOneCard = item
			} else if item == 3 {
				greaterOneCard = item
			}
		}
		if greaterOneCard != 0 && greaterTwoCard != 0 {
			cardGroup = []model.Card{
				model.Card{
					Type:  curCard.Type,
					Value: 1,
				},
				model.Card{
					Type:  curCard.Type,
					Value: 2,
				},
				model.Card{
					Type:  curCard.Type,
					Value: 3,
				},
			}
		}
	case 2:
		for _, item := range cardInfo[curCard.Type] {
			if item == 1 {
				lessOneCard = item
			} else if item == 3 {
				greaterOneCard = item
			}
		}
		if lessOneCard != 0 && greaterOneCard != 0 {

		}
		if lessOneCard != 0 && lessTwoCard != 0 {

		}
	case 8:
		for _, item := range cardInfo[curCard.Type] {
			if item == 7 {
				lessOneCard = item
			} else if item == 9 {
				greaterOneCard = item
			}
		}
		if lessOneCard != 0 && greaterOneCard != 0 {

		}
		if lessOneCard != 0 && lessTwoCard != 0 {

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
	}
	result[player] = map[string]interface{}{
		"action":    []string{"eatCard"},
		"cardGroup": cardGroup,
	}

	return result
}

//碰牌
func TouchCard(roomNum int, curCard string) map[string][]string {
	curCardType := curCard[1:]
	curCardNum, _ := strconv.Atoi(curCard[:1])

	result := make(map[string][]string)
	for i := 1; i <= 4; i++ {
		cardInfo := GetPlayerCardInfo(roomNum, fmt.Sprintf(`player%d`, i))
		total := 0
		for _, item := range cardInfo[curCardType] {
			if curCardNum == item {
				total++
			}
		}
		if total == 2 {
			key := fmt.Sprintf(`player%d`, i)
			result[key] = append(result[key], "touchCard")
			break
		}
	}
	return result
}

//抓一张牌
func GrabOneCard(roomNum int, curPlayer string) {
	surplusCard := GetSurplusCard(roomNum)
	curCard := surplusCard[0]

	cardInfo := GetPlayerCardInfo(roomNum, curPlayer)
	curCardTypeArr := cardInfo[curCard.Type]
	curCardTypeArr = append(curCardTypeArr, curCard.Value)
	sort.Ints(curCardTypeArr)

	cardInfo[curCard.Type] = curCardTypeArr
	//todo 判断当前用户是否可以胡牌,不行重新存入剩余的牌堆

	surplusCard = append(surplusCard[:0], surplusCard[1:]...)
	key := fmt.Sprintf(`%d-surplusCard`, roomNum)
	redis.DelKey(key)
	redis.SetValue(key, utils.ToJSON(surplusCard), time.Hour)
}

//出一张手牌
func PlayOneCard(roomNum int, curCard, curPlayer string) {
	cardInfo := GetPlayerCardInfo(roomNum, curPlayer)

	curCardType := curCard[1:]
	curCardNum, _ := strconv.Atoi(curCard[:1])

	arr := cardInfo[curCardType]
	for i, item := range arr {
		if curCardNum == item {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	cardInfo[curCardType] = arr
	//todo 判断除该用户外其他用户是否可以胡牌，碰牌，下家胡牌

}

//胡牌
func HuCard(roomNum int, curPlayer string) bool {
	cardInfo := GetPlayerCardInfo(roomNum, curPlayer)
	goldNum := len(cardInfo["金"])
	if goldNum == 3 {
		return true
	}
	pairNum := 0
	for _, arr := range cardInfo {
		//1-9 每张牌的数量
		cardsNum := make([]int, 10)
		for _, card := range arr {
			cardsNum[card]++
		}
		isHu := ComputeCards(cardsNum, pairNum, goldNum)
		if !isHu {
			return false
		}
	}
	return true
}

/*
	对子只能存在一对,如果存在金，则可以用金做抵扣
	从第一张牌开始计算，假如一个牌有4张，在整个牌里面他只能做刻字和一个顺子；除开 333344445555 这种特殊情况，但是拆分出来也是判断可以胡的。
	所以减去三张牌，ComputeCards，这个时候它的第一张牌就只有一张，自然而然的就走找顺子的道路上了。
	但是减去三张发现后面也没有办法胡，看代码继续走下面，再减去2张试试呢。比如 22223344 这种牌
	一张牌它就只能去找后面的顺子，没有就不能胡。
	这里还有一个问题，就是有重复计算的部分
	比如 33334567 的牌，减去三个 3 剩下 34567，减去345剩67 则可以用金做抵扣，如果没有则不能糊；
	在回来减去两个 3 剩下 334567 ，在减去345剩下367不能胡；
	在回来到下面减一个345 剩33367，减去333 剩下67 ，这里和第一次其实是一样的算法，只是顺序不同。
*/
func ComputeCards(cardsNum []int, pairNum, goldNum int) bool {
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
			fallthrough
		case 3:
			//这种存在这几种情况，可以加后面成顺子，取两张为对子，或取一个刻字
			//减掉后再传入SplitCards
			cardsNum[i] -= 3
			if ComputeCards(cardsNum, pairNum, goldNum) {
				return true
			}
			cardsNum[i] += 3
			//这种不行就向下传递。。。
			fallthrough
		case 2:
			if pairNum == 0 {
				pairNum++
				cardsNum[i] -= 2
				if ComputeCards(cardsNum, pairNum, goldNum) {
					return true
				}
				cardsNum[i] += 2
			}
			if pairNum == 1 && goldNum == 0 {
				return false
			}
			if pairNum > 0 && goldNum > 0 {
				cardsNum[i] -= 2
				goldNum--
				if ComputeCards(cardsNum, pairNum, goldNum) {
					return true
				}
				cardsNum[i] += 2
				goldNum++
			}
			fallthrough
		case 1:
			if i+2 < len(cardsNum) && cardsNum[i+1] > 0 && cardsNum[i+2] > 0 {
				cardsNum[i]--
				cardsNum[i+1]--
				cardsNum[i+2]--
				if ComputeCards(cardsNum, pairNum, goldNum) {
					return true
				}
				cardsNum[i]++
				cardsNum[i+1]++
				cardsNum[i+2]++
			}
			//如果发现普通的值不够则使用金
			if i+2 < len(cardsNum) && cardsNum[i+1] > 0 && goldNum > 0 {
				cardsNum[i]--
				cardsNum[i+1]--
				goldNum--
				if ComputeCards(cardsNum, pairNum, goldNum) {
					return true
				}
				cardsNum[i]++
				cardsNum[i+1]++
				goldNum++
			}
			if i+2 < len(cardsNum) && cardsNum[i+2] > 0 && goldNum > 0 {
				cardsNum[i]--
				cardsNum[i+2]--
				goldNum--
				if ComputeCards(cardsNum, pairNum, goldNum) {
					return true
				}
				cardsNum[i]++
				cardsNum[i+2]++
				goldNum++
			}
		}
	}
	return false
}
