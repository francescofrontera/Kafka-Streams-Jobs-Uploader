package utils

import (
	"strings"
)


func GetPathToJar(jarName string) (string, string) {
	path := "/Users/francescofrontera/Go/src/github.com/francescofrontera/ks-job-uploader"

	sourcePath := strings.Join([]string{path, "jars", jarName}, "/")

	targetPath := strings.Join([]string{"/jar", jarName}, "/")

	return sourcePath, targetPath
}
