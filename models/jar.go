package models

type RunJar struct {
	JarName   string `json:jarName`
	MainClass string `json:mainClass`
}
