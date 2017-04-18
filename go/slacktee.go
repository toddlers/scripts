package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"io"
	"os"
)

func main() {
	token := flag.String("token", "", "Please provide slackbot token")
	channel := flag.String("channel", "", "Please provide channel id")
	flag.Parse()

	if *token == "" || *channel == "" {
		fmt.Println("Input is missing")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	info, _ := os.Stdin.Stat()

	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		fmt.Println("The command is intended to work with pipes")
		fmt.Println("Usage:")
		fmt.Println(" echo \"your text\" OR cat yourfile.txt| searchr -pattern=<your_pattern> -channel=<channel_id>")
	} else if info.Size() > 0 {
		reader := bufio.NewReader(os.Stdin)
		sendMessage(*token, *channel, reader)
	}
}

func sendMessage(token, channel string, reader *bufio.Reader) {
	line := 1
	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		api := slack.New(token)
		params := slack.PostMessageParameters{}
		channelID, timestamp, err := api.PostMessage(channel, input, params)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	}
	line++
}
