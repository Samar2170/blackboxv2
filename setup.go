package main

import (
	"os"
	"strings"
)

var UploadDir = "./uploads"

func setup() {
	dirs, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	uploadDirName := strings.Split(UploadDir, "/")[1]
	for _, dir := range dirs {
		if dir.IsDir() {
			if dir.Name() == uploadDirName {
				return
			}
		}
	}
	createUploadDir()
	createLogDir()
}

func createUploadDir() {
	err := os.Mkdir(UploadDir, 0755)
	if err != nil {
		panic(err)
	}
}

func createLogDir() {
	fds, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, fd := range fds {
		if fd.IsDir() {
			if fd.Name() == "logs" {
				return
			}
		}
	}
	err = os.Mkdir("logs", 0755)
	if err != nil {
		panic(err)
	}
}
