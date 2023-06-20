package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// GetFileByte def
func GetFileByte(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil && os.IsNotExist(err) {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// FileWalk def
func FileWalk(fileDir string) []string {
	start, err := os.Stat(fileDir)
	if err != nil || !start.IsDir() {
		return nil
	}
	var targets []string
	filepath.Walk(fileDir, func(fpath string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		targets = append(targets, fpath)
		return nil
	})
	return targets
}

// String2Int def
func String2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}
