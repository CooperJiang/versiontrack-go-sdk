package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lilithgames/versiontrack-go-sdk/internal/archive"
	"github.com/lilithgames/versiontrack-go-sdk/internal/http"
	"github.com/lilithgames/versiontrack-go-sdk/internal/utils"
)

// Updater 更新器接口
type Updater interface {
	// CheckForUpdates 检查是否有可用更新
	CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error)
	
	// Download 下载更新文件
	Download(ctx context.Context, info *UpdateInfo, destPath string, callback ProgressCallback) error
	
	// Update 执行更新
	Update(ctx context.Context, info *UpdateInfo, downloadPath string) error
	
	// GetUpdateHistory 获取更新历史
	GetUpdateHistory() []UpdateRecord
	
	// Rollback 回滚到指定版本
	Rollback(ctx context.Context, version string) error
}

// Client VersionTrack客户端
type Client struct {
	config     *Config
	httpClient *http.Client
	history    []UpdateRecord
}

// NewClient 创建新的客户端实例
func NewClient(config *Config) (*Client, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 设置默认值
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.PreserveFiles == nil {
		config.PreserveFiles = []string{"config.yaml", "config.yml", "*.conf"}
	}
	if config.BackupCount == 0 {
		config.BackupCount = 3
	}

	httpClient := http.NewClient(config.ServerURL, config.Timeout)

	return &Client{
		config:     config,
		httpClient: httpClient,
		history:    make([]UpdateRecord, 0),
	}, nil
}

// CheckForUpdates 检查是否有可用更新
func (c *Client) CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error) {
	url := fmt.Sprintf("/api/v1/public/versions/check?projectId=%s&platform=%s&arch=%s&currentVersion=%s",
		c.config.ProjectID, c.config.Platform, c.config.Arch, currentVersion)

	var result struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    *UpdateInfo `json:"data"`
	}

	if err := c.httpClient.Get(ctx, url, &result); err != nil {
		return nil, NewClientError("CHECK_FAILED", "Failed to check for updates", err)
	}

	if result.Code != 200 {
		return nil, NewClientError("API_ERROR", result.Message, nil)
	}

	return result.Data, nil
}

// Download 下载更新文件
func (c *Client) Download(ctx context.Context, info *UpdateInfo, destPath string, callback ProgressCallback) error {
	if info == nil {
		return NewClientError("INVALID_INFO", "Update info is nil", nil)
	}

	if !info.HasUpdate {
		return ErrNoUpdateAvailable
	}

	// 创建目标目录
	if err := utils.EnsureDir(filepath.Dir(destPath)); err != nil {
		return NewClientError("CREATE_DIR_FAILED", "Failed to create destination directory", err)
	}

	// 下载文件
	if err := c.httpClient.DownloadFile(ctx, info.DownloadURL, destPath, func(downloaded, total int64) {
		if callback != nil {
			progress := &DownloadProgress{
				Downloaded: downloaded,
				Total:      total,
				Speed:      0, // TODO: 计算下载速度
				Percentage: float64(downloaded) / float64(total) * 100,
			}
			callback(progress)
		}
	}); err != nil {
		return NewClientError("DOWNLOAD_FAILED", "Failed to download update file", err)
	}

	// 验证文件
	if err := utils.VerifyFileMD5(destPath, info.MD5Hash); err != nil {
		return NewClientError("VERIFY_FAILED", "File verification failed", err)
	}

	return nil
}

// Update 执行更新
func (c *Client) Update(ctx context.Context, info *UpdateInfo, downloadPath string) error {
	if info == nil {
		return NewClientError("INVALID_INFO", "Update info is nil", nil)
	}

	// 1. 创建备份
	backupPath, err := c.createBackup()
	if err != nil {
		return NewClientError("BACKUP_FAILED", "Failed to create backup", err)
	}

	// 2. 解压更新文件
	tempDir, err := utils.CreateTempDir("versiontrack-update")
	if err != nil {
		return NewClientError("CREATE_TEMP_FAILED", "Failed to create temp directory", err)
	}
	defer utils.RemoveTempDir(tempDir)

	if err := archive.ExtractTarGz(downloadPath, tempDir); err != nil {
		return NewClientError("EXTRACT_FAILED", "Failed to extract update file", err)
	}

	// 3. 应用更新
	if err := c.applyUpdate(tempDir); err != nil {
		// 更新失败，尝试回滚
		if rollbackErr := c.restoreBackup(backupPath); rollbackErr != nil {
			return NewClientError("UPDATE_AND_ROLLBACK_FAILED", 
				fmt.Sprintf("Update failed: %v, Rollback also failed: %v", err, rollbackErr), nil)
		}
		return NewClientError("UPDATE_FAILED", "Update failed, rolled back successfully", err)
	}

	// 4. 记录更新历史
	record := UpdateRecord{
		Version:    info.LatestVersion,
		UpdatedAt:  time.Now(),
		Status:     "success",
		BackupPath: backupPath,
	}
	c.history = append(c.history, record)

	// 5. 清理旧备份
	c.cleanupOldBackups()

	return nil
}

