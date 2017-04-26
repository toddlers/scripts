package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Content struct {
	Username  string `json:"username"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
	Reason    string `json:"reason"`
}

type Payload struct {
	Path    string `json:"path"`
	Content `json:"content"`
}

type Config struct {
	Address    string
	Scheme     string
	HTTPClient *http.Client
}

type API struct {
	config Config
}

//stashes uri

const StashesURI string = "/stashes"

// DefaultConfig

func DefaultConfig() *Config {
	config := &Config{
		Scheme:     "http",
		Address:    "127.0.0.1:4567",
		HTTPClient: http.DefaultClient,
	}
	return config
}

// Generic GET Request. Decoded JSON is set in the out interface{} passed in.

func (c *API) getStashes() {
	method := "GET"
	var out []Payload
	url := &url.URL{
		Scheme: c.config.Scheme,
		Host:   c.config.Address,
		Path:   StashesURI,
	}

	log.Println("No option spefied getting all Stashes")
	req, err := http.NewRequest(method, url.String(), nil)
	resp, err := c.config.HTTPClient.Do(req)

	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode != 404 {

		response, err := ioutil.ReadAll(resp.Body)

		resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
		if err := json.Unmarshal(response, &out); err != nil {
			log.Println("error : ", err)
		}
		data, err := json.MarshalIndent(out, "", " ")
		if err != nil {
			log.Println("error : ", err)
		}
		data_out := append(data, '\n')
		os.Stdout.Write(data_out)
	} else {
		log.Println("Not valid request")
	}
}

func (c *API) getStash(uri string) {
	method := "GET"
	var out Payload

	s := []string{StashesURI, uri}
	path := strings.Join(s, "/")
	url := &url.URL{
		Scheme: c.config.Scheme,
		Host:   c.config.Address,
		Path:   path,
	}
	log.Println("Getting Stash for : ", path)
	request, _ := http.NewRequest(method, url.String(), nil)
	resp, err := c.config.HTTPClient.Do(request)
	if err != nil {
		log.Println("error:", err)
	}
	if resp.StatusCode != 404 {
		response, err := ioutil.ReadAll(resp.Body)

		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(response, &out); err != nil {
			log.Println("error : ", err)
		}
		data, err := json.MarshalIndent(out, "", " ")
		if err != nil {
			log.Println("error : ", err)
		}
		data_out := append(data, '\n')
		os.Stdout.Write(data_out)
	} else {
		log.Println("There is no stash for", uri)
	}
}

func (c *API) deleteStash(uri string) {
	method := "DELETE"

	s := []string{StashesURI, uri}
	path := strings.Join(s, "/")

	url := &url.URL{
		Scheme: c.config.Scheme,
		Host:   c.config.Address,
		Path:   path,
	}

	log.Println("Deleting Stash for : ", path)
	request, _ := http.NewRequest(method, url.String(), nil)
	resp, err := c.config.HTTPClient.Do(request)
	if err != nil {
		log.Println("error :", err)
	}
	if resp.StatusCode != 404 {

		out, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		log.Printf("%s", out)
	} else {
		log.Println("No stash exists for : ", uri)
	}
}

func (c *API) createStash(path string, payload interface{}) {
	var out interface{}
	method := "POST"
	s := []string{StashesURI, path}
	uri := strings.Join(s, "/")

	url := &url.URL{
		Scheme: c.config.Scheme,
		Host:   c.config.Address,
		Path:   uri,
	}

	// Encode payload struct into JSON and create a reader for it
	encodedPayload, err := json.Marshal(payload)
	payloadReader := bytes.NewReader(encodedPayload)
	if err != nil {
		log.Println("error :", err)
	}
	log.Println("Creating Stash for : ", path)
	req, err := http.NewRequest(method, url.String(), payloadReader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		log.Println("error: ", err)
	}
	if resp.StatusCode != 404 {

		response, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err := json.Unmarshal(response, &out); err != nil {
			log.Println("error :", err)
		}
		data, err := json.MarshalIndent(out, "", " ")
		if err != nil {
			log.Println("error :", err)
		}
		data_out := append(data, '\n')
		log.Println("Created Stash for : ")
		os.Stdout.Write(data_out)
	} else {
		log.Println("Not valid request")
	}

}

func main() {

	path := flag.String("path", "", "HOSTNAME/CHECK_NAME")
	source := flag.String("source", "API", "Name of the source")
	reason := flag.String("reason", "", "Reason for creating the stash")
	uname := flag.String("uname", "toddlers@example.com", "Name of the user creating stash")
	create := flag.Bool("create", false, "Create Stash")
	delete := flag.Bool("delete", false, "Delete Stash")
	gets := flag.Bool("gets", false, "Get Stash")
	flag.Parse()

	defConfig := DefaultConfig()

	apiClient := &API{
		config: *defConfig,
	}

	if *create == true {

		if *uname != "" && *reason != "" && *path != "" {
			now := time.Now()

			data := Payload{
				Path: *path,
				Content: Content{
					Username:  *uname,
					Timestamp: now.Unix(),
					Source:    *source,
					Reason:    *reason,
				},
			}
			apiClient.createStash(*path, data)
		} else {
			flag.PrintDefaults()
		}
	} else if *delete == true {
		apiClient.deleteStash(*path)
	} else if *gets == true {
		apiClient.getStash(*path)
	} else {
		apiClient.getStashes()
	}
}
