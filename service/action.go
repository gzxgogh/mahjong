package service

import (
	"fmt"
	"github.com/gzxgogh/ggin/models"
	"github.com/gzxgogh/ggin/utils"
	"mahjong/model"
	"mahjong/redis"
	"sort"
	"time"
)

type action struct{}

var Action *action

// 摇骰子
func (ac *action) Dice() int64 {
	a := utils.GetRandomWithAll(1, 6)
	b := utils.GetRandomWithAll(1, 6)
	sum := a + b
	return sum
}

// 洗牌分牌
func (ac *action) ShuffleCards(roomNum, diceNum int, player string) models.Result {
	redis.DelKey(fmt.Sprintf(`%d-assistantCards`, roomNum))
	var totalCardsArr, finalCardsArr []model.Card
	typeArr := []string{model.CardTypeW, model.CardTypeT, model.CardTypeS}
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
			Type:  model.CardTypeZ,
			Value: 1,
		})
	}
	for i := 112; i > 0; i-- {
		randomNum := utils.GetRandomWithAll(0, i-1)
		finalCardsArr = append(finalCardsArr, totalCardsArr[randomNum])
		totalCardsArr = append(totalCardsArr[:randomNum], totalCardsArr[(randomNum+1):]...)
	}

	//确定是那个用户摇的骰子，并且根据点数开始抓牌
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

	GrabTheCard(roomNum, startGroupNum, startNum, finalCardsArr)

	//判断能否抢金
	var resList []model.Action
	for i := 1; i <= 4; i++ {
		player := fmt.Sprintf(`player%d`, i)
		cardInfo := GetPlayerCardInfo(roomNum, player)
		if robGold(cardInfo) {
			res := model.Action{
				Player: player,
				Action: []string{model.WinRobGold},
			}
			resList = append(resList, res)
		}
	}
	return models.Success(resList)
}

// 获取该局的金
func (ac *action) GetGoldCard(roomNum int) models.Result {
	gold := GetGoldCard(roomNum)
	return models.Success(gold)
}

// 抓一张牌
func (ac *action) GrabOneCard(roomNum int, curPlayer string) models.Result {
	surplusCard := GetSurplusCard(roomNum)
	curCard := surplusCard[0]
	gold := GetGoldCard(roomNum)

	cardInfo := GetPlayerCardInfo(roomNum, curPlayer)
	if gold.String() == curCard.String() {
		curCard = model.Card{
			Type:  model.CardTypeG,
			Value: 1,
		}
	}
	fmt.Println("摸到的牌为：", curCard)
	curCardTypeArr := cardInfo[curCard.Type]
	curCardTypeArr = append(curCardTypeArr, curCard.Value)
	sort.Ints(curCardTypeArr)
	cardInfo[curCard.Type] = curCardTypeArr

	var result model.Action
	result.Player = curPlayer
	result.GardCard = &curCard
	var actionArr []string
	flag, actionStr := ziMoCard(cardInfo, curCard)
	if flag {
		actionArr = append(actionArr, actionStr)
	}
	flag, cardGroup := darkBarCard(cardInfo)
	if flag {
		actionArr = append(actionArr, model.ActionBar)
		result.BarCards = cardGroup
	}

	result.Action = actionArr
	surplusCard = append(surplusCard[:0], surplusCard[1:]...)
	key := fmt.Sprintf(`%d-surplusCard`, roomNum)
	redis.SetValue(key, utils.ToJSON(surplusCard), time.Hour)
	redis.SetValue(fmt.Sprintf("%d-%s", roomNum, curPlayer), utils.ToJSON(cardInfo), time.Hour)

	return models.Success(result)
}

