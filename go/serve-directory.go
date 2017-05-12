package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	port := 80
	ip := fmt.Sprintf("%s:%d", "127.0.0.1", port)

	fmt.Printf("binding web server at %s\n", ip)
	err := http.ListenAndServe(ip, http.FileServer(http.Dir(dir)))
	if err != nil {
		panic(err)
	}
}
