package main

import (
	"mahjong/service"
)

func main() {

	//redis.InItRedisCoon()
	service.ShuffleCards(1001, 3, "player1")
	//service.ComputeAction(1001)
	//service.TouchCard(1001,"8万")
	//service.EatCard(1001,"5万","player3")
	//service.Dice("a")
	//service.GrabOneCard(1001,"player3")
	//service.PlayOneCard(1001,"9筒","player3")
	//service.HuCard(1001,"player4")
}
