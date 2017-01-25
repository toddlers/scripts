package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UniqueNames(arg []string) []string {
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

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	//searchDir := "/Users/suresh.prajapati/src/golang/unix"
	flag.Parse()
	searchDir := flag.Arg(0)
	fileList := []string{}
	var aPath string
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if Exists(path) {
			fileInfo, err := os.Lstat(path)
			if err != nil {
				fmt.Println(err)
			}
			if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				aPath, _ = filepath.EvalSymlinks(path)
			} else {
				aPath = path
			}

			fileNames := strings.Split(aPath, "/")
			fileList = append(fileList, fileNames[len(fileNames)-1])
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Printf("%v\n", fileList)
	uFileList := []string{}

	for _, fname := range fileList {

		if len(fname) > 0 {
			fExt := filepath.Ext(fname)
			fName := fname[0 : len(fname)-len(fExt)]
			uFileList = append(uFileList, fName)
		}
	}
	for _, file := range UniqueNames(uFileList) {
		if len(file) > 0 {
			fmt.Println(file)
		}
	}
}
