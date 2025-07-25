# VersionTrack Go SDK

一个用于Go应用程序版本管理和热更新的SDK，与VersionTrack系统集成。

**版本：v1.0.1** - 🆕 全新API设计，支持新的VersionTrack后端系统

## 特性

- 🚀 **自动版本检查**：定期检查是否有新版本发布
- 📦 **智能更新**：支持tar.gz包的解压和安装
- 🔒 **配置保护**：更新时自动保护重要配置文件
- 📊 **进度跟踪**：实时显示下载和更新进度
- 🔄 **回滚支持**：更新失败时自动回滚
- 📝 **更新历史**：记录所有更新操作的历史
- 🛡️ **安全验证**：MD5校验确保文件完整性
- 🔑 **API密钥认证**：使用安全的API密钥进行认证
- 🎯 **多版本管理**：支持多版本检查和选择更新
- ⚡ **强制更新**：支持强制更新策略

## 快速开始

### 安装

```bash
go get github.com/CooperJiang/versiontrack-go-sdk@v1.0.1
```

### 基本使用

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/CooperJiang/versiontrack-go-sdk/internal/utils"
    "github.com/CooperJiang/versiontrack-go-sdk/pkg/client"
)

func main() {
    // 配置客户端
    config := &client.Config{
        ServerURL:     "https://your-versiontrack-server.com",
        APIKey:        "your-api-key-here",  // 🆕 使用API密钥替代ProjectID
        Platform:      utils.GetPlatform(),  // 自动检测平台
        Arch:          utils.GetArch(),      // 自动检测架构
        Timeout:       30 * time.Second,
        PreserveFiles: []string{"config.yaml", "*.conf"},
        BackupCount:   3,
        UpdateMode:    client.UpdateModeAuto, // 🆕 支持多种更新模式
    }

    // 创建客户端
    updater, err := client.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // 🆕 检查多版本更新
    ctx := context.Background()
    updatesInfo, err := updater.CheckForMultipleUpdates(ctx, "1.0.0")
    if err != nil {
        log.Fatal(err)
    }

    if updatesInfo.HasUpdate {
        // 🆕 获取推荐更新版本
        recommendedVersion, err := updater.GetRecommendedUpdate(ctx, "1.0.0")
        if err != nil {
            log.Fatal(err)
        }
        
        if recommendedVersion != nil {
            // 🆕 更新到指定版本
            err = updater.UpdateToVersion(ctx, recommendedVersion.Version, func(progress *client.DownloadProgress) {
                log.Printf("下载进度: %.1f%%", progress.Percentage)
            })
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}
```

## 配置说明

### Config 结构

```go
type Config struct {
    ServerURL     string        // VersionTrack服务器地址
    APIKey        string        // API密钥 🆕 (替代ProjectID)
    Platform      string        // 平台 (windows/linux/macos)
    Arch          string        // 架构 (amd64/arm64)
    Timeout       time.Duration // HTTP请求超时时间
    PreserveFiles []string      // 需要保护的文件列表
    BackupCount   int          // 备份保留数量
    UpdateMode    UpdateMode   // 🆕 更新模式
    SkipVersions  []string     // 🆕 跳过的版本列表
}
```

### 参数说明

- **ServerURL**: VersionTrack服务器的API地址
- **APIKey**: 🆕 在VersionTrack管理后台项目设置中获取的API密钥
- **Platform**: 目标平台，支持 `windows`、`linux`、`macos`
- **Arch**: 目标架构，支持 `amd64`、`arm64`
- **Timeout**: HTTP请求超时时间，默认30秒
- **PreserveFiles**: 更新时不覆盖的文件模式列表，默认包含 `config.yaml`
- **BackupCount**: 保留的备份数量，默认3个
- **UpdateMode**: 🆕 更新模式，支持 `auto`/`manual`/`prompt`
- **SkipVersions**: 🆕 跳过的版本列表，这些版本不会被自动更新

### 🆕 更新模式说明

```go
const (
    UpdateModeAuto   UpdateMode = "auto"   // 自动更新到推荐版本
    UpdateModeManual UpdateMode = "manual" // 手动选择版本
    UpdateModePrompt UpdateMode = "prompt" // 提示用户选择
)
```

## 主要接口

### Updater 接口

```go
type Updater interface {
    // 🆕 新版本API - 支持多版本管理
    CheckForMultipleUpdates(ctx context.Context, currentVersion string) (*UpdatesInfo, error)
    GetRecommendedUpdate(ctx context.Context, currentVersion string) (*VersionInfo, error)
    UpdateToVersion(ctx context.Context, targetVersion string, callback ProgressCallback) error
    HasForcedUpdate(ctx context.Context, currentVersion string) (*VersionInfo, error)
    DownloadVersion(ctx context.Context, versionInfo *VersionInfo, destPath string, callback ProgressCallback) error
    
    // 兼容旧版本API - 保持向后兼容
    CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error)
    Download(ctx context.Context, info *UpdateInfo, destPath string, callback ProgressCallback) error
    Update(ctx context.Context, info *UpdateInfo, downloadPath string) error
    
    // 通用功能
    GetUpdateHistory() []UpdateRecord
    Rollback(ctx context.Context, version string) error
}
```

### 🆕 新版本数据结构

#### UpdatesInfo - 多版本更新信息
```go
type UpdatesInfo struct {
    HasUpdate         bool          `json:"hasUpdate"`         // 是否有更新
    CurrentVersion    string        `json:"currentVersion"`    // 当前版本
    LatestVersion     string        `json:"latestVersion"`     // 最新版本
    AvailableVersions []VersionInfo `json:"availableVersions"` // 可用版本列表
    UpdateStrategy    UpdateStrategy `json:"updateStrategy"`   // 更新策略
}
```

#### VersionInfo - 版本信息
```go
type VersionInfo struct {
    Version     string `json:"version"`     // 版本号
    VersionCode int64  `json:"versionCode"` // 版本代码
    Changelog   string `json:"changelog"`   // 更新日志
    ReleaseDate string `json:"releaseDate"` // 发布日期
    IsForced    bool   `json:"isForced"`    // 是否强制更新
    DownloadURL string `json:"downloadUrl"` // 下载地址
    FileSize    int64  `json:"fileSize"`    // 文件大小
    FileHash    string `json:"fileHash"`    // 文件哈希
}
```

#### UpdateStrategy - 更新策略
```go
type UpdateStrategy struct {
    HasForced          bool   `json:"hasForced"`          // 是否有强制更新
    MinRequiredVersion string `json:"minRequiredVersion"` // 最低要求版本
}
```

## 使用示例

### 1. 基础示例 - 多版本管理

参见 [examples/basic/main.go](examples/basic/main.go)

```go
// 配置客户端
config := &client.Config{
    ServerURL:    "http://localhost:9000",
    APIKey:       "your-api-key-here",
    Platform:     utils.GetPlatform(),
    Arch:         utils.GetArch(),
    UpdateMode:   client.UpdateModeAuto,
}

updater, _ := client.NewClient(config)

// 🆕 检查多版本更新
updatesInfo, err := updater.CheckForMultipleUpdates(ctx, "1.0.0")
if updatesInfo.HasUpdate {
    // 自动获取推荐版本并更新
    recommendedVersion, _ := updater.GetRecommendedUpdate(ctx, "1.0.0")
    if recommendedVersion != nil {
        updater.UpdateToVersion(ctx, recommendedVersion.Version, progressCallback)
    }
}
```

### 2. Web服务示例 - 优雅更新

参见 [examples/web-service/main.go](examples/web-service/main.go)

这个示例展示了如何在Web服务中集成自动更新功能，包括：
- 🆕 使用新的API密钥认证
- 定时检查更新
- 优雅关闭服务
- 手动触发更新的API接口
- 🆕 强制更新处理

### 3. CLI工具示例 - 手动选择版本

参见 [examples/cli-tool/main.go](examples/cli-tool/main.go)

这个示例展示了如何为命令行工具添加更新功能：
- 🆕 手动选择版本更新
- 用户交互确认
- 配置文件保护
- 🆕 跳过指定版本

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

## 🆕 v1.0.1 版本变更

### 重要变更
- **API密钥认证**: 使用 `APIKey` 替代 `ProjectID` 进行认证
- **多版本管理**: 支持检查和选择多个可用版本
- **更新模式**: 新增自动、手动、提示三种更新模式
- **强制更新**: 支持强制更新策略和最低版本要求
- **版本跳过**: 支持跳过指定版本的更新

### 向后兼容性
为了保持向后兼容，SDK 同时提供新旧两套 API：
- 旧版本 API（如 `CheckForUpdates`）仍然可用
- 新版本 API 提供更丰富的功能

### 迁移指南
从 v1.0.0 升级到 v1.0.1：

1. **更新配置结构**：
```go
// 旧版本
config := &client.Config{
    ProjectID: "your-project-id",  // ❌ 已弃用
}

// 新版本
config := &client.Config{
    APIKey: "your-api-key-here",   // ✅ 新的认证方式
    UpdateMode: client.UpdateModeAuto, // ✅ 新增更新模式
}
```

2. **使用新的API方法**：
```go
// 推荐使用新的多版本API
updatesInfo, err := updater.CheckForMultipleUpdates(ctx, currentVersion)
recommendedVersion, err := updater.GetRecommendedUpdate(ctx, currentVersion)
err = updater.UpdateToVersion(ctx, targetVersion, callback)
```

3. **获取API密钥**：
   - 登录 VersionTrack 管理后台
   - 进入项目设置页面
   - 在 API Keys 部分生成新的密钥

## 最佳实践

1. **定期检查**: 建议定时检查更新，而不是每次启动都检查
2. **优雅关闭**: Web服务更新前应优雅关闭，避免数据丢失
3. **配置保护**: 合理配置 `PreserveFiles` 以保护重要文件
4. **错误处理**: 妥善处理各种错误情况
5. **用户体验**: CLI工具应询问用户确认后再执行更新
6. **🆕 强制更新**: 对于安全补丁等重要更新，建议使用强制更新策略
7. **🆕 版本策略**: 合理设置更新模式，平衡自动化和用户控制

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！