package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// _avbdeviced:*:229:-2:Ethernet AVB Device Daemon:/var/empty:/usr/bin/false

type User struct {
	Name    string `json:"name"`
	Uid     string `json:"uid"`
	Comment string `json:"comment"`
	Home    string `json:"home"`
	Shell   string `json:"shell"`
}

type ResultMsg struct {
	User  `json:"user"`
	Error string `json:"error"`
}

// UserToMap returns a mapping of user names to user's attribute

func UserToMap() (UserMap map[string]User) {
	UserMap = make(map[string]User)
	data, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	details := strings.Split(string(data), "\n")
	for _, d := range details {
		var user User
		if len(d) == 0 || string(d[0]) == "#" {
			continue
		} else {
			l := strings.Split(d, ":")
			user.Name = l[0]
			user.Uid = l[2]
			user.Comment = l[4]
			user.Home = l[5]
			user.Shell = l[6]
		}
		UserMap[user.Name] = user
	}
	return
}

func findUser(user string) (udata User, err error) {
	userMap := UserToMap()
	for u, d := range userMap {
		if u == user {
			udata = d
			err = nil
			break
		} else {
			err = fmt.Errorf("no user found named %s", user)
		}
	}
	return
}

func SendOutput(w http.ResponseWriter, output ResultMsg) {
	jsonStr, err := json.MarshalIndent(output, "", "   ")
	if err != nil {
		panic("Can't marshal JSON: " + err.Error())
	}
	log.Printf("output: %s", jsonStr)
	if output.Error != "" {
		log.Printf("status code: %d", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		log.Printf("status code: %d", http.StatusOK)
	}
	w.Write(jsonStr)
}

func handler(w http.ResponseWriter, r *http.Request) {
	output := ResultMsg{}
	param := r.FormValue("u")
	udata, err := findUser(param)
	if err != nil {
		output.Error = "No user exists user : '" + param + "'"
		log.Println(output.Error)
	}

	log.Printf("user = %s", param)

	if output.Error == "" {
		output.User = udata
	}

	SendOutput(w, output)
}

func main() {
	port := flag.Int("port", 0, "TCP port for the HTTP server to listen on")
	flag.Parse()

	http.HandleFunc("/user", handler)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
