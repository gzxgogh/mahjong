package main

import (
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/gzxgogh/ggin/models"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"mahjong/controller"
	_ "mahjong/docs"
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
		c.JSON(http.StatusOK, models.Error(404, "无效的路由"))
	})

	engine.GET("/mahjong/dice", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.Dice(c))
	})

	engine.POST("/mahjong/shuffle/cards", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.ShuffleCards(c))
	})

	engine.GET("/mahjong/gold/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.GetGoldCard(c))
	})

	engine.POST("/mahjong/grab/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.GrabOneCard(c))
	})

	engine.POST("/mahjong/play/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.PlayOneCard(c))
	})

	engine.POST("/mahjong/eat/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.EatCard(c))
	})

	engine.POST("/mahjong/touch/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.TouchCard(c))
	})

	engine.POST("/mahjong/bar/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.BarCard(c))
	})

	engine.POST("/mahjong/record/abandon/card", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.RecordAbandonCard(c))
	})

	engine.GET("/mahjong/abandon/cards", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.GetAbandonCards(c))
	})

	engine.GET("/mahjong/player/cards", func(c *gin.Context) {
		c.JSON(http.StatusOK, controller.Action.GetPlayerCards(c))
	})

	return engine
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"msg":  "系统异常，请联系客服",
		"code": 1001,
	})
}
