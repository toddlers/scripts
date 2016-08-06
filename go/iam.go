package main

import (
	"fmt"
	//	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session", err)
		return
	}

	svc := iam.New(sess)

	//params := &iam.ListAccessKeysInput{
	//	MaxItems: aws.Int64(1000),
	//}

	//resp, err := svc.ListAccessKeys(params)
	resp, err := svc.ListAccessKeys(nil)

	if err != nil {
		log.Fatal(err)
		return
	}
	var keyToFind string

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "access_key,k",
			Usage:       "iam.go --access_key <ACCESS_KEY>",
			Destination: &keyToFind,
		},
	}

	app.Action = func(c *cli.Context) error {
		matchFound := false
		uname := ""
		for _, v := range resp.AccessKeyMetadata {
			if *v.AccessKeyId == keyToFind {
				matchFound = true
				uname = *v.UserName

			}
		}
		if matchFound {
			fmt.Println("Key Belongs to username  : ", uname)
		} else {
			fmt.Println("No Match found , please try another account")
		}
		return nil
	}
	app.Run(os.Args)
}
