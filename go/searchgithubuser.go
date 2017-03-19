package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	apiURL       = "https://api.github.com"
	userEndpoint = "/users/"
)

type User struct {
	Login             string      `json:"login"`
	ID                int         `json:"id"`
	AvatarURL         string      `json:"avatar_url"`
	GravatarID        string      `json:"gravatar_id"`
	URL               string      `json:"url"`
	HTMLURL           string      `json:"html_url"`
	FollowersURL      string      `json:"followers_url"`
	FollowingURL      string      `json:"following_url"`
	GistsURL          string      `json:"gists_url"`
	StarredURL        string      `json:"starred_url"`
	SubscriptionsURL  string      `json:"subscriptions_url"`
	OrganizationsURL  string      `json:"organizations_url"`
	ReposURL          string      `json:"repos_url"`
	EventsURL         string      `json:"events_url"`
	ReceivedEventsURL string      `json:"received_events_url"`
	Type              string      `json:"type"`
	SiteAdmin         bool        `json:"site_admin"`
	Name              string      `json:"name"`
	Company           string      `json:"company"`
	Blog              string      `json:"blog"`
	Location          string      `json:"location"`
	Email             string      `json:"email"`
	Hireable          interface{} `json:"hireable"`
	Bio               string      `json:"bio"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

func getUsers(name string) User {
	resp, err := http.Get(apiURL + userEndpoint + name)
	if err != nil {
		log.Fatalf("Error receiving data: %s\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading data : %s\n", err)
	}
	var user User
	json.Unmarshal(body, &user)
	return user
}

var (
	user  string
	emoji string
)

func init() {
	flag.StringVarP(&user, "user", "u", "", "Search Users")
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	// parse flags
	flag.Parse()

	// if user does not supply flags , print usage
	// we can clean this up later by putting this
	// into its own function
	if flag.NFlag() == 0 {
		printUsage()
	}
	users := strings.Split(user, ",")
	fmt.Printf("Search user(s): %s\n", users)
	for _, u := range users {
		result := getUsers(u)
		fmt.Println(`Username:	`, result.Login)
		fmt.Println(`Name:		`, result.Name)
		fmt.Println(`Email:		`, result.Email)
		fmt.Println(`Bio:		`, result.Bio)
		fmt.Println("")
	}