// 出一张手牌
func (ac *action) PlayOneCard(roomNum int, curPlayer string, curCard model.Card) models.Result {
	key := fmt.Sprintf("%d-%s", roomNum, curPlayer)
	cardInfo := GetPlayerCardInfo(roomNum, curPlayer)
	curCardArr := cardInfo[curCard.Type]
	for i, item := range curCardArr {
		if curCard.Value == item {
			curCardArr = append(curCardArr[:i], curCardArr[i+1:]...) //用户手牌移除该牌
			break
		}
	}
	cardInfo[curCard.Type] = curCardArr
	redis.SetValue(key, utils.ToJSON(cardInfo), time.Hour)

	var resList []model.Action
	//如果是第一张出牌则可以判断是否可以枪金
	surplusCardArr := GetSurplusCard(roomNum)
	if len(surplusCardArr) == 47 {
		var res model.Action
		res.Player = curPlayer
		if robGold(cardInfo) {
			res.Action = []string{model.WinRobGold}
		}
		resList = append(resList, res)
	} else {
		resList = append(resList, model.Action{
			Player: curPlayer,
		})
	}
	//判断其他用户如果获取到该牌是否能胡牌
	nextPlayer := ""
	for i := 0; i < 3; i++ {
		var res model.Action
		nextPlayer = GetNextPlayer(curPlayer)
		res.Player = nextPlayer
		nextCardInfo := GetPlayerCardInfo(roomNum, nextPlayer)

		var actionArr []string
		if huCard(curCard, nextCardInfo) {
			actionArr = append(actionArr, model.WinHu)
		}
		flag, cardGroup := rightBarCard(curCard, nextCardInfo)
		if flag {
			actionArr = append(actionArr, model.ActionBar)
			res.BarCards = cardGroup
		}
		if touchCard(curCard, nextCardInfo) {
			actionArr = append(actionArr, model.ActionTouch)
		}
		if i == 0 {
			flag, eatCards := eatCard(curCard, nextCardInfo)
			if flag {
				actionArr = append(actionArr, model.ActionEat)
				res.EatCards = eatCards
			}
		}
		res.Action = actionArr
		resList = append(resList, res)
		curPlayer = nextPlayer
	}
	return models.Success(resList)
}

// 吃牌
func (ac *action) EatCard(roomNum int, player string, curCard model.Card, cardGroup []model.Card) models.Result {
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
			break
		}
	}
	for i, item := range arr {
		if item == b {
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}

	cardInfo[curCard.Type] = arr
	key := fmt.Sprintf(`%d-%s`, roomNum, player)
	redis.SetValue(key, utils.ToJSON(cardInfo), 1*time.Hour)
	assistantCards(roomNum, player, "eatCard", cardGroup)

	return models.Success(nil)
}

// 碰牌
func (ac *action) TouchCard(roomNum int, player string, curCard model.Card) models.Result {
	cardInfo := GetPlayerCardInfo(roomNum, player)
	arr := cardInfo[curCard.Type]
	var newArr []int
	total := 0
	cardGroup := []model.Card{curCard}
	for _, value := range arr {
		if value == curCard.Value && total <= 2 {
			cardGroup = append(cardGroup, model.Card{
				Type:  curCard.Type,
				Value: value,
			})
			total++
			continue
		}
		newArr = append(newArr, value)
	}
	cardInfo[curCard.Type] = newArr

	key := fmt.Sprintf(`%d-%s`, roomNum, player)
	fmt.Println("key", key)
	fmt.Println("cardInfo", cardInfo)
	redis.SetValue(key, utils.ToJSON(cardInfo), 1*time.Hour)
	assistantCards(roomNum, player, "touchCard", cardGroup)

	return models.Success(nil)
}

