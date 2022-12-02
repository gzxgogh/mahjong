package main

import (
	"mahjong/redis"
	"mahjong/service"
)

func main() {

	redis.InItRedisCoon()
	service.Action{}.ShuffleCards(1001, 3, "player1")

	//result :=service.PlayOneCard(1001,"player3",model.Card{
	//	Type: model.CardType_Z,
	//	Value: 1,
	//})
	//
	//fmt.Println(utils.ToJSON(result))
}
