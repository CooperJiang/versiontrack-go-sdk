package client

import (
	"time"
)

// Config 客户端配置
type Config struct {
	// VersionTrack服务器地址
	ServerURL string
	// 项目ID
	ProjectID string
	// 平台信息 (windows/linux/macos)
	Platform string
	// 架构信息 (amd64/arm64)
	Arch string
	// HTTP请求超时时间
	Timeout time.Duration
	// 需要保护的文件列表（更新时不覆盖）
	PreserveFiles []string
	// 备份保留数量
	BackupCount int
}

// UpdateInfo 更新信息
type UpdateInfo struct {
	// 是否有更新
	HasUpdate bool `json:"hasUpdate"`
	// 最新版本号
	LatestVersion string `json:"latestVersion"`
	// 下载URL
	DownloadURL string `json:"downloadUrl"`
	// 文件大小
	FileSize int64 `json:"fileSize"`
	// MD5哈希值
	MD5Hash string `json:"md5Hash"`
	// 发布说明
	ReleaseNotes string `json:"releaseNotes"`
	// 发布时间
	PublishedAt string `json:"publishedAt"`
}

// UpdateRecord 更新记录
type UpdateRecord struct {
	// 版本号
	Version string `json:"version"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt"`
	// 更新状态
	Status string `json:"status"`
	// 备份路径
	BackupPath string `json:"backupPath"`
}

// DownloadProgress 下载进度信息
type DownloadProgress struct {
	// 已下载字节数
	Downloaded int64
	// 总字节数
	Total int64
	// 下载速度 (bytes/second)
	Speed int64
	// 百分比
	Percentage float64
}

// ProgressCallback 下载进度回调函数
type ProgressCallback func(progress *DownloadProgress)