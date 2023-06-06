package model

import (
	"fmt"
	"sort"
)

// 牌类型
type CardType string

const (
	CardTypeW = "万"
	CardTypeT = "筒"
	CardTypeS = "条"
	CardTypeZ = "中"
	CardTypeG = "金"
)

const (
	WinHu          = "胡"
	WinZiMo        = "自摸"
	WinRobGold     = "抢金"
	WinThreeGold   = "三金"
	WinIdleGold    = "闲金"
	WinGoldSparrow = "金雀"
	WinGoldDragon  = "金龙"
)

const (
	ActionEat   = "吃"
	ActionTouch = "碰"
	ActionBar   = "杠"
)

// 实现string
func (st CardType) String() string {
	str := ""
	switch st {
	case CardTypeW:
		str += "万"
	case CardTypeT:
		str += "筒"
	case CardTypeS:
		str += "条"
	case CardTypeZ:
		str = "中"
	case CardTypeG:
		str = "金"
	}
	return str
}

// 牌定义
type Card struct {
	Value int    `json:"value"`
	Type  string `json:"type"`
}

// 实现string
func (c *Card) String() string {
	return fmt.Sprint("Card: ", c.Value, " ", c.Type)
}

// 自定义排序
type SortCards []*Card

// 实现
func (sc SortCards) Sort() {
	sort.Slice(sc, func(i, j int) bool {
		return sc[i].Value < sc[j].Value
	})
}

type Action struct {
	Player   string   `json:"player"`
	Action   []string `json:"action"`
	GardCard *Card    `json:"gardCard"`
	EatCards [][]Card `json:"eatCards"`
	BarCards [][]Card `json:"barCards"`
}

type ShuffleCardsReq struct {
	RoomNum int    `json:"roomNum"`
	DiceNum int    `json:"diceNum"`
	Player  string `json:"player"`
}

type JudgeRobGoldReq struct {
	RoomNum int `json:"roomNum"`
}

type GrabOneCardReq struct {
	RoomNum int    `json:"roomNum"`
	Player  string `json:"player"`
}

type PlayOneCardReq struct {
	RoomNum int    `json:"roomNum"`
	Player  string `json:"player"`
	CurCard Card   `json:"curCard"`
}

type EatCardReq struct {
	RoomNum   int    `json:"roomNum"`
	Player    string `json:"player"`
	CurCard   Card   `json:"curCard"`
	CardGroup []Card `json:"cardGroup"`
}

type TouchCardReq struct {
	RoomNum int    `json:"roomNum"`
	Player  string `json:"player"`
	CurCard Card   `json:"curCard"`
}

type BarCardReq struct {
	RoomNum int    `json:"roomNum"`
	Player  string `json:"player"`
	BarType string `json:"barType"` //rightBar/darkBar
	CurCard Card   `json:"curCard"`
}
