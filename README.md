# VersionTrack Go SDK

一个用于Go应用程序版本管理和热更新的SDK，与VersionTrack系统集成。

**版本：v1.0.4** - 🆕 支持预览发布功能，增强版本状态管理

## 特性

- 🚀 **自动版本检查**：定期检查是否有新版本发布
- 📦 **智能更新**：支持tar.gz包的解压和安装
- 🔒 **配置保护**：更新时自动保护重要配置文件
- 📊 **进度跟踪**：实时显示下载和更新进度
- 🔄 **回滚支持**：更新失败时自动回滚
- 📝 **更新历史**：记录所有更新操作的历史
- 🛡️ **安全验证**：MD5校验确保文件完整性
- 🔑 **API密钥认证**：仅使用安全的API密钥进行认证，简化配置
- 🎯 **多版本管理**：支持多版本检查和选择更新
- ⚡ **强制更新**：支持强制更新策略
- 🕐 **预览发布**：支持scheduled状态版本的预览和下载时间管理
- 📋 **版本状态**：完整的版本状态管理（draft, published, scheduled, recalled, archived）
- ⏰ **定时发布**：支持版本预定发布时间和下载可用时间控制

## 快速开始

### 安装

```bash
go get github.com/CooperJiang/versiontrack-go-sdk@v1.0.2
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
    // 配置客户端 - 简化配置，仅需API密钥
    config := &client.Config{
        ServerURL:     "https://your-versiontrack-server.com",
        APIKey:        "your-api-key-here",  // 🆕 仅使用API密钥，无需ProjectID
        Platform:      utils.GetPlatform(),  // 自动检测平台
        Arch:          utils.GetArch(),      // 自动检测架构
        Timeout:       30 * time.Second,
        PreserveFiles: []string{"config.yaml", "*.conf"},
        BackupCount:   3,
        UpdateMode:    client.UpdateModeAuto,
    }

    // 创建客户端
    updater, err := client.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // 检查更新 - SDK内部使用API密钥进行身份验证
    ctx := context.Background()
    updateInfo, err := updater.CheckForUpdates(ctx, "1.0.0")
    if err != nil {
        log.Fatal(err)
    }

    if updateInfo.HasUpdate {
        // 下载更新
        err = updater.Download(ctx, updateInfo, "/tmp/update.tar.gz", func(progress *client.DownloadProgress) {
            log.Printf("下载进度: %.1f%%", progress.Percentage)
        })
        if err != nil {
            log.Fatal(err)
        }
        
        // 执行更新
        err = updater.Update(ctx, updateInfo, "/tmp/update.tar.gz")
        if err != nil {
            log.Fatal(err)
        }
    }
}
```

## 配置说明

### Config 结构

```go
type Config struct {
    ServerURL     string        // VersionTrack服务器地址
    APIKey        string        // API密钥 (必须)
    Platform      string        // 平台 (windows/linux/macos)
    Arch          string        // 架构 (amd64/arm64)
    Timeout       time.Duration // HTTP请求超时时间
    PreserveFiles []string      // 需要保护的文件列表
    BackupCount   int          // 备份保留数量
    UpdateMode    UpdateMode   // 更新模式
    SkipVersions  []string     // 跳过的版本列表
}
```

### 参数说明

- **ServerURL**: VersionTrack服务器的API地址
- **APIKey**: 🆕 在VersionTrack管理后台项目设置中获取的API密钥 (必须)
- **Platform**: 目标平台，支持 `windows`、`linux`、`macos`
- **Arch**: 目标架构，支持 `amd64`、`arm64`
- **Timeout**: HTTP请求超时时间，默认30秒
- **PreserveFiles**: 更新时不覆盖的文件模式列表，默认包含 `config.yaml`
- **BackupCount**: 保留的备份数量，默认3个
- **UpdateMode**: 更新模式，支持 `auto`/`manual`/`prompt`
- **SkipVersions**: 跳过的版本列表，这些版本不会被自动更新

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
    // 基础API - 推荐使用
    CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error)
    Download(ctx context.Context, info *UpdateInfo, destPath string, callback ProgressCallback) error
    Update(ctx context.Context, info *UpdateInfo, downloadPath string) error
    
    // 高级API - 多版本管理
    CheckForMultipleUpdates(ctx context.Context, currentVersion string) (*UpdatesInfo, error)
    GetRecommendedUpdate(ctx context.Context, currentVersion string) (*VersionInfo, error)
    UpdateToVersion(ctx context.Context, targetVersion string, callback ProgressCallback) error
    HasForcedUpdate(ctx context.Context, currentVersion string) (*VersionInfo, error)
    DownloadVersion(ctx context.Context, versionInfo *VersionInfo, destPath string, callback ProgressCallback) error
    
    // 通用功能
    GetUpdateHistory() []UpdateRecord
    Rollback(ctx context.Context, version string) error
}
```

### 数据结构

#### UpdateInfo - 更新信息
```go
type UpdateInfo struct {
    HasUpdate      bool           `json:"hasUpdate"`      // 是否有更新
    LatestVersion  *VersionDetail `json:"latestVersion"`  // 最新版本详情
    CurrentVersion string         `json:"currentVersion"` // 当前版本
    UpdateFiles    []UpdateFile   `json:"updateFiles"`    // 更新文件列表
    IsForced       bool          `json:"isForced"`       // 是否强制更新
    // 以下字段从API响应中自动填充
    DownloadURL    string        `json:"-"`              // 下载URL
    FileSize       int64         `json:"-"`              // 文件大小
    MD5Hash        string        `json:"-"`              // MD5哈希值
    ReleaseNotes   string        `json:"-"`              // 发布说明
    PublishedAt    string        `json:"-"`              // 发布时间
}
```

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

## 使用示例

### 1. 基础示例 - 简单更新检查

```go
// 配置客户端 - 仅需API密钥
config := &client.Config{
    ServerURL:    "http://localhost:9000",
    APIKey:       "your-api-key-here",      // 仅需提供API密钥
    Platform:     utils.GetPlatform(),
    Arch:         utils.GetArch(),
    UpdateMode:   client.UpdateModeAuto,
}

