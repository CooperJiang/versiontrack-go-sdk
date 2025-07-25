package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CooperJiang/versiontrack-go-sdk/internal/utils"
	"github.com/CooperJiang/versiontrack-go-sdk/pkg/client"
)

func main() {
	// 配置更新客户端
	config := &client.Config{
		ServerURL:     "http://localhost:9000",
		APIKey:        "your-api-key-here",                      // 使用API密钥替代ProjectID
		Platform:      utils.GetPlatform(),                     // 自动检测平台
		Arch:          utils.GetArch(),                          // 自动检测架构
		Timeout:       30 * time.Second,
		PreserveFiles: []string{"config.yaml", "config.yml", "*.conf", "data.db"},
		BackupCount:   3,
		UpdateMode:    client.UpdateModeAuto,                    // 设置更新模式
		SkipVersions:  []string{},                               // 跳过的版本列表
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

	// 检查多版本更新
	fmt.Println("检查多版本更新中...")
	ctx := context.Background()
	updatesInfo, err := updater.CheckForMultipleUpdates(ctx, currentVersion)
	if err != nil {
		log.Fatalf("Failed to check for updates: %v", err)
	}

	if !updatesInfo.HasUpdate {
		fmt.Println("当前已是最新版本")
		return
	}

	fmt.Printf("发现 %d 个可用更新版本，最新版本: %s\n", len(updatesInfo.AvailableVersions), updatesInfo.LatestVersion)
	
	// 显示所有可用版本
	for i, version := range updatesInfo.AvailableVersions {
		fmt.Printf("%d. 版本 %s - %s", i+1, version.Version, version.Changelog)
		if version.IsForced {
			fmt.Print(" [强制更新]")
		}
		fmt.Println()
	}

	// 检查强制更新
	if updatesInfo.UpdateStrategy.HasForced {
		fmt.Printf("检测到强制更新，最低要求版本: %s\n", updatesInfo.UpdateStrategy.MinRequiredVersion)
	}

	// 根据更新模式执行不同操作
	switch config.UpdateMode {
	case client.UpdateModeAuto:
		// 自动更新到推荐版本
		recommendedVersion, err := updater.GetRecommendedUpdate(ctx, currentVersion)
		if err != nil {
			log.Fatalf("Failed to get recommended update: %v", err)
		}
		
		if recommendedVersion != nil {
			fmt.Printf("自动更新到推荐版本: %s\n", recommendedVersion.Version)
			err = updater.UpdateToVersion(ctx, recommendedVersion.Version, func(progress *client.DownloadProgress) {
				if progress.Total > 0 {
					fmt.Printf("\r下载进度: %.1f%% (%d/%d bytes)", progress.Percentage, progress.Downloaded, progress.Total)
				}
			})
			if err != nil {
				log.Fatalf("Failed to update: %v", err)
			}
			fmt.Printf("\n自动更新成功，已升级到版本: %s\n", recommendedVersion.Version)
		}
		
	case client.UpdateModeManual:
		// 手动选择版本（这里选择最新版本作为示例）
		if len(updatesInfo.AvailableVersions) > 0 {
			targetVersion := updatesInfo.AvailableVersions[0].Version
			fmt.Printf("手动选择更新到版本: %s\n", targetVersion)
			
			err = updater.UpdateToVersion(ctx, targetVersion, func(progress *client.DownloadProgress) {
				if progress.Total > 0 {
					fmt.Printf("\r下载进度: %.1f%% (%d/%d bytes)", progress.Percentage, progress.Downloaded, progress.Total)
				}
			})
			if err != nil {
				log.Fatalf("Failed to update: %v", err)
			}
			fmt.Printf("\n手动更新成功，已升级到版本: %s\n", targetVersion)
		}
		
	default:
		fmt.Println("提示模式：请手动选择要更新的版本")
	}

	// 兼容旧版API的示例
	fmt.Println("\n=== 使用旧版API检查更新 ===")
	updateInfo, err := updater.CheckForUpdates(ctx, currentVersion)
	if err != nil {
		log.Fatalf("Failed to check for updates (legacy): %v", err)
	}

	if updateInfo.HasUpdate && updateInfo.LatestVersion != nil {
		fmt.Printf("旧版API检测到更新: %s -> %s\n", currentVersion, updateInfo.LatestVersion.Version)
	}
	
	// 显示更新历史
	history := updater.GetUpdateHistory()
	if len(history) > 0 {
		fmt.Println("\n更新历史:")
		for _, record := range history {
			fmt.Printf("- %s: %s (%s)\n", record.Version, record.Status, record.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}
}