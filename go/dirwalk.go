// traverse a given dir with maxdepth 1(like find) and prints
// file names without extension
// Usage : dirwalk.go <dir_name>

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Given an array , this will return unique elements in array

func uniqueNames(arg []string) []string {
	tempMap := make(map[string]uint8)
	for idx := range arg {
		tempMap[arg[idx]] = 0
	}
	tempArray := make([]string, 0)
	for key := range tempMap {
		tempArray = append(tempArray, key)
	}
	return tempArray
}

// Checks if the path provided exists or not

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()
	searchDir := flag.Arg(0)
	fileList := []string{}
	var aPath string

	if len(searchDir) == 0 {
		fmt.Println("Please specify the dir name")
		os.Exit(1)
	}

	if !exists(searchDir) {
		fmt.Println("Dir doesnt exists")
		os.Exit(1)
	}
	walkDir := func(path string, f os.FileInfo, err error) error {
		if exists(path) {
			fileInfo, err := os.Lstat(path)
			if err != nil {
				fmt.Println(err)
			}
			if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				aPath, _ = filepath.EvalSymlinks(path)
			}

			aPath = path

			fileList = append(fileList, filepath.Base(aPath))
		}
		return nil
	}

	err := filepath.Walk(searchDir, walkDir)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	uFileList := []string{}

	for _, fname := range fileList {
		if len(fname) > 0 {
			uFileList = append(uFileList, filepath.Base(fname))
		}
	}
	for _, file := range uniqueNames(uFileList) {
		fmt.Println(file)
	}
}
