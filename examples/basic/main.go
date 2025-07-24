package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lilithgames/versiontrack-go-sdk/internal/utils"
	"github.com/lilithgames/versiontrack-go-sdk/pkg/client"
)

func main() {
	// 配置更新客户端
	config := &client.Config{
		ServerURL:     "https://your-versiontrack-server.com",
		ProjectID:     "your-project-id",
		Platform:      utils.GetPlatform(), // 自动检测平台
		Arch:          utils.GetArch(),     // 自动检测架构
		Timeout:       30 * time.Second,
		PreserveFiles: []string{"config.yaml", "config.yml", "*.conf", "data.db"},
		BackupCount:   3,
	}

	// 创建客户端
	updater, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 当前版本（通常从你的应用程序中获取）
	currentVersion := "1.0.0"

	fmt.Printf("当前版本: %s\n", currentVersion)
	fmt.Printf("平台: %s\n", config.Platform)
	fmt.Printf("架构: %s\n", config.Arch)

	// 检查更新
	fmt.Println("检查更新中...")
	ctx := context.Background()
	updateInfo, err := updater.CheckForUpdates(ctx, currentVersion)
	if err != nil {
		log.Fatalf("Failed to check for updates: %v", err)
	}

	if !updateInfo.HasUpdate {
		fmt.Println("当前已是最新版本")
		return
	}

	fmt.Printf("发现新版本: %s\n", updateInfo.LatestVersion)
	fmt.Printf("文件大小: %d bytes\n", updateInfo.FileSize)
	fmt.Printf("发布说明: %s\n", updateInfo.ReleaseNotes)

	// 下载更新
	fmt.Println("开始下载更新...")
	downloadPath := fmt.Sprintf("/tmp/update_%s.tar.gz", updateInfo.LatestVersion)
	
	err = updater.Download(ctx, updateInfo, downloadPath, func(progress *client.DownloadProgress) {
		if progress.Total > 0 {
			fmt.Printf("\r下载进度: %.1f%% (%d/%d bytes)", progress.Percentage, progress.Downloaded, progress.Total)
		}
	})
	
	if err != nil {
		log.Fatalf("Failed to download update: %v", err)
	}
	fmt.Println("\n下载完成")

	// 执行更新
	fmt.Println("开始执行更新...")
	err = updater.Update(ctx, updateInfo, downloadPath)
	if err != nil {
		log.Fatalf("Failed to update: %v", err)
	}

	fmt.Printf("更新成功，已升级到版本: %s\n", updateInfo.LatestVersion)
	
	// 显示更新历史
	history := updater.GetUpdateHistory()
	if len(history) > 0 {
		fmt.Println("\n更新历史:")
		for _, record := range history {
			fmt.Printf("- %s: %s (%s)\n", record.Version, record.Status, record.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}
}