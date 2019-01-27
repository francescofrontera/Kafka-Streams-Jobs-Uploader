package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func NormalizeJarName(jarName string) string {
	const internalImage = "internal_image"
	normalizeString := strings.TrimSuffix(strings.ToLower(jarName), filepath.Ext(jarName))

	return strings.Join([]string{internalImage, normalizeString}, "/")
}

func GetPathToJar(jarName string) (string, string) {
	workDirPath, err := os.Getwd(); if err != nil {
		log.Panic(err)
	}
	sourcePath := strings.Join([]string{workDirPath, "jars", jarName}, "/")
	targetPath := strings.Join([]string{"/jar", jarName}, "/")

	return sourcePath, targetPath
}
