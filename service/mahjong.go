package service

import (
	"fmt"
	"mahjong/redis"
	"mahjong/utils"
	"sort"
	"strconv"
	"time"
)

var tenThousandArr = []string{"1万", "1万", "1万", "1万", "2万", "2万", "2万", "2万", "3万", "3万", "3万", "3万", "4万", "4万", "4万", "4万", "5万", "5万", "5万", "5万", "6万", "6万", "6万", "6万", "7万", "7万", "7万", "7万", "8万", "8万", "8万", "8万", "9万", "9万", "9万", "9万"}
var canisterArr = []string{"1筒", "1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "5筒", "6筒", "6筒", "6筒", "6筒", "7筒", "7筒", "7筒", "7筒", "8筒", "8筒", "8筒", "8筒", "9筒", "9筒", "9筒", "9筒"}
var stripArr = []string{"1条", "1条", "1条", "1条", "2条", "2条", "2条", "2条", "3条", "3条", "3条", "3条", "4条", "4条", "4条", "4条", "5条", "5条", "5条", "5条", "6条", "6条", "6条", "6条", "7条", "7条", "7条", "7条", "8条", "8条", "8条", "8条", "9条", "9条", "9条", "9条"}
var decorArr = []string{"中", "中", "中", "中"}

//摇骰子
func Dice(player string) int64 {
	a := utils.GetRandomWithAll(1, 6)
	b := utils.GetRandomWithAll(1, 6)
	sum := a + b
	return sum
}

//洗牌分牌
func ShuffleCards(roomNum, diceNum int, player string) {
	var cardsArr []string
	cardsArr = append(cardsArr, tenThousandArr...)
	cardsArr = append(cardsArr, canisterArr...)
	cardsArr = append(cardsArr, stripArr...)
	cardsArr = append(cardsArr, decorArr...)

	//洗牌
	total := len(cardsArr)
	var newCardsArr []string
	for i := total; i > 0; i-- {
		randomNum := utils.GetRandomWithAll(0, i-1)
		newCardsArr = append(newCardsArr, cardsArr[randomNum])
		cardsArr = append(cardsArr[:randomNum], cardsArr[(randomNum+1):]...)
	}

	//分成四组，每组28张
	var groupA, groupB, groupC, groupD []string
	for i, item := range newCardsArr {
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
	fmt.Println("用户1面前的牌堆", groupA)
	fmt.Println("用户2面前的牌堆", groupB)
	fmt.Println("用户3面前的牌堆", groupC)
	fmt.Println("用户4面前的牌堆", groupD)
	fmt.Println("从", startGroupNum, "个用户牌堆的第", startNum+1, "开始抓牌")
	GrabTheCard(roomNum, startGroupNum, startNum, newCardsArr)
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
func GrabTheCard(roomNum, startGroupNum, startNum int, allCardsArr []string) {
	var newCardsArr, surplusCard []string
	var grabTheCardArr []int
	switch startGroupNum {
	case 1:
		startNum = 6
		grabTheCardArr = []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	case 2:
		startNum = 6 + 28
		grabTheCardArr = []int{2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1}
	case 3:
		startNum = 6 + (28 * 2)
		grabTheCardArr = []int{3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2}
	case 4:
		startNum = 6 + (28 * 3)
		grabTheCardArr = []int{4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3}
	}
	//重新排序，用户从第0个牌开始抓取即可
	newCardsArr = append(newCardsArr, allCardsArr[startNum:]...)
	newCardsArr = append(newCardsArr, allCardsArr[:startNum]...)
	keyArr := make(map[int][]string)
	curNum := 0
	for _, v := range grabTheCardArr {
		for i := 0; i < 4; i++ {
			keyArr[v] = append(keyArr[v], newCardsArr[curNum])
			curNum++
		}
	}

	//剩余的牌数,最后一张为金,庄家多模第一张门牌
	surplusCard = append(surplusCard, newCardsArr[curNum:]...)
	gold := surplusCard[len(surplusCard)-1]

	keyArr[startGroupNum] = append(keyArr[startGroupNum], surplusCard[0])
	surplusCard = append(surplusCard[:0], surplusCard[(0+1):]...)

	for k, arr := range keyArr {
		kInfo := make(map[string][]int)
		for _, item := range arr {
			if item == gold {
				kInfo["金"] = append(kInfo["金"], 1)
			} else if item != "中" && item != gold {
				cardType := item[1:]
				cardNum, _ := strconv.Atoi(item[:1])
				kInfo[cardType] = append(kInfo[cardType], cardNum)
			} else {
				kInfo["中"] = append(kInfo["中"], 1)
			}
		}
		for _, v := range kInfo {
			sort.Ints(v)
		}
		fmt.Println("用户", k, "的手牌为:", utils.ToJSON(kInfo))
		//存入用户手牌
		redis.SetValue(fmt.Sprintf(`%d-player%d`, roomNum, k), utils.ToJSON(kInfo), 1*time.Hour)
	}

	//存入分配玩后各个玩家手里的牌，和场上现有的牌
	redis.SetValue(fmt.Sprintf(`%d-glod`, roomNum), gold, 1*time.Hour)
	redis.SetValue(fmt.Sprintf(`%d-surplusCard`, roomNum), utils.ToJSON(surplusCard), 1*time.Hour)
}

//获取剩余牌堆的牌
func GetSurplusCard(roomNum int) []string {
	value := redis.GetValue(fmt.Sprintf(`%d-surplusCard`, roomNum))
	var surplusCard []string
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
func EatCard(roomNum int, curCard, curPlayer string) map[string]interface{} {
	curCardType := curCard[1:]
	curCardNum, _ := strconv.Atoi(curCard[:1])

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
	var cardGroup []string
	switch curCardNum {
	case 1:
		for _, item := range cardInfo[curCardType] {
			if item == 2 {
				greaterOneCard = item
			} else if item == 3 {
				greaterOneCard = item
			}
		}
		if greaterOneCard != 0 && greaterTwoCard != 0 {
			cardGroup = append(cardGroup, fmt.Sprintf(`1%s2%s3%s`, curCardType, curCardType, curCardType))
		}
	case 2:
		for _, item := range cardInfo[curCardType] {
			if item == 1 {
				lessOneCard = item
			} else if item == 3 {
				greaterOneCard = item
			}
		}
		if lessOneCard != 0 && greaterOneCard != 0 {
			cardGroup = append(cardGroup, fmt.Sprintf(`1%s2%s3%s`, curCardType, curCardType, curCardType))
		}
		if lessOneCard != 0 && lessTwoCard != 0 {
			cardGroup = append(cardGroup, fmt.Sprintf(`2%s3%s4%s`, curCardType, curCardType, curCardType))
		}
	case 8:
		for _, item := range cardInfo[curCardType] {
			if item == 7 {
				lessOneCard = item
			} else if item == 9 {
				greaterOneCard = item
			}
		}
		if lessOneCard != 0 && greaterOneCard != 0 {
			cardGroup = append(cardGroup, fmt.Sprintf(`7%s8%s9%s`, curCardType, curCardType, curCardType))
		}
		if lessOneCard != 0 && lessTwoCard != 0 {
			cardGroup = append(cardGroup, fmt.Sprintf(`6%s7%s8%s`, curCardType, curCardType, curCardType))
		}
	case 9:
		for _, item := range cardInfo[curCardType] {
			if item == 7 {
				lessTwoCard = item
			} else if item == 8 {
				lessOneCard = item
			}
		}
		if lessOneCard != 0 && lessTwoCard != 0 {
			if lessOneCard != 0 && greaterOneCard != 0 {
				cardGroup = append(cardGroup, fmt.Sprintf(`7%s8%s9%s`, curCardType, curCardType, curCardType))
			}
		}
	default:
		for _, item := range cardInfo[curCardType] {
			if item == curCardNum-2 {
				lessTwoCard = item
			} else if item == curCardNum-1 {
				lessOneCard = item
			} else if item == curCardNum+1 {
				greaterOneCard = item
			} else if item == curCardNum+2 {
				greaterTwoCard = item
			}
		}
		if lessTwoCard != 0 && lessOneCard != 0 {
			str := fmt.Sprintf(`%d%s%d%s%d%s`, lessTwoCard, curCardType, lessOneCard, curCardType, curCardNum, curCardType)
			cardGroup = append(cardGroup, str)
		}
		if lessOneCard != 0 && greaterOneCard != 0 {
			str := fmt.Sprintf(`%d%s%d%s%d%s`, lessOneCard, curCardType, curCardNum, curCardType, greaterOneCard, curCardType)
			cardGroup = append(cardGroup, str)
		}
		if greaterOneCard != 0 && greaterTwoCard != 0 {
			str := fmt.Sprintf(`%d%s%d%s%d%s`, curCardNum, curCardType, greaterOneCard, curCardType, greaterTwoCard, curCardType)
			cardGroup = append(cardGroup, str)
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
	curCardType := curCard[1:]
	curCardNum, _ := strconv.Atoi(curCard[:1])

	cardInfo := GetPlayerCardInfo(roomNum, curPlayer)
	curCardTypeArr := cardInfo[curCardType]
	curCardTypeArr = append(curCardTypeArr, curCardNum)
	sort.Ints(curCardTypeArr)

	cardInfo[curCardType] = curCardTypeArr
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
