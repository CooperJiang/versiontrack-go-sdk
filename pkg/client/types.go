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

// VersionDetail 版本详细信息
type VersionDetail struct {
	ID          string `json:"id"`
	ProjectID   string `json:"projectId"`
	Version     string `json:"version"`
	VersionCode int    `json:"versionCode"`
	VersionName string `json:"versionName"`
	Description string `json:"description"`
	Changelog   string `json:"changelog"`
	ForceUpdate bool   `json:"forceUpdate"`
	MinVersion  string `json:"minVersion"`
	Status      string `json:"status"`
	CreatedBy   string `json:"createdBy"`
	PublishedAt string `json:"publishedAt"`
}

// UpdateFile 更新文件信息
type UpdateFile struct {
	ID                string `json:"id"`
	VersionID         string `json:"versionId"`
	FileName          string `json:"fileName"`
	FilePath          string `json:"filePath"`
	FileSize          int64  `json:"fileSize"`
	FileHash          string `json:"fileHash"`
	Platform          string `json:"platform"`
	Arch              string `json:"arch"`
	FileType          string `json:"fileType"`
	DownloadURL       string `json:"downloadUrl"`
	IsCompressed      bool   `json:"isCompressed"`
	CompressionType   string `json:"compressionType"`
	Signature         string `json:"signature"`
	SignatureAlgorithm string `json:"signatureAlgorithm"`
	UploadStatus      string `json:"uploadStatus"`
}

// UpdateInfo 更新信息
type UpdateInfo struct {
	// 是否有更新
	HasUpdate bool `json:"hasUpdate"`
	// 最新版本详情
	LatestVersion *VersionDetail `json:"latestVersion"`
	// 当前版本
	CurrentVersion string `json:"currentVersion"`
	// 更新文件列表
	UpdateFiles []UpdateFile `json:"updateFiles"`
	// 是否强制更新
	IsForced bool `json:"isForced"`
	// 下载URL (从第一个匹配的文件获取)
	DownloadURL string `json:"-"`
	// 文件大小 (从第一个匹配的文件获取)
	FileSize int64 `json:"-"`
	// MD5哈希值 (从第一个匹配的文件获取)
	MD5Hash string `json:"-"`
	// 发布说明 (从版本详情获取)
	ReleaseNotes string `json:"-"`
	// 发布时间 (从版本详情获取)
	PublishedAt string `json:"-"`
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