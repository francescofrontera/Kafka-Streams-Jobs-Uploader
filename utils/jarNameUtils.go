package utils

import (
	"os"
	"strings"
)

var (
	jarFolder = os.Getenv("JARS_FOLDER")
)

func GetPathToJar(jarName string) (string, string) {
	sourcePath := strings.Join([]string{jarFolder, "jars", jarName}, "/")
	targetPath := strings.Join([]string{"/jar", jarName}, "/")

	return sourcePath, targetPath
}
