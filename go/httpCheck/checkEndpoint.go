package main

import (
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"time"
)

func checkUrl(url string) (*http.Response, error) {
	log.Println("Default ConnectTimeout : ", 500*time.Millisecond)
	log.Println("Default ReadWriteTimeout : ", 1*time.Second)
	httpClient := NewTimeoutClient(500*time.Millisecond, 1*time.Second)
	return httpClient.Get(url)
}

type slackConfig struct {
	token   string
	channel string
}

func main() {
	token := flag.String("token", "", "Please provide slackbot token")
	channel := flag.String("channel", "", "Please provide channel id")
	url := flag.String("url", "", "Please provide a url to monitor")
	flag.Parse()

	if *token == "" || *channel == "" || *url == "" {
		fmt.Println("Input is missing")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	var sc slackConfig

	sc.token = *token
	sc.channel = *channel

	resp, err := checkUrl(*url)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Request Timed out")
		sc.sendMessage(*url)
	}

}

func (sConfig *slackConfig) sendMessage(url string) {
	api := slack.New(sConfig.token)
	params := slack.PostMessageParameters{Username: "SlackBot"}
	attachment := slack.Attachment{
		Title:     "Alert ! HTTP Endpoint Not Available !",
		TitleLink: url,
		Text:      "Please check endpoint",
	}
	params.Attachments = []slack.Attachment{attachment}
	message := "HTTP Endpoint Availability Check !"
	channelID, timestamp, err := api.PostMessage(sConfig.channel, message, params)
	if err != nil {
		log.Printf("%s\n", err)
	} else {
		log.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	}
}
