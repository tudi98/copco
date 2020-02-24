package utils

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func TrimLines(b []byte) []byte {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	var lines [][]byte
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		lines = append(lines, line)
	}
	return bytes.Join(lines, []byte{'\n'})
}

func GetFilesWithExtension(ext string) []string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	path = path + separator + "tests"

	var files []string
	_ = filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() && filepath.Ext(path) == ext {
			files = append(files, "."+separator+"tests"+separator+f.Name())
		}
		return nil
	})

	return files
}