// GetUpdateHistory 获取更新历史
func (c *Client) GetUpdateHistory() []UpdateRecord {
	return c.history
}

// Rollback 回滚到指定版本
func (c *Client) Rollback(ctx context.Context, version string) error {
	// 查找对应版本的备份
	var targetRecord *UpdateRecord
	for i := range c.history {
		if c.history[i].Version == version {
			targetRecord = &c.history[i]
			break
		}
	}

	if targetRecord == nil {
		return NewClientError("BACKUP_NOT_FOUND", "Backup for version not found", nil)
	}

	// 执行回滚
	if err := c.restoreBackup(targetRecord.BackupPath); err != nil {
		return NewClientError("ROLLBACK_FAILED", "Failed to rollback", err)
	}

	return nil
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config == nil {
		return ErrInvalidConfig
	}
	if config.ServerURL == "" {
		return fmt.Errorf("ServerURL is required")
	}
	if config.ProjectID == "" {
		return fmt.Errorf("ProjectID is required")
	}
	if config.Platform == "" {
		return fmt.Errorf("Platform is required")
	}
	if config.Arch == "" {
		return fmt.Errorf("Arch is required")
	}

	// 验证平台和架构
	validPlatforms := []string{"windows", "linux", "macos"}
	validArchs := []string{"amd64", "arm64"}

	if !contains(validPlatforms, config.Platform) {
		return fmt.Errorf("invalid platform: %s, must be one of %v", config.Platform, validPlatforms)
	}
	if !contains(validArchs, config.Arch) {
		return fmt.Errorf("invalid arch: %s, must be one of %v", config.Arch, validArchs)
	}

	return nil
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// createBackup 创建当前版本的备份
func (c *Client) createBackup() (string, error) {
	// 获取当前执行文件路径
	execPath, err := utils.GetExecutablePath()
	if err != nil {
		return "", err
	}

	// 创建备份目录
	backupDir := filepath.Join(filepath.Dir(execPath), ".versiontrack", "backups")
	if err := utils.EnsureDir(backupDir); err != nil {
		return "", err
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("backup_%s.tar.gz", timestamp))

	// 创建备份
	currentDir := filepath.Dir(execPath)
	if err := archive.CreateTarGz(currentDir, backupPath, c.config.PreserveFiles); err != nil {
		return "", err
	}

	return backupPath, nil
}

// applyUpdate 应用更新
func (c *Client) applyUpdate(updateDir string) error {
	execPath, err := utils.GetExecutablePath()
	if err != nil {
		return err
	}

	currentDir := filepath.Dir(execPath)

	// 遍历更新文件
	return filepath.Walk(updateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(updateDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(currentDir, relPath)

		// 检查是否是需要保护的文件
		if c.shouldPreserveFile(relPath) {
			// 如果目标文件已存在，跳过覆盖
			if utils.FileExists(targetPath) {
				return nil
			}
		}

		// 复制文件
		return utils.CopyFile(path, targetPath)
	})
}

// shouldPreserveFile 检查文件是否需要保护
func (c *Client) shouldPreserveFile(filename string) bool {
	for _, pattern := range c.config.PreserveFiles {
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
		if strings.Contains(filename, pattern) {
			return true
		}
	}
	return false
}

// restoreBackup 恢复备份
func (c *Client) restoreBackup(backupPath string) error {
	execPath, err := utils.GetExecutablePath()
	if err != nil {
		return err
	}

	currentDir := filepath.Dir(execPath)

	// 解压备份到当前目录
	return archive.ExtractTarGz(backupPath, currentDir)
}

// cleanupOldBackups 清理旧备份
func (c *Client) cleanupOldBackups() {
	if len(c.history) <= c.config.BackupCount {
		return
	}

	// 删除超过保留数量的备份
	for i := 0; i < len(c.history)-c.config.BackupCount; i++ {
		backupPath := c.history[i].BackupPath
		if backupPath != "" {
			utils.RemoveFile(backupPath)
		}
	}

	// 更新历史记录
	c.history = c.history[len(c.history)-c.config.BackupCount:]
}