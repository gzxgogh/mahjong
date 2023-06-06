package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gzxgogh/ggin/models"
	"mahjong/model"
	"mahjong/service"
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
func (ac *action) Dice(c *gin.Context) models.Result {
	num := service.Action.Dice()
	return models.Success(num)
}

// ShuffleCards	godoc
// @Summary		洗牌分牌
// @Description	洗牌分牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	param body models.ShuffleCardsReq true "请求体"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/shuffle/cards [post]
func (ac *action) ShuffleCards(c *gin.Context) models.Result {
	var params model.ShuffleCardsReq
	if err := c.ShouldBindJSON(&params); err != nil {
		return models.Error(-1, "解析失败")
	}

	return service.Action.ShuffleCards(params.RoomNum, params.DiceNum, params.Player)
}

// GetGoldCard	godoc
// @Summary		获取该局的金
// @Description	获取该局的金
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum query int true "房间号"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/gold/get [get]
func (ac *action) GetGoldCard(c *gin.Context) models.Result {
	roomNum, err := strconv.Atoi(c.Query("roomNum"))
	if err != nil {
		return models.Error(-1, "无效的参数：roomNum")
	}
	return service.Action.GetGoldCard(roomNum)
}

// GrabOneCard	godoc
// @Summary		抓一张牌
// @Description	抓一张牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	param body models.GrabOneCardReq true "请求体"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/grab/card [post]
func (ac *action) GrabOneCard(c *gin.Context) models.Result {
	var params model.GrabOneCardReq
	if err := c.ShouldBindJSON(&params); err != nil {
		return models.Error(-1, "解析失败")
	}
	return service.Action.GrabOneCard(params.RoomNum, params.Player)
}

// PlayOneCard	godoc
// @Summary		出一张手牌
// @Description	出一张手牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	param body models.PlayOneCardReq true "请求体"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/play/card [post]
func (ac *action) PlayOneCard(c *gin.Context) models.Result {
	var params model.PlayOneCardReq
	if err := c.ShouldBindJSON(&params); err != nil {
		return models.Error(-1, "解析失败")
	}
	return service.Action.PlayOneCard(params.RoomNum, params.Player, params.CurCard)
}

// EatCard	godoc
// @Summary		吃牌
// @Description	吃牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	param body models.EatCardReq true "请求体"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/eat/card [post]
func (ac *action) EatCard(c *gin.Context) models.Result {
	var params model.EatCardReq
	if err := c.ShouldBindJSON(&params); err != nil {
		return models.Error(-1, "解析失败")
	}

	return service.Action.EatCard(params.RoomNum, params.Player, params.CurCard, params.CardGroup)
}

// TouchCard	godoc
// @Summary		碰牌
// @Description	碰牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	param body models.TouchCardReq true "请求体"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/touch/card [post]
func (ac *action) TouchCard(c *gin.Context) models.Result {
	var params model.TouchCardReq
	if err := c.ShouldBindJSON(&params); err != nil {
		return models.Error(-1, "解析失败")
	}

	return service.Action.TouchCard(params.RoomNum, params.Player, params.CurCard)
}

// BarCard	godoc
// @Summary		杠牌
// @Description	杠牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	param body models.BarCardReq true "请求体"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/bar/card [post]
func (ac *action) BarCard(c *gin.Context) models.Result {
	var params model.BarCardReq
	if err := c.ShouldBindJSON(&params); err != nil {
		return models.Error(-1, "解析失败")
	}

	return service.Action.BarCard(params.RoomNum, params.Player, params.BarType, params.CurCard)
}

// RecordAbandonCard	godoc
// @Summary		记录弃牌堆
// @Description	记录弃牌堆
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	param body models.PlayOneCardReq true "请求体"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/record/abandon/card [post]
func (ac *action) RecordAbandonCard(c *gin.Context) models.Result {
	var params model.PlayOneCardReq
	if err := c.ShouldBindJSON(&params); err != nil {
		return models.Error(-1, "解析失败")
	}
	return service.Action.RecordAbandonCard(params.RoomNum, params.Player, params.CurCard)
}

// GetAbandonCards	godoc
// @Summary		获取弃牌堆
// @Description	获取弃牌堆
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum query int true "房间号"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/abandon/cards/get [get]
func (ac *action) GetAbandonCards(c *gin.Context) models.Result {
	roomNum, err := strconv.Atoi(c.Query("roomNum"))
	if err != nil {
		return models.Error(-1, "无效的参数：roomNum")
	}

	return service.Action.GetAbandonCards(roomNum)
}

// GetPlayerCards	godoc
// @Summary		获取用户手牌
// @Description	获取用户手牌
// @Tags         麻将
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	roomNum query int true "房间号"
// @Success 200 {string} string	"ok"
// @Router	/mahjong/player/cards [get]
func (ac *action) GetPlayerCards(c *gin.Context) models.Result {
	roomNum, err := strconv.Atoi(c.Query("roomNum"))
	if err != nil {
		return models.Error(-1, "无效的参数：roomNum")
	}

	return service.Action.GetPlayerCards(roomNum)
}
