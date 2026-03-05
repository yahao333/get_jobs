package main

import (
	"fmt"
	"os"

	"github.com/yahao333/get_jobs/internal/config"
	"github.com/yahao333/get_jobs/internal/storage"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("config.yaml"); err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	config.Info("启动 get_jobs 应用...")

	// 初始化数据库
	if err := storage.InitDB(); err != nil {
		config.Error("初始化数据库失败: ", err)
		os.Exit(1)
	}

	config.Info("应用启动成功!")
	config.Info("Web 服务地址: http://", config.GetString("web.host"), ":", config.GetInt("web.port"))

	// TODO: 启动 Web 服务
	// TODO: 启动求职任务

	// 保持运行
	select {}
}
