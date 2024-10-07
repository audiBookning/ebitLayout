package utils

import (
	"path/filepath"
	"runtime"
)

var (
	FilePathTxt          string
	Assets_Relative_Path = "../../"
)

func InitGetFilepath() {
	_, FilePathTxt, _, _ = runtime.Caller(0)
}

// GetFilePath constructs the full path for asset files.
func GetFilePath(fileName string) string {
	dir := filepath.Dir(FilePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}