updater, _ := client.NewClient(config)

// 检查更新 - SDK自动使用API密钥进行认证
updateInfo, err := updater.CheckForUpdates(ctx, "1.0.0")
if err != nil {
    log.Fatal(err)
}

if updateInfo.HasUpdate {
    log.Printf("发现新版本: %s", updateInfo.LatestVersion.Version)
    log.Printf("更新说明: %s", updateInfo.LatestVersion.Changelog)
    
    // 下载并安装更新
    err = updater.Download(ctx, updateInfo, "/tmp/update.tar.gz", nil)
    if err == nil {
        err = updater.Update(ctx, updateInfo, "/tmp/update.tar.gz")
    }
}
```

### 2. Web服务示例 - 优雅更新

参见 [examples/web-service/main.go](examples/web-service/main.go)

这个示例展示了如何在Web服务中集成自动更新功能，包括：
- 🆕 使用简化的API密钥认证
- 定时检查更新
- 优雅关闭服务
- 手动触发更新的API接口
- 强制更新处理

### 3. CLI工具示例 - 手动选择版本

参见 [examples/cli-tool/main.go](examples/cli-tool/main.go)

这个示例展示了如何为命令行工具添加更新功能：
- 🆕 简化配置，仅使用API密钥
- 手动选择版本更新
- 用户交互确认
- 配置文件保护
- 跳过指定版本

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
- **API密钥认证**: 安全的身份验证机制

## 版本历史

### 🆕 v1.0.2 版本变更

#### 重要变更
- **简化配置**: 移除 `ProjectID` 字段，仅使用 `APIKey` 进行认证
- **统一认证**: SDK内部统一使用API密钥进行所有API调用的身份验证
- **精简配置**: 减少配置参数，提高易用性
- **API优化**: 更新API调用URL，移除projectId参数

#### 配置变更
```go
// v1.0.1 (旧版本)
config := &client.Config{
    APIKey:    "your-api-key",
    ProjectID: "your-project-id",  // ❌ 已移除
}

// v1.0.2 (新版本)
config := &client.Config{
    APIKey: "your-api-key",  // ✅ 仅需API密钥
}
```

#### 迁移指南
从 v1.0.1 升级到 v1.0.2：

1. **更新配置结构**：
   - 移除配置中的 `ProjectID` 字段
   - 确保 `APIKey` 字段正确设置

2. **API调用变更**：
   - SDK内部API调用不再使用projectId参数
   - 服务器通过API密钥自动识别项目信息

3. **无需额外操作**：
   - 其他API接口保持不变
   - 功能和使用方式完全兼容

### v1.0.1 版本变更

#### 重要变更
- **API密钥认证**: 引入 `APIKey` 字段进行认证
- **多版本管理**: 支持检查和选择多个可用版本
- **更新模式**: 新增自动、手动、提示三种更新模式
- **强制更新**: 支持强制更新策略和最低版本要求
- **版本跳过**: 支持跳过指定版本的更新

## 最佳实践

1. **定期检查**: 建议定时检查更新，而不是每次启动都检查
2. **优雅关闭**: Web服务更新前应优雅关闭，避免数据丢失
3. **配置保护**: 合理配置 `PreserveFiles` 以保护重要文件
4. **错误处理**: 妥善处理各种错误情况
5. **用户体验**: CLI工具应询问用户确认后再执行更新
6. **强制更新**: 对于安全补丁等重要更新，建议使用强制更新策略
7. **版本策略**: 合理设置更新模式，平衡自动化和用户控制
8. **🆕 简化配置**: 使用v1.0.2的简化配置，减少出错可能性

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！