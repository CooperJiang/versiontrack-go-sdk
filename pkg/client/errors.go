package client

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfig 配置无效错误
	ErrInvalidConfig = errors.New("invalid configuration")
	
	// ErrInvalidVersion 版本格式无效错误
	ErrInvalidVersion = errors.New("invalid version format")
	
	// ErrNetworkTimeout 网络超时错误
	ErrNetworkTimeout = errors.New("network timeout")
	
	// ErrDownloadFailed 下载失败错误
	ErrDownloadFailed = errors.New("download failed")
	
	// ErrVerificationFailed 文件校验失败错误
	ErrVerificationFailed = errors.New("file verification failed")
	
	// ErrExtractionFailed 解压失败错误
	ErrExtractionFailed = errors.New("extraction failed")
	
	// ErrUpdateFailed 更新失败错误
	ErrUpdateFailed = errors.New("update failed")
	
	// ErrBackupFailed 备份失败错误
	ErrBackupFailed = errors.New("backup failed")
	
	// ErrNoUpdateAvailable 无可用更新错误
	ErrNoUpdateAvailable = errors.New("no update available")
)

// ClientError 客户端错误类型
type ClientError struct {
	Code    string
	Message string
	Cause   error
}

// Error 实现error接口
func (e *ClientError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *ClientError) Unwrap() error {
	return e.Cause
}

// NewClientError 创建客户端错误
func NewClientError(code, message string, cause error) *ClientError {
	return &ClientError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}