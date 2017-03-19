package main

import (
	"fmt"
	"log"
	"os"
)

func du(currentPath string, info os.FileInfo) int64 {
	size := info.Size()
	if !info.IsDir() {
		return size
	}

	dir, err := os.Open(currentPath)

	if err != nil {
		log.Print(err)
		return size
	}
	defer dir.Close()
	fis, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fis {
		if fi.Name() == "." || fi.Name() == ".." {
			continue
		}
		size += du(currentPath+"/"+fi.Name(), fi)
	}
	fmt.Printf("%d %s\n", size, currentPath)
	return size
}
func main() {
	log.SetFlags(log.Llongfile)
	dir := os.Args[1]
	info, err := os.Lstat(dir)
	if err != nil {
		log.Fatal(err)
	}
	du(dir, info)
}
