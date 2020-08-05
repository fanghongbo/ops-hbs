package main

import (
	"context"
	"github.com/fanghongbo/ops-hbs/common/g"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	g.InitAll()

	// 等待中断信号以优雅地关闭 Hbs（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Hbs ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := g.Shutdown(ctx); err != nil {
		log.Fatal("Hbs Shutdown:", err)
	} else {
		log.Println("Hbs Exiting")
	}
}
