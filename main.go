package main

import (
	"context"
	"flag"
	"github.com/gzxgogh/ggin"
	"github.com/gzxgogh/ggin/logs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func parseArgs() string {
	var configFile string
	flag.StringVar(&configFile, "f", os.Args[0]+".yml", "yml配置文件名")
	flag.Parse()
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if !strings.Contains(configFile, "/") {
		configFile = path + "/" + configFile
	}
	return configFile
}

func main() {
	cfgFile := parseArgs()
	ggin.Init(cfgFile)

	engine := setupRouter()
	server := &http.Server{
		Addr:    ":8011",
		Handler: engine,
	}

	go func() {
		var err error
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logs.Error("HTTP server listen: {}", err.Error())
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	logs.Info("Get Signal:" + sig.String())
	logs.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logs.Error("Server Shutdown:" + err.Error())
	}
	logs.Info("Server exiting")
}
