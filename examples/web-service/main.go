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

	"github.com/lilithgames/versiontrack-go-sdk/internal/utils"
	"github.com/lilithgames/versiontrack-go-sdk/pkg/client"
)

var (
	VERSION = "1.0.0"
	server  *http.Server
)

func main() {
	// 启动Web服务
	go startWebServer()

	// 启动更新检查器
	go startUpdateChecker()

	// 等待信号
	waitForSignal()
}

func startWebServer() {
	mux := http.NewServeMux()
	
	// 版本信息接口
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"version": "%s", "status": "running"}`, VERSION)
	})

	// 健康检查接口
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok"}`)
	})

	// 手动更新接口
	mux.HandleFunc("/update", handleManualUpdate)

	server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Printf("Web服务已启动，端口: 8080，版本: %s\n", VERSION)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Web服务启动失败: %v", err)
	}
}

func startUpdateChecker() {
	// 配置更新客户端
	config := &client.Config{
		ServerURL:     "https://your-versiontrack-server.com",
		ProjectID:     "your-web-service-project-id",
		Platform:      utils.GetPlatform(),
		Arch:          utils.GetArch(),
		Timeout:       30 * time.Second,
		PreserveFiles: []string{"config.yaml", "config.yml", "*.conf", "data.db", "logs/*"},
		BackupCount:   5,
	}

	updater, err := client.NewClient(config)
	if err != nil {
		log.Printf("Failed to create update client: %v", err)
		return
	}

	// 定时检查更新（每30分钟）
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkAndUpdate(updater)
		}
	}
}

func checkAndUpdate(updater *client.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("检查更新中...")
	
	updateInfo, err := updater.CheckForUpdates(ctx, VERSION)
	if err != nil {
		log.Printf("检查更新失败: %v", err)
		return
	}

	if !updateInfo.HasUpdate {
		log.Println("当前已是最新版本")
		return
	}

	log.Printf("发现新版本: %s，准备更新", updateInfo.LatestVersion)

	// 下载更新
	downloadPath := fmt.Sprintf("/tmp/web_service_update_%s.tar.gz", updateInfo.LatestVersion)
	err = updater.Download(ctx, updateInfo, downloadPath, func(progress *client.DownloadProgress) {
		if progress.Total > 0 {
			log.Printf("下载进度: %.1f%%", progress.Percentage)
		}
	})

	if err != nil {
		log.Printf("下载更新失败: %v", err)
		return
	}

	log.Println("开始执行更新...")

	// 优雅关闭服务器
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("服务器关闭失败: %v", err)
		return
	}

	// 执行更新
	err = updater.Update(ctx, updateInfo, downloadPath)
	if err != nil {
		log.Printf("更新失败: %v", err)
		// 重新启动服务器
		go startWebServer()
		return
	}

	log.Printf("更新成功，版本: %s", updateInfo.LatestVersion)
	
	// 更新版本号
	VERSION = updateInfo.LatestVersion

	// 重新启动服务器
	go startWebServer()
}

func handleManualUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 配置更新客户端
	config := &client.Config{
		ServerURL:     "https://your-versiontrack-server.com",
		ProjectID:     "your-web-service-project-id",
		Platform:      utils.GetPlatform(),
		Arch:          utils.GetArch(),
		Timeout:       30 * time.Second,
		PreserveFiles: []string{"config.yaml", "config.yml", "*.conf", "data.db", "logs/*"},
		BackupCount:   5,
	}

	updater, err := client.NewClient(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create client: %v", err), http.StatusInternalServerError)
		return
	}

	// 在后台执行更新
	go func() {
		checkAndUpdate(updater)
	}()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Update started", "status": "ok"}`)
}

func waitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\n收到退出信号，正在关闭服务...")

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if server != nil {
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("服务器关闭失败: %v", err)
		}
	}

	fmt.Println("服务已关闭")
}