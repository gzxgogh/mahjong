package main

import (
	"context"
	"log"
	"mahjong/redis"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	redis.InItRedisCoon()
	engine := setupRouter()

	server := &http.Server{
		Addr:    ":8011",
		Handler: engine,
	}

	go func() {
		var err error
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("HTTP server listen: {}", err.Error())
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	log.Println("Get Signal:" + sig.String())
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:" + err.Error())
	}
	log.Println("Server exiting")
}
