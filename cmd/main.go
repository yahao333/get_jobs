package main

import (
	"fmt"
	"os"

	"github.com/yahao333/get_jobs/internal/config"
	"github.com/yahao333/get_jobs/internal/storage"
	"github.com/yahao333/get_jobs/internal/web"
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

	// 启动 Web 服务
	port := config.GetInt("web.port")
	if port == 0 {
		port = 8080
	}

	config.Info("Web 服务地址: http://localhost:", port)

	// 启动 Web 服务器
	server := web.NewServer(port)
	if err := server.Start(); err != nil {
		config.Error("Web 服务启动失败: ", err)
		os.Exit(1)
	}

	// 保持运行
	select {}
}
