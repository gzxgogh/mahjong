package model

import (
	"fmt"
	"sort"
)

//牌类型
type CardType int

const (
	CardType_Unknown CardType = iota
	CardType_W                //万
	CardType_T                //筒
	CardType_S                //条
	CardType_Z                //中
	CardType_G                //金

)

//实现string
func (st CardType) String() string {
	str := ""
	switch st {
	case CardType_W:
		str += "万"
	case CardType_T:
		str += "筒"
	case CardType_S:
		str += "条"
	case CardType_Z:
		str = "中"
	case CardType_G:
		str = "金"
	}
	return str
}

//牌定义
type Card struct {
	Value int
	Type  CardType
}

//实现string
func (c *Card) String() string {
	return fmt.Sprint("Card: ", c.Value, " ", c.Type)
}

//自定义排序
type SortCards []*Card

//实现
func (sc SortCards) Sort() {
	sort.Slice(sc, func(i, j int) bool {
		return sc[i].Value < sc[j].Value
	})
}
