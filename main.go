package main

import (
	"fmt"
	"mahjong/model"
	"mahjong/utils"
)

func main() {

	var cardsArr, newCardsArr []model.Card
	typeArr := []model.CardType{model.CardType_W, model.CardType_T, model.CardType_S}
	for _, item := range typeArr {
		for i := 0; i < 4; i++ {
			for value := 1; value <= 9; value++ {
				cardsArr = append(cardsArr, model.Card{
					Type:  item,
					Value: value,
				})
			}
		}
	}
	for i := 0; i < 4; i++ {
		cardsArr = append(cardsArr, model.Card{
			Type:  model.CardType_Z,
			Value: 1,
		})
	}

	for i := 112; i > 0; i-- {
		randomNum := utils.GetRandomWithAll(0, i-1)
		newCardsArr = append(newCardsArr, cardsArr[randomNum])
		cardsArr = append(cardsArr[:randomNum], cardsArr[(randomNum+1):]...)
	}
	fmt.Println("newCardsArr", len(newCardsArr), newCardsArr)

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
