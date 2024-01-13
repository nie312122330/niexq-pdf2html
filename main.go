package main

import (
	"fmt"
	"niexq-html2pdf/config"
	"niexq-html2pdf/controllers"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/nie312122330/niexq-gotools/logext"
	"github.com/nie312122330/niexq-gowebapi/ginext"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger = logext.DefaultLogger(config.AppConf.Server.AppName)
}

func main() {
	if config.AppConf.Server.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	//创建PDF结果目录
	os.MkdirAll(config.AppConf.ChromeConf.Pdfdir, 0777)

	crontab := newWithSeconds()
	defer crontab.Stop()

	// 每秒钟显示当前时间
	// crontab.AddFunc("*/1 * * * * ?", func() {
	// 	dateStr, _ := dateext.Now().Format("yyyy-MM-dd HH:mm:ss", true)
	// 	fmt.Printf("%s 程序当前使用内存量:%s\n", dateStr, curMemStats())
	// })
	// // 启动定时器
	// crontab.Start()

	ginEngine := gin.New()
	//PDF结果目录
	ginEngine.Static("/"+config.AppConf.ChromeConf.Pdfdir, config.AppConf.ChromeConf.Pdfdir)

	//设置自定义日志,错误恢复中间件
	ginEngine.Use(ginext.GinLogger(), ginext.GinRecovery())
	//注册路由
	controllers.PubCtrRegisterRouter(ginEngine, config.AppConf.Server.AppName)
	//运行
	logger.Info(fmt.Sprintf("开始启动web服务,端口:%d", config.AppConf.Server.ServerPort))

	ginEngine.Run(fmt.Sprintf(":%d", config.AppConf.Server.ServerPort))
}

// 当前的内存状态
func curMemStats() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	var maloc string
	if m.Alloc > 1024*1024 {
		maloc = fmt.Sprintf("%0.2fMb", float32(m.Alloc)/1024/1024)
	} else {
		maloc = fmt.Sprintf("%f0.2Kb", float32(m.Alloc)/1024)
	}
	return fmt.Sprintf(" 内存:%s", maloc)
}

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
