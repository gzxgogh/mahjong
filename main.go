package main

import (
	"fmt"
	"mahjong/model"
	"mahjong/utils"
)

func main() {

	var totalCardsArr, finalCardsArr []model.Card
	typeArr := []model.CardType{model.CardType_W, model.CardType_T, model.CardType_S}
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
	var groupA, groupB, groupC, groupD []model.Card
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
	fmt.Println("newCardsArr", len(finalCardsArr), finalCardsArr)

	//redis.InItRedisCoon()
	//service.ShuffleCards(1001, 3, "player1")
	//service.ComputeAction(1001)
	//service.TouchCard(1001,"8万")
	//service.EatCard(1001,"5万","player3")
	//service.Dice("a")
	//service.GrabOneCard(1001,"player3")
	//service.PlayOneCard(1001,"9筒","player3")
	//service.HuCard(1001,"player4")
}
