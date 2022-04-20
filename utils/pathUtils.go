package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// RootPath 获取文件的路径
func RootPath() string {
	s, err := exec.LookPath(os.Args[0])

	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	i := strings.LastIndex(s, "\\")

	path := s[0 : i+1]

	return path
}
