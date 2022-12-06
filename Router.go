package main

import (
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"mahjong/controller"
	_ "mahjong/docs"
	"mahjong/utils"
	"net/http"
)

func setupRouter() *gin.Engine {
	engine := gin.Default()
	//添加swagger支持
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//处理全局异常
	engine.Use(nice.Recovery(recoveryHandler))
	//设置404返回的内容
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, utils.Error(404, "无效的路由"))
	})

	engine.GET("/mahjong/dice", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.Dice(GinParamMap(c)))
	})
	engine.POST("/mahjong/shuffle/cards", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.ShuffleCards(GinParamMap(c)))
	})
	engine.GET("/mahjong/gold/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.GetGoldCard(GinParamMap(c)))
	})
	engine.POST("/mahjong/grab/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.GrabOneCard(GinParamMap(c)))
	})
	engine.POST("/mahjong/play/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.PlayOneCard(GinParamMap(c)))
	})
	engine.POST("/mahjong/eat/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.EatCard(GinParamMap(c)))
	})
	engine.POST("/mahjong/touch/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.TouchCard(GinParamMap(c)))
	})
	engine.POST("/mahjong/bar/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.BarCard(GinParamMap(c)))
	})
	engine.POST("/mahjong/record/abandon/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.RecordAbandonCard(GinParamMap(c)))
	})
	engine.GET("/mahjong/abandon/cards", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.GetAbandonCards(GinParamMap(c)))
	})

	return engine
}

func GinParamMap(c *gin.Context) map[string]string {
	params := make(map[string]string)
	if c.Request.Method == "GET" {
		for k, v := range c.Request.URL.Query() {
			params[k] = v[0]
		}
		return params
	} else {
		c.Request.ParseForm()
		for k, v := range c.Request.PostForm {
			params[k] = v[0]
		}
		for k, v := range c.Request.URL.Query() {
			params[k] = v[0]
		}
		return params
	}
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"msg":  "系统异常，请联系客服",
		"code": 1001,
	})
}
