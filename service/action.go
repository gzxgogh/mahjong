package service

import (
	"fmt"
	"mahjong/model"
	"mahjong/redis"
	"mahjong/utils"
	"sort"
	"time"
)

type Action struct{}

//摇骰子
func (ac Action) Dice() int64 {
	a := utils.GetRandomWithAll(1, 6)
	b := utils.GetRandomWithAll(1, 6)
	sum := a + b
	return sum
}

//洗牌分牌
func (ac Action) ShuffleCards(roomNum, diceNum int, player string) {
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

//获取该局的金
func (ac Action) GetGoldCard(roomNum int) string {
	gold := redis.GetValue(fmt.Sprintf(`%d-glod`, roomNum))
	return gold
}

//抓一张牌
func (ac Action) GrabOneCard(roomNum int, curPlayer string) model.Result {
	surplusCard := GetSurplusCard(roomNum)
	curCard := surplusCard[0]
	fmt.Println("摸到的牌为：", curCard)
	cardInfo := GetPlayerCardInfo(roomNum, curPlayer)
	curCardTypeArr := cardInfo[curCard.Type]
	curCardTypeArr = append(curCardTypeArr, curCard.Value)
	sort.Ints(curCardTypeArr)
	cardInfo[curCard.Type] = curCardTypeArr

	var result model.Result
	result.Player = curPlayer
	var actionArr []string
	if ziMoCard(cardInfo) {
		actionArr = append(actionArr, "ziMo")
	}
	flag, cardGroup := barkBarCard(cardInfo)
	if flag {
		actionArr = append(actionArr, "barCard")
		result.BarCards = cardGroup
	}

	result.Action = actionArr
	surplusCard = append(surplusCard[:0], surplusCard[1:]...)
	key := fmt.Sprintf(`%d-surplusCard`, roomNum)
	redis.SetValue(key, utils.ToJSON(surplusCard), time.Hour)
	redis.SetValue(fmt.Sprintf("%d-%s", roomNum, curPlayer), utils.ToJSON(cardInfo), time.Hour)

	return result
}

//出一张手牌
func (ac Action) PlayOneCard(roomNum int, curPlayer string, curCard model.Card) []model.Result {
	key := fmt.Sprintf("%d-%s", roomNum, curPlayer)
	cardInfo := GetPlayerCardInfo(roomNum, curPlayer)
	curCardArr := cardInfo[curCard.Type]
	for i, item := range curCardArr {
		if curCard.Value == item {
			curCardArr = append(curCardArr[:i], curCardArr[i+1:]...) //用户手牌移除该牌
		}
	}
	cardInfo[curCard.Type] = curCardArr
	redis.SetValue(key, utils.ToJSON(cardInfo), time.Hour)

	var resList []model.Result
	//如果是第一张出牌则可以判断是否可以枪金
	surplusCardArr := GetSurplusCard(roomNum)
	if len(surplusCardArr) == 47 {
		var res model.Result
		res.Player = curPlayer
		if robGold(cardInfo) {
			res.Action = []string{"robGold"}
		}
		resList = append(resList, res)
	} else {
		resList = append(resList, model.Result{})
	}
	//判断其他用户如果获取到该牌是否能胡牌
	nextPlayer := ""
	for i := 0; i < 3; i++ {
		var res model.Result
		nextPlayer = GetNextPlayer(curPlayer)
		res.Player = nextPlayer
		nextCardInfo := GetPlayerCardInfo(roomNum, nextPlayer)

		var actionArr []string
		if huCard(curCard, nextCardInfo) {
			actionArr = append(actionArr, "huCard")
		}
		flag, cardGroup := rightBarCard(curCard, nextCardInfo)
		if flag {
			actionArr = append(actionArr, "barCard")
			res.BarCards = cardGroup
		}
		if touchCard(curCard, nextCardInfo) {
			actionArr = append(actionArr, "touchCard")
		}
		if i == 0 {
			flag, eatCards := eatCard(curCard, nextCardInfo)
			if flag {
				actionArr = append(actionArr, "eatCard")
				res.EatCards = eatCards
			}
		}
		res.Action = actionArr
		resList = append(resList, res)
		curPlayer = nextPlayer
	}
	return resList
}

//吃牌
func (ac Action) EatCard(roomNum int, curCard model.Card, cardGroup []model.Card, player string) {
	cardInfo := GetPlayerCardInfo(roomNum, player)
	arr := cardInfo[curCard.Type]

	var a, b int
	for _, card := range cardGroup {
		if card.Value == curCard.Value {
			continue
		}
		if a == 0 {
			a = card.Value
		} else {
			b = card.Value
		}
	}

	for i, item := range arr {
		if item == a {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	for i, item := range arr {
		if item == b {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}

	cardInfo[curCard.Type] = arr
	key := fmt.Sprintf(`%d-%s`, roomNum, player)
	redis.SetValue(key, utils.ToJSON(cardInfo), 1*time.Hour)

	return
}

//碰牌
func (ac Action) TouchCard(roomNum int, curCard model.Card, player string) {
	cardInfo := GetPlayerCardInfo(roomNum, player)
	arr := cardInfo[curCard.Type]
	var newArr []int
	total := 0
	for _, item := range arr {
		if item == curCard.Value && total <= 2 {
			total++
			continue
		}
		newArr = append(newArr, item)
	}
	cardInfo[curCard.Type] = newArr

	key := fmt.Sprintf(`%d-%s`, roomNum, player)
	fmt.Println("key", key)
	fmt.Println("cardInfo", cardInfo)
	redis.SetValue(key, utils.ToJSON(cardInfo), 1*time.Hour)
}

//杠牌
func (ac Action) BarCard(roomNum int, curCard model.Card, player string) model.Result {
	cardInfo := GetPlayerCardInfo(roomNum, player)
	arr := cardInfo[curCard.Type]
	var newArr []int
	total := 0
	for _, item := range arr {
		if item == curCard.Value && item <= 3 {
			total++
			continue
		}
		newArr = append(newArr, item)
	}
	cardInfo[curCard.Type] = newArr

	//杠玩后马上从背后摸一张，判断能不能杠上开花
	res := model.Result{
		Player: player,
	}
	surplusCardArr := GetSurplusCard(roomNum)
	length := len(surplusCardArr)
	newCard := surplusCardArr[length-1]

	caryTypeArr := cardInfo[newCard.Type]
	caryTypeArr = append(caryTypeArr, newCard.Value)
	sort.Ints(caryTypeArr)
	cardInfo[newCard.Type] = caryTypeArr

	if ziMoCard(cardInfo) {
		res.Action = []string{"ziMo"}
	}
	//扣减牌数，重进记录
	surplusCardArr = append(surplusCardArr[:length-1], surplusCardArr[length:]...)
	key := fmt.Sprintf(`%d-surplusCard`, roomNum)
	redis.SetValue(key, utils.ToJSON(surplusCardArr), 1*time.Hour)

	//重新存入用户手牌
	key = fmt.Sprintf(`%d-%s`, roomNum, player)
	redis.SetValue(key, utils.ToJSON(cardInfo), 1*time.Second)

	return res
}
