package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"log"
	"time"
)

func check(err error) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Get error details
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		}
	}
}

func delayMinute(n time.Duration) {
	time.Sleep(n * time.Minute)
}

func runCommand(svc *opsworks.OpsWorks, resp *opsworks.DescribeStacksOutput, start, end int) {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Println("recovering", err)
		}
	}()
	for i := start; i < end; i++ {
		fmt.Println(i)
		fmt.Println(*resp.Stacks[i].Name)

	/* support command
	deploy, undeploy, rollback, start, stop, restart, setup,
	configure, update_dependencies, install_dependencies,
	update_custom_cookbooks, execute_recipes, sync_remote_users
	*/
		pCook := &opsworks.CreateDeploymentInput{
			Command: &opsworks.DeploymentCommand{ // Required
				Name: aws.String("update_custom_cookbooks"), // Required
			},
			StackId: aws.String(*resp.Stacks[i].StackId), // Required
			Comment: aws.String("Updating ssh_users recipe to latest"),
		}
		_, err := svc.CreateDeployment(pCook)
		check(err)
	}

}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := opsworks.New(sess, &aws.Config{Region: aws.String("us-east-1")})
	params := &opsworks.DescribeStacksInput{}
	resp, err := svc.DescribeStacks(params)
	check(err)
  
/* let's not bombard gitlab, and put batching for 
update custom cookbooks
*/


	start := 0
	end := 10
	for start < len(resp.Stacks) {
		if end >= len(resp.Stacks) {
			end = len(resp.Stacks) - 1
		}
		runCommand(svc, resp, start, end)
		start += 10
		end += 10
		delayMinute(2)
	}
}
