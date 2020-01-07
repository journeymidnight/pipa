package main

import (
	"github.com/journeymidnight/pipa/handler"
	"os"
	"os/signal"
	"syscall"

	"github.com/journeymidnight/pipa/helper"
	"github.com/journeymidnight/pipa/redis"
)

func main() {

	helper.SetupGlobalConfig()
	logLevel := helper.ParseLevel(helper.Config.LogLevel)
	helper.Log = helper.NewFileLogger(helper.Config.LogPath, logLevel)
	defer helper.Log.Close()

	helper.Log.Info("Pipa start!")
	redis.Initialize()
	defer redis.Close()

	handler.StartWorker()

	signal.Ignore()
	signalQueue := make(chan os.Signal)
	signal.Notify(signalQueue, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGUSR1)
	for {
		s := <-signalQueue
		switch s {
		case syscall.SIGHUP:
			helper.SetupGlobalConfig()
		case syscall.SIGUSR1:
			// TODO: Dump something?
		default:
			// TODO: Stop pipa server with graceful shutdown
			handler.Stop()
			return
		}
	}
}
