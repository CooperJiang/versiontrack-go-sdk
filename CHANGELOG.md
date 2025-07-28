# VersionTrack Go SDK 更新日志

## [v1.0.4] - 2025-07-28

### 新增功能
- 🎉 **预览发布支持**: 新增对VersionTrack预览发布功能的完整支持
- 📊 **版本权重**: 添加`VersionWeight`字段用于精确的版本排序和比较  
- 🕒 **预定发布时间**: 支持`ScheduledReleaseAt`字段显示版本的预定发布时间
- 📋 **版本状态**: 新增`Status`字段显示版本状态 (draft, published, scheduled, recalled, archived)
- ⬇️ **下载控制**: 添加`IsDownloadable`和`DownloadableStatus`字段控制版本下载权限

### 重要变更
- 🔄 **字段重命名**: `VersionInfo.VersionCode` → `VersionInfo.VersionWeight` (更语义化)
- 🗑️ **字段移除**: 移除了`VersionDetail`中的一些不常用字段:
  - `Description` (使用`Changelog`替代)
  - `ForceUpdate` (使用`IsForced`替代) 
  - `MinVersion` (简化版本管理)

### API增强
- ✨ **扩展版本信息**: `VersionInfo`结构体新增多个字段支持预览发布工作流
- 📱 **兼容性保持**: 保持向后兼容，现有API调用仍然有效
- 🔐 **认证优化**: 改进API密钥认证和下载权限控制

### 使用示例

```go
// 检查更新现在会返回更详细的版本信息
updates, err := client.CheckForMultipleUpdates(ctx, currentVersion)
if err != nil {
    log.Fatal(err)
}

for _, version := range updates.AvailableVersions {
    fmt.Printf("版本: %s\n", version.Version)
    fmt.Printf("状态: %s\n", version.Status)
    fmt.Printf("是否可下载: %v\n", version.IsDownloadable)
    fmt.Printf("下载状态: %s\n", version.DownloadableStatus)
    
    if version.ScheduledReleaseAt != "" {
        fmt.Printf("预定发布时间: %s\n", version.ScheduledReleaseAt)
    }
}
```

### 升级指南

从v1.0.3升级到v1.0.4：

1. **字段名变更**: 如果直接访问`VersionInfo.VersionCode`，请改为`VersionInfo.VersionWeight`
2. **新字段利用**: 利用新的状态和下载控制字段优化用户体验
3. **预览发布**: 可以显示预览版本但控制其下载权限

### 兼容性
- ✅ Go 1.21+
- ✅ 向后兼容v1.0.3的API调用
- ✅ 支持VersionTrack服务端v1.5.0+

---

**完整更新详情请参考**: [GitHub Releases](https://github.com/CooperJiang/versiontrack-go-sdk/releases/tag/v1.0.4)