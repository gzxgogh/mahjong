package controller

import (
	"fmt"
	"mahjong/model"
	"mahjong/service"
	"mahjong/utils"
	"strconv"
)

type action struct{}

var Action action

// Dice	godoc
// @Summary		摇骰子
// @Description	摇骰子
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Success 200 {string} string	"ok"
// @Router	/mahjong/dice [get]
func (ac *action) Dice(params map[string]string) model.Result {
	num := service.Action.Dice()
	return utils.Success(num)
}

// ShuffleCards	godoc
// @Summary		洗牌分牌
// @Description	洗牌分牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Param	diceNum formData int true "点数"
// @Param	player formData string true "当前的用户"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/shuffle/cards [post]
func (ac *action) ShuffleCards(params map[string]string) model.Result {
	fmt.Println("params", params)
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}
	diceNum, err := strconv.Atoi(params["diceNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：diceNum")
	}
	return service.Action.ShuffleCards(roomNum, diceNum, params["player"])
}

// GetGoldCard	godoc
// @Summary		获取该局的金
// @Description	获取该局的金
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/gold/get [get]
func (ac *action) GetGoldCard(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}
	return service.Action.GetGoldCard(roomNum)
}

// GrabOneCard	godoc
// @Summary		抓一张牌
// @Description	抓一张牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Param	player formData string true "当前的用户"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/grab/card [post]
func (ac *action) GrabOneCard(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}
	if params["player"] == "" {
		return utils.Error(-1, "无效的参数：player")
	}
	return service.Action.GrabOneCard(roomNum, params["player"])
}

// PlayOneCard	godoc
// @Summary		出一张手牌
// @Description	出一张手牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Param	player formData string true "当前的用户"
// @Param	curCard formData string true "牌{'type':'万','value':5}"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/play/card [post]
func (ac *action) PlayOneCard(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}
	if params["player"] == "" {
		return utils.Error(-1, "无效的参数：player")
	}
	var curdCard model.Card
	utils.FromJSON(params["curCurd"], &curdCard)
	return service.Action.PlayOneCard(roomNum, params["player"], curdCard)
}

// EatCard	godoc
// @Summary		吃牌
// @Description	吃牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Param	player formData string true "当前的用户"
// @Param	curCard formData string true "牌{'type':'万','value':5}"
// @Param	cardGroup formData string true "牌[{'type':'万','value':5},{'type':'万','value':6}，{'type':'万','value':7}]"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/eat/card [post]
func (ac *action) EatCard(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}
	if params["player"] == "" {
		return utils.Error(-1, "无效的参数：player")
	}
	var curdCard model.Card
	utils.FromJSON(params["curCurd"], &curdCard)

	var cardGroup []model.Card
	utils.FromJSON(params["cardGroup"], &cardGroup)

	return service.Action.EatCard(roomNum, curdCard, cardGroup, params["player"])
}

// TouchCard	godoc
// @Summary		碰牌
// @Description	碰牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Param	player formData string true "当前的用户"
// @Param	curCard formData string true "牌{'type':'万','value':5}"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/touch/card [post]
func (ac *action) TouchCard(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}
	if params["player"] == "" {
		return utils.Error(-1, "无效的参数：player")
	}
	var curCard model.Card
	utils.FromJSON(params["curCard"], &curCard)

	return service.Action.TouchCard(roomNum, curCard, params["player"])
}

// BarCard	godoc
// @Summary		杠牌
// @Description	杠牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Param	player formData string true "当前的用户"
// @Param	barType formData string true "rightBar/darkBar"
// @Param	curCard formData string true "牌{'type':'万','value':5}"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/bar/card [post]
func (ac *action) BarCard(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}
	if params["player"] == "" {
		return utils.Error(-1, "无效的参数：player")
	}
	if params["barType"] != "rightBar" && params["barType"] != "darkBar" {
		return utils.Error(-1, "无效的参数：barType")
	}
	var curdCard model.Card
	utils.FromJSON(params["curCurd"], &curdCard)

	return service.Action.BarCard(roomNum, curdCard, params["player"], params["barType"])
}

// RecordAbandonCard	godoc
// @Summary		记录弃牌堆
// @Description	记录弃牌堆
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Param	player formData string true "当前的用户"
// @Param	curCard formData string true "牌{'type':'万','value':5}"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/record/abandon/card [post]
func (ac *action) RecordAbandonCard(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}
	if params["player"] == "" {
		return utils.Error(-1, "无效的参数：player")
	}
	var curdCard model.Card
	utils.FromJSON(params["curCurd"], &curdCard)

	return service.Action.RecordAbandonCard(roomNum, curdCard, params["player"])
}

// GetAbandonCards	godoc
// @Summary		获取弃牌堆
// @Description	获取弃牌堆
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/abandon/cards/get [get]
func (ac *action) GetAbandonCards(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}

	return service.Action.GetAbandonCards(roomNum)
}

// GetPlayerCards	godoc
// @Summary		获取用户手牌
// @Description	获取用户手牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum formData int true "房间号"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/player/cards [get]
func (ac *action) GetPlayerCards(params map[string]string) model.Result {
	roomNum, err := strconv.Atoi(params["roomNum"])
	if err != nil {
		return utils.Error(-1, "无效的参数：roomNum")
	}

	return service.Action.GetPlayerCards(roomNum)
}
