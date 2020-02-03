package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func printFile(ignoreDirs []string) filepath.WalkFunc {
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		if info.IsDir() {
			dir := filepath.Base(path)
			for _, d := range ignoreDirs {
				if d == dir {
					return filepath.SkipDir
				}
			}
		}
		fmt.Println(path)
		return nil
	}
	return fn
}

func main() {

	log.SetFlags(log.Lshortfile)
	dir := os.Args[1]
	ignoreDirs := []string{".idea", ".vscode", ".git"}
	err := filepath.Walk(dir, printFile(ignoreDirs))
	if err != nil {
		log.Fatal(err)
	}
}
