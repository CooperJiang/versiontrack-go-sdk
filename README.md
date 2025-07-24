# VersionTrack Go SDK

一个用于Go应用程序版本管理和热更新的SDK，与VersionTrack系统集成。

## 特性

- 🚀 **自动版本检查**：定期检查是否有新版本发布
- 📦 **智能更新**：支持tar.gz包的解压和安装
- 🔒 **配置保护**：更新时自动保护重要配置文件
- 📊 **进度跟踪**：实时显示下载和更新进度
- 🔄 **回滚支持**：更新失败时自动回滚
- 📝 **更新历史**：记录所有更新操作的历史
- 🛡️ **安全验证**：MD5校验确保文件完整性

## 快速开始

### 安装

```bash
go get github.com/lilithgames/versiontrack-go-sdk
```

### 基本使用

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/lilithgames/versiontrack-go-sdk/internal/utils"
    "github.com/lilithgames/versiontrack-go-sdk/pkg/client"
)

func main() {
    // 配置客户端
    config := &client.Config{
        ServerURL:     "https://your-versiontrack-server.com",
        ProjectID:     "your-project-id",
        Platform:      utils.GetPlatform(), // 自动检测平台
        Arch:          utils.GetArch(),     // 自动检测架构
        Timeout:       30 * time.Second,
        PreserveFiles: []string{"config.yaml", "*.conf"},
        BackupCount:   3,
    }

    // 创建客户端
    updater, err := client.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // 检查更新
    ctx := context.Background()
    updateInfo, err := updater.CheckForUpdates(ctx, "1.0.0")
    if err != nil {
        log.Fatal(err)
    }

    if updateInfo.HasUpdate {
        // 执行更新...
    }
}
```

## 配置说明

### Config 结构

```go
type Config struct {
    ServerURL     string        // VersionTrack服务器地址
    ProjectID     string        // 项目ID
    Platform      string        // 平台 (windows/linux/macos)
    Arch          string        // 架构 (amd64/arm64)
    Timeout       time.Duration // HTTP请求超时时间
    PreserveFiles []string      // 需要保护的文件列表
    BackupCount   int          // 备份保留数量
}
```

### 参数说明

- **ServerURL**: VersionTrack服务器的API地址
- **ProjectID**: 在VersionTrack系统中创建的项目ID
- **Platform**: 目标平台，支持 `windows`、`linux`、`macos`
- **Arch**: 目标架构，支持 `amd64`、`arm64`
- **Timeout**: HTTP请求超时时间，默认30秒
- **PreserveFiles**: 更新时不覆盖的文件模式列表，默认包含 `config.yaml`
- **BackupCount**: 保留的备份数量，默认3个

## 主要接口

### Updater 接口

```go
type Updater interface {
    // 检查是否有可用更新
    CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error)
    
    // 下载更新文件
    Download(ctx context.Context, info *UpdateInfo, destPath string, callback ProgressCallback) error
    
    // 执行更新
    Update(ctx context.Context, info *UpdateInfo, downloadPath string) error
    
    // 获取更新历史
    GetUpdateHistory() []UpdateRecord
    
    // 回滚到指定版本
    Rollback(ctx context.Context, version string) error
}
```

### 更新信息结构

```go
type UpdateInfo struct {
    HasUpdate     bool   `json:"hasUpdate"`     // 是否有更新
    LatestVersion string `json:"latestVersion"` // 最新版本号
    DownloadURL   string `json:"downloadUrl"`   // 下载地址
    FileSize      int64  `json:"fileSize"`      // 文件大小
    MD5Hash       string `json:"md5Hash"`       // MD5校验值
    ReleaseNotes  string `json:"releaseNotes"`  // 发布说明
    PublishedAt   string `json:"publishedAt"`   // 发布时间
}
```

## 使用示例

### 1. 基础示例

参见 [examples/basic/main.go](examples/basic/main.go)

### 2. Web服务示例

参见 [examples/web-service/main.go](examples/web-service/main.go)

这个示例展示了如何在Web服务中集成自动更新功能，包括：
- 定时检查更新
- 优雅关闭服务
- 手动触发更新的API接口

### 3. CLI工具示例

参见 [examples/cli-tool/main.go](examples/cli-tool/main.go)

这个示例展示了如何为命令行工具添加更新功能：
- 命令行参数控制
- 用户交互确认
- 配置文件保护

## 更新包结构

SDK支持包含以下文件的tar.gz更新包：

```
update-package.tar.gz
├── binary-file          # 主程序二进制文件 (会被更新)
├── README.md           # 说明文件 (会被更新)
├── config.yaml         # 配置文件 (默认不会被覆盖)
└── script.sh           # 脚本文件 (会被更新)
```

### 文件处理规则

- **二进制文件**: 直接替换
- **README.md**: 直接替换
- **脚本文件**: 直接替换
- **配置文件**: 仅在不存在时创建，存在时保留原文件

## 错误处理

SDK提供了详细的错误类型：

```go
var (
    ErrInvalidConfig        = errors.New("invalid configuration")
    ErrInvalidVersion      = errors.New("invalid version format")
    ErrNetworkTimeout      = errors.New("network timeout")
    ErrDownloadFailed      = errors.New("download failed")
    ErrVerificationFailed  = errors.New("file verification failed")
    ErrExtractionFailed    = errors.New("extraction failed")
    ErrUpdateFailed        = errors.New("update failed")
    ErrBackupFailed        = errors.New("backup failed")
    ErrNoUpdateAvailable   = errors.New("no update available")
)
```

## 安全特性

- **MD5校验**: 验证下载文件的完整性
- **路径安全**: 防止路径遍历攻击
- **备份机制**: 更新前自动创建备份
- **回滚支持**: 更新失败时自动恢复

## 最佳实践

1. **定期检查**: 建议定时检查更新，而不是每次启动都检查
2. **优雅关闭**: Web服务更新前应优雅关闭，避免数据丢失
3. **配置保护**: 合理配置 `PreserveFiles` 以保护重要文件
4. **错误处理**: 妥善处理各种错误情况
5. **用户体验**: CLI工具应询问用户确认后再执行更新

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！