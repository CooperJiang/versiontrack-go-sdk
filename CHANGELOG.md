# Changelog

所有对 VersionTrack Go SDK 的重要更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [v1.0.1] - 2024-01-XX

### 🆕 新增功能
- **API密钥认证**: 使用更安全的 APIKey 替代 ProjectID 进行认证
- **多版本管理**: 新增 `CheckForMultipleUpdates()` 方法支持检查多个可用版本
- **智能推荐**: 新增 `GetRecommendedUpdate()` 方法自动选择推荐版本
- **手动版本选择**: 新增 `UpdateToVersion()` 方法支持手动指定更新版本
- **强制更新策略**: 新增 `HasForcedUpdate()` 方法检查强制更新要求
- **更新模式**: 支持自动、手动、提示三种更新模式
- **版本跳过**: 支持配置跳过指定版本的更新
- **增强的示例**: 更新所有示例代码展示新功能

### 🔄 变更
- 配置结构中 `ProjectID` 字段已弃用，改为使用 `APIKey`
- 优化了错误处理和用户体验
- 改进了示例代码的交互性和可读性

### 🔧 修复
- 修复了下载进度显示的问题
- 改进了网络超时处理
- 优化了文件校验逻辑

### 📚 文档
- 全面更新 README.md 文档
- 新增版本迁移指南
- 更新所有使用示例
- 添加最佳实践建议

### ⚠️ 弃用警告
- `ProjectID` 配置参数已弃用，请使用 `APIKey`
- 旧版本 API 仍然可用但建议迁移到新版本

### 🚀 向后兼容
- 保持与 v1.0.0 的完全向后兼容
- 旧版本 API（如 `CheckForUpdates`）继续可用
- 逐步迁移策略，无需强制升级

---

## [v1.0.0] - 2024-01-XX

### 🎉 首次发布
- 基础的版本检查和更新功能
- 支持 Windows、Linux、macOS 平台
- tar.gz 格式的更新包支持
- 文件保护和备份机制
- MD5 校验确保文件完整性
- 自动回滚功能
- 基础示例和文档

[v1.0.1]: https://github.com/CooperJiang/versiontrack-go-sdk/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/CooperJiang/versiontrack-go-sdk/releases/tag/v1.0.0