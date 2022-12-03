package main

import (
	"fmt"
	"mahjong/redis"
	"mahjong/service"
	"mahjong/utils"
)

func main() {

	redis.InItRedisCoon()
	//service.Action{}.ShuffleCards(1001, 3, "player1")

	//出牌
	//result :=service.Action{}.PlayOneCard(1001,"player2",model.Card{
	//	Type: model.CardType_Z,
	//	Value: 1,
	//})

	//摸排
	result := service.Action{}.GrabOneCard(1001, "player3")

	//碰排
	//service.Action{}.TouchCard(1001,model.Card{Type: model.CardType_Z, Value: 1},"player4")

	//吃排
	//cardGroup :=[]model.Card{
	//	{Type: model.CardType_W, Value: 2},
	//	{Type: model.CardType_W, Value: 3},
	//	{Type: model.CardType_W, Value: 4},
	//}
	//fmt.Println(cardGroup)
	//service.Action{}.EatCard(1001,model.Card{Type: model.CardType_W, Value: 4}, cardGroup,"player1")

	fmt.Println(utils.ToJSON(result))
}
