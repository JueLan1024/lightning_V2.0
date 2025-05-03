package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/kafka"

	_ "web_app/docs" // 导入生成的 Swagger 文档
	"web_app/logger"
	"web_app/pkg/bloom"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"

	"go.uber.org/zap"
	// _ "net/http/pprof" // 导入 pprof 包
)

func main() {
	// 1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("settings.Init() failed,err:%v\n", err)
		return
	}
	// 2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("logger.Init() failed,err:%v\n", err)
		return
	}
	zap.L().Debug("logger init success...")
	// 3.初始化MySQL连接
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		zap.L().Error("mysql.Init() failed", zap.Error(err))
		return
	}
	defer mysql.Close()
	// 4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		zap.L().Error("redis.Init() failed", zap.Error(err))
		return
	}
	defer redis.Close()
	// 5.初始化雪花ID生成器
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		zap.L().Error("snowflake.Init() failed", zap.Error(err))
		return
	}
	// 6.初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		zap.L().Error("controller.InitTrans failed", zap.Error(err))
		return
	}
	// 7.初始化布隆过滤器
	if err := bloom.InitBloomFilter(); err != nil {
		zap.L().Error("bloom.InitBloomFilter() failed", zap.Error(err))
		return
	}
	// 背景context
	ctx, cancel := context.WithCancel(context.Background())
	// 8.初始化kafka.Reader
	kafka.Init(ctx, settings.Conf.KafkaConfig)
	// 注册路由
	r := routes.Setup(settings.Conf.Mode, settings.Conf.RatelimitConfig)
	// 启动服务(优雅关机)
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			settings.Conf.Host,
			settings.Conf.Port),
		Handler: r,
	}

	// // 启动 pprof 服务
	// go func() {
	// 	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
	// 		zap.L().Error("pprof server failed to start", zap.Error(err))
	// 	}
	// }()

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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")

	// 通知 Kafka 消费者停止
	cancel()
	// 创建一个5秒超时的context
	ctx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeoutCancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
