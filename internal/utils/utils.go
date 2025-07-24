package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// GetExecutablePath 获取当前可执行文件路径
func GetExecutablePath() (string, error) {
	return os.Executable()
}

// EnsureDir 确保目录存在
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	// 确保目标目录存在
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// 复制权限
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// VerifyFileMD5 验证文件MD5
func VerifyFileMD5(filePath, expectedMD5 string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	actualMD5 := fmt.Sprintf("%x", hash.Sum(nil))
	if actualMD5 != expectedMD5 {
		return fmt.Errorf("MD5 mismatch: expected %s, got %s", expectedMD5, actualMD5)
	}

	return nil
}

// CreateTempDir 创建临时目录
func CreateTempDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}

// RemoveTempDir 删除临时目录
func RemoveTempDir(dir string) error {
	return os.RemoveAll(dir)
}

// RemoveFile 删除文件
func RemoveFile(path string) error {
	return os.Remove(path)
}

// GetPlatform 获取当前平台
func GetPlatform() string {
	switch runtime.GOOS {
	case "darwin":
		return "macos"
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	default:
		return runtime.GOOS
	}
}

// GetArch 获取当前架构
func GetArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "amd64"
	case "arm64":
		return "arm64"
	default:
		return runtime.GOARCH
	}
}