package main

import (
	"app/config"
	"app/initialize"
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed assets/*
var embededFiles embed.FS

func main() {
	e := gin.New()
	initialize.Init(e, embededFiles)

	if gin.Mode() == gin.ReleaseMode {
		// 生产模块下实现 gracefully shutdown
		srv := &http.Server{
			Addr: ":" + config.Config.Server.Port,
			Handler: e,
		}
		go func() {
			// service connections
			log.Printf("start to listen %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen server error: %s\n", err)
			}
		}()

	// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscanll.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		// catching ctx.Done(). timeout of 5 seconds.
		select {
		case <-ctx.Done():
			log.Println("timeout of 5 seconds.")
		}
		log.Println("Server exiting")
	} else {
		// 开发模式下使用 air 监听变化自动重启
		e.Run(":" + config.Config.Server.Port)
	}
}
