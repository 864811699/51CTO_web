package main

import (
	"51ctoGoWeb/dao/mysql"
	"51ctoGoWeb/dao/redis"
	"51ctoGoWeb/logger"
	"51ctoGoWeb/routers"
	"51ctoGoWeb/settings"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//加载配置
	if err:=settings.Init();err != nil{
		panic(fmt.Errorf("settings.Init fail : %v",err))
	}

	//初始化日志
	if err:=logger.Init();err!=nil {
		panic(fmt.Errorf("logger init fail : %v",err))
	}
	defer zap.L().Sync()
	zap.L().Info("logger init success .. ")

	//初始化mysql
	if err:=mysql.Init();err!=nil {
		panic(fmt.Errorf("mysql init fail : %v",err))
	}
	zap.L().Info("mysql init success .. ")
	defer mysql.Close()

	//初始化redis
	if err:=redis.Init();err!=nil {
		panic(fmt.Errorf("redis init fail : %v",err))
	}
	zap.L().Info("redis init success .. ")
	defer redis.Close()

	//注册路由
	engine:=routers.Setup()
	//启动服务

	srv := &http.Server{
		Addr:   fmt.Sprintf(":%d",viper.GetInt("app.port")),
		Handler: engine,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)  // 此处不会阻塞
	<-quit  // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")

}