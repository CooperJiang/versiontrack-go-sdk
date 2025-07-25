package main

import "fmt"

// Version 信息
var (
	Version   = "v1.0.1"
	BuildTime = "2024-01-XX"
	GitCommit = "development"
)

func GetVersionInfo() string {
	return fmt.Sprintf(`VersionTrack Go SDK %s
Build Time: %s
Git Commit: %s
Go Version: Built with Go 1.21+
Platform Support: Windows, Linux, macOS
Architecture Support: amd64, arm64`, Version, BuildTime, GitCommit)
}

func main() {
	fmt.Println(GetVersionInfo())
}