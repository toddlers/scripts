package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	flag "github.com/ogier/pflag"
	"log"
	"os"
	"unsafe"
)

var arn string

func init() {
	flag.StringVarP(&arn, "arn", "a", "", "IAM ARN")
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

func check(err error) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Get error details
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		}
	}
}

func getSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	check(err)
	return sess
}


func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		printUsage()
	}

	sess := getSession()
	svc := opsworks.New(sess)
	params := &opsworks.DescribeStacksInput{}
	resp, err := svc.DescribeStacks(params)
	check(err)

	for i := 0; i < len(resp.Stacks); i++ {
		params := &opsworks.SetPermissionInput{
			IamUserArn: aws.String(arn),
			StackId:    aws.String(*resp.Stacks[i].StackId),
			AllowSsh:   aws.Bool(true),
			AllowSudo:  aws.Bool(true),
		}
		setPerms, err := svc.SetPermission(params)
		check(err)
		if unsafe.Sizeof(setPerms) == 8 {
			fmt.Printf("Added sudo permissions to stack : %s\n", *resp.Stacks[i].Name)
		}
	}
}
