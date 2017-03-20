package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"log"
	"unsafe"
)

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

const (
	iamArn = "arn:aws:iam::503805144375:user/suresh.prajapati"
)

func main() {
	sess := getSession()
	svc := opsworks.New(sess)
	params := &opsworks.DescribeStacksInput{}
	resp, err := svc.DescribeStacks(params)
	check(err)

	for i := 0; i < len(resp.Stacks); i++ {
		params := &opsworks.SetPermissionInput{
			IamUserArn: aws.String(iamArn),
			StackId:    aws.String(*resp.Stacks[i].StackId),
			AllowSsh:   aws.Bool(true),
			AllowSudo:  aws.Bool(true),
		}
		setPerms, err := svc.SetPermission(params)
		check(err)
		if unsafe.Sizeof(setPerms) == 8 {
			fmt.Printf("Added sudo permissions to stack : %s", *resp.Stacks[i].Name)
		}
	}
}
