package pathEx

import (
	"path"
	"strings"
)

// c#习惯下的path
func GetExtension(filePath string) string {
	return path.Ext(filePath)
}

func GetFileName(filePath string) string {
	return path.Base(filePath)
}

func GetFileNameWithoutExtension(filePath string) string {
	return strings.TrimSuffix(path.Base(filePath), path.Ext(filePath))
}

func GetFullPath(filePath string) string {
	return path.Dir(filePath)
}
