package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CooperJiang/versiontrack-go-sdk/internal/utils"
	"github.com/CooperJiang/versiontrack-go-sdk/pkg/client"
)

var VERSION = "1.0.0"

func main() {
	var (
		checkUpdate = flag.Bool("check-update", false, "检查并执行更新")
		showVersion = flag.Bool("version", false, "显示版本信息")
		configFile  = flag.String("config", "config.yaml", "配置文件路径")
	)
	flag.Parse()

	// 显示版本信息
	if *showVersion {
		fmt.Printf("CLI Tool Version: %s\n", VERSION)
		fmt.Printf("Platform: %s\n", utils.GetPlatform())
		fmt.Printf("Arch: %s\n", utils.GetArch())
		return
	}

	// 检查更新
	if *checkUpdate {
		performUpdate(*configFile)
		return
	}

	// 正常业务逻辑
	fmt.Printf("CLI Tool v%s 正在运行...\n", VERSION)
	
	// 读取配置文件（如果存在）
	if utils.FileExists(*configFile) {
		fmt.Printf("使用配置文件: %s\n", *configFile)
		// 这里可以读取和处理配置文件
	} else {
		fmt.Printf("配置文件 %s 不存在，使用默认配置\n", *configFile)
	}

	// 执行主要业务逻辑
	runMainLogic()
}

func performUpdate(configFile string) {
	fmt.Println("🔍 开始检查更新...")

	// 🆕 配置更新客户端
	config := &client.Config{
		ServerURL:     "https://your-versiontrack-server.com",
		APIKey:        "your-api-key-here", // 🆕 使用API密钥替代ProjectID
		Platform:      utils.GetPlatform(),
		Arch:          utils.GetArch(),
		Timeout:       30 * time.Second,
		PreserveFiles: []string{"config.yaml", "config.yml", "*.conf", "data/*", "logs/*"},
		BackupCount:   3,
		UpdateMode:    client.UpdateModePrompt, // 🆕 提示模式
		SkipVersions:  []string{"1.0.2"},       // 🆕 跳过指定版本
	}

	updater, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("创建更新客户端失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 🆕 检查多版本更新
	updatesInfo, err := updater.CheckForMultipleUpdates(ctx, VERSION)
	if err != nil {
		log.Fatalf("检查更新失败: %v", err)
	}

	if !updatesInfo.HasUpdate {
		fmt.Println("✅ 当前已是最新版本")
		return
	}

	fmt.Printf("🎉 发现 %d 个可用更新版本:\n", len(updatesInfo.AvailableVersions))
	fmt.Printf("📋 当前版本: %s\n", updatesInfo.CurrentVersion)
	fmt.Printf("🚀 最新版本: %s\n\n", updatesInfo.LatestVersion)

	// 🆕 检查强制更新
	if updatesInfo.UpdateStrategy.HasForced {
		fmt.Printf("⚠️  检测到强制更新，最低要求版本: %s\n", updatesInfo.UpdateStrategy.MinRequiredVersion)
		
		forcedVersion, err := updater.HasForcedUpdate(ctx, VERSION)
		if err != nil {
			log.Printf("检查强制更新失败: %v", err)
		} else if forcedVersion != nil {
			fmt.Printf("🔴 必须更新到版本 %s 或更新版本\n\n", forcedVersion.Version)
		}
	}

	// 显示所有可用版本
	fmt.Println("📋 可用版本列表:")
	for i, version := range updatesInfo.AvailableVersions {
		status := ""
		if version.IsForced {
			status += " [强制更新]"
		}
		for _, skipVersion := range config.SkipVersions {
			if version.Version == skipVersion {
				status += " [已跳过]"
				break
			}
		}
		
		fmt.Printf("  %d. 版本 %s%s\n", i+1, version.Version, status)
		fmt.Printf("     📝 更新日志: %s\n", version.Changelog)
		fmt.Printf("     📅 发布日期: %s\n", version.ReleaseDate)
		fmt.Printf("     📦 文件大小: %s\n", formatBytes(version.FileSize))
		fmt.Println()
	}

	// 🆕 获取推荐更新版本
	recommendedVersion, err := updater.GetRecommendedUpdate(ctx, VERSION)
	if err != nil {
		log.Printf("获取推荐版本失败: %v", err)
	} else if recommendedVersion != nil {
		fmt.Printf("💡 推荐更新版本: %s\n\n", recommendedVersion.Version)
	}

	// 提示用户选择
	choice := promptUserChoice(updatesInfo.AvailableVersions, recommendedVersion)
	
	if choice == -1 {
		fmt.Println("👋 取消更新")
		return
	}

	selectedVersion := updatesInfo.AvailableVersions[choice]
	fmt.Printf("✅ 选择更新到版本: %s\n", selectedVersion.Version)

	// 确认更新
	if !confirmUpdate(selectedVersion) {
		fmt.Println("👋 取消更新")
		return
	}

	// 🆕 执行更新
	fmt.Printf("🚀 开始更新到版本 %s...\n", selectedVersion.Version)
	
	err = updater.UpdateToVersion(ctx, selectedVersion.Version, func(progress *client.DownloadProgress) {
		if progress.Total > 0 {
			fmt.Printf("\r📥 下载进度: %.1f%% (%s/%s)", 
				progress.Percentage, 
				formatBytes(progress.Downloaded), 
				formatBytes(progress.Total))
		}
	})
	
	if err != nil {
		log.Fatalf("更新失败: %v", err)
	}

	fmt.Printf("\n🎉 更新成功！已升级到版本: %s\n", selectedVersion.Version)
	fmt.Println("请重新运行程序以使用新版本")
	
	// 显示更新历史
	history := updater.GetUpdateHistory()
	if len(history) > 0 {
		fmt.Println("\n📚 最近的更新历史:")
		for i, record := range history {
			if i >= 3 { // 只显示最近3次更新
				break
			}
			fmt.Printf("  - %s: %s (%s)\n", 
				record.Version, 
				record.Status, 
				record.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}
}

func promptUserChoice(versions []client.VersionInfo, recommended *client.VersionInfo) int {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("请选择要更新的版本 (输入序号，0取消): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "0" {
			return -1
		}
		
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(versions) {
			fmt.Printf("❌ 无效选择，请输入 1-%d 或 0\n", len(versions))
			continue
		}
		
		return choice - 1
	}
}

func confirmUpdate(version client.VersionInfo) bool {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Printf("\n📋 更新详情:\n")
	fmt.Printf("  版本: %s\n", version.Version)
	fmt.Printf("  更新日志: %s\n", version.Changelog)
	fmt.Printf("  文件大小: %s\n", formatBytes(version.FileSize))
	if version.IsForced {
		fmt.Printf("  ⚠️  这是强制更新\n")
	}
	
	for {
		fmt.Print("\n确认更新? (y/n): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		
		switch input {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("❌ 请输入 y 或 n")
		}
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func runMainLogic() {
	// 这里是你的CLI工具的主要业务逻辑
	fmt.Println("执行业务逻辑...")
	
	// 模拟一些工作
	tasks := []string{
		"初始化配置",
		"连接数据库", 
		"处理数据",
		"生成报告",
		"清理临时文件",
	}

	for i, task := range tasks {
		fmt.Printf("[%d/%d] %s...\n", i+1, len(tasks), task)
		time.Sleep(500 * time.Millisecond) // 模拟工作耗时
	}

	fmt.Println("所有任务完成！")
}