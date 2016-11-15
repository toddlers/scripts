package main

import (
	//	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	//	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := opsworks.New(sess, &aws.Config{Region: aws.String("us-east-1")})
	params := &opsworks.DescribeStacksInput{
	//StackIds: []*string{
	//	aws.String("<STACK_ID>"),
	//},
	}
	resp, err := svc.DescribeStacks(params)
	//resp, err := svc.DescribeStacks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//f, err := os.Create("stack_name.txt")
	//check(err)
	//defer f.Close()
	//w := bufio.NewWriter(f)
	for i := 0; i < len(resp.Stacks); i++ {
		params := &opsworks.SetPermissionInput{
			IamUserArn: aws.String("IAM_ARN"), // Required
			StackId:    aws.String(*resp.Stacks[i].StackId),                           // Required
			AllowSsh:   aws.Bool(true),
			AllowSudo:  aws.Bool(true),
			//Level:      aws.String("iam_only"),
		}

		resp, err := svc.SetPermission(params)
		check(err)
		fmt.Println(resp)
		//_, err := fmt.Fprintf(w, " %s => %s\n", *resp.Stacks[i].Name, *resp.Stacks[i].StackId)
		//check(err)
		//	w.Flush()
	}
	//fmt.Println(resp)
}
