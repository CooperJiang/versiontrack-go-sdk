package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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
	fmt.Println("开始检查更新...")

	// 配置更新客户端
	config := &client.Config{
		ServerURL:     "https://your-versiontrack-server.com",
		ProjectID:     "your-cli-tool-project-id", 
		Platform:      utils.GetPlatform(),
		Arch:          utils.GetArch(),
		Timeout:       30 * time.Second,
		PreserveFiles: []string{"config.yaml", "config.yml", "*.conf", "data/*", "logs/*"},
		BackupCount:   3,
	}

	updater, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("创建更新客户端失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 检查更新
	updateInfo, err := updater.CheckForUpdates(ctx, VERSION)
	if err != nil {
		log.Fatalf("检查更新失败: %v", err)
	}

	if !updateInfo.HasUpdate {
		fmt.Println("当前已是最新版本")
		return
	}

	fmt.Printf("发现新版本: %s\n", updateInfo.LatestVersion)
	fmt.Printf("文件大小: %d bytes\n", updateInfo.FileSize)
	
	if updateInfo.ReleaseNotes != "" {
		fmt.Printf("更新说明:\n%s\n", updateInfo.ReleaseNotes)
	}

	// 询问用户是否要更新
	fmt.Print("是否要更新到新版本? (y/n): ")
	var response string
	fmt.Scanln(&response)
	
	if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
		fmt.Println("更新已取消")
		return
	}

	// 下载更新
	fmt.Println("正在下载更新...")
	downloadPath := fmt.Sprintf("/tmp/cli_tool_update_%s.tar.gz", updateInfo.LatestVersion)
	
	err = updater.Download(ctx, updateInfo, downloadPath, func(progress *client.DownloadProgress) {
		if progress.Total > 0 {
			fmt.Printf("\r下载进度: %.1f%% (%d/%d bytes)", progress.Percentage, progress.Downloaded, progress.Total)
		}
	})
	
	if err != nil {
		log.Fatalf("下载更新失败: %v", err)
	}
	fmt.Println("\n下载完成")

	// 执行更新
	fmt.Println("正在执行更新...")
	err = updater.Update(ctx, updateInfo, downloadPath)
	if err != nil {
		log.Fatalf("更新失败: %v", err)
	}

	fmt.Printf("更新成功！已升级到版本: %s\n", updateInfo.LatestVersion)
	fmt.Println("请重新运行程序以使用新版本")
	
	// 显示更新历史
	history := updater.GetUpdateHistory()
	if len(history) > 0 {
		fmt.Println("\n最近的更新历史:")
		for i, record := range history {
			if i >= 3 { // 只显示最近3次更新
				break
			}
			fmt.Printf("- %s: %s (%s)\n", 
				record.Version, 
				record.Status, 
				record.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	// 清理下载的临时文件
	os.Remove(downloadPath)
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