// 明杠
func (ac *action) BarCard(roomNum int, player, barType string, curCard model.Card) models.Result {
	cardInfo := GetPlayerCardInfo(roomNum, player)
	arr := cardInfo[curCard.Type]
	var newArr []int
	cardGroup := []model.Card{curCard}
	total := 3
	if barType == "darkBar" {
		total = 4
	}
	for _, value := range arr {
		if value == curCard.Value && total > 0 {
			cardGroup = append(cardGroup, model.Card{
				Type:  curCard.Type,
				Value: value,
			})
			total--
			continue
		}
		newArr = append(newArr, value)
	}
	if total != 0 {
		return models.Error(-1, "无效的杠牌")
	}
	cardInfo[curCard.Type] = newArr

	//杠玩后马上从背后摸一张，判断能不能杠上开花
	res := model.Action{
		Player: player,
	}
	surplusCardArr := GetSurplusCard(roomNum)
	length := len(surplusCardArr)
	newCard := surplusCardArr[length-1]
	fmt.Println("摸到的牌:", newCard.String())
	caryTypeArr := cardInfo[newCard.Type]
	caryTypeArr = append(caryTypeArr, newCard.Value)
	sort.Ints(caryTypeArr)
	cardInfo[newCard.Type] = caryTypeArr

	flag, actionStr := ziMoCard(cardInfo, newCard)
	if flag {
		res.Action = []string{actionStr}
	}
	//扣减牌数，重进记录
	surplusCardArr = append(surplusCardArr[:length-1], surplusCardArr[length:]...)
	key := fmt.Sprintf(`%d-surplusCard`, roomNum)
	redis.SetValue(key, utils.ToJSON(surplusCardArr), 1*time.Hour)

	//重新存入用户手牌
	key = fmt.Sprintf(`%d-%s`, roomNum, player)
	redis.SetValue(key, utils.ToJSON(cardInfo), 1*time.Hour)
	assistantCards(roomNum, player, "barType", cardGroup)
	res.GardCard = &newCard
	return models.Success(res)
}

// 记录弃牌
func (ac *action) RecordAbandonCard(roomNum int, player string, curCard model.Card) models.Result {
	key := fmt.Sprintf("%d-abandonCards", roomNum)
	value := redis.GetValue(key)
	if value == "" {
		obj := make(map[string][]model.Card)
		obj[player] = []model.Card{curCard}
		redis.SetValue(key, utils.ToJSON(obj), time.Hour)
	} else {
		obj := make(map[string][]model.Card)
		utils.FromJSON(value, &obj)
		if len(obj[player]) == 0 {
			obj[player] = []model.Card{curCard}
		} else {
			obj[player] = append(obj[player], curCard)
		}
		redis.SetValue(key, utils.ToJSON(obj), time.Hour)
	}
	return models.Success(nil)
}

// 获取弃牌堆
func (ac *action) GetAbandonCards(roomNum int) models.Result {
	key := fmt.Sprintf("%d-abandonCards", roomNum)
	value := redis.GetValue(key)
	obj := make(map[string][]model.Card)
	utils.FromJSON(value, &obj)

	return models.Success(obj)
}

// 获取用户手牌
func (ac *action) GetPlayerCards(roomNum int) models.Result {
	assistantStr := redis.GetValue(fmt.Sprintf(`%d-assistantCards`, roomNum))
	assistantInfo := make(map[string]interface{})
	utils.FromJSON(assistantStr, &assistantInfo)

	result := make(map[string]interface{})
	gold := GetGoldCard(roomNum)
	for i := 1; i <= 4; i++ {
		player := fmt.Sprintf("player%d", i)
		cardInfo := GetPlayerCardInfo(roomNum, player)
		playerCard := make(map[string]interface{})
		cardsArr := make([]map[string]interface{}, 0)
		for j := 0; j < len(cardInfo[model.CardTypeG]); j++ {
			cardsArr = append(cardsArr, map[string]interface{}{
				"type":  gold.Type,
				"value": gold.Value,
				"gold":  true,
			})
		}
		for typ, arr := range cardInfo {
			if typ == model.CardTypeG {
				continue
			}
			for _, value := range arr {
				cardsArr = append(cardsArr, map[string]interface{}{
					"type":  typ,
					"value": value,
				})
			}
		}
		playerCard["main"] = cardsArr
		playerCard["assistant"] = assistantInfo[player]
		result[player] = playerCard
	}
	return models.Success(result)
}
