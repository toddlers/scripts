package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
	"os"
)

func getSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("ap-southeast-1")})
	if err != nil {
		fmt.Println("failed to create session", err)
	}
	return sess
}

func findInstance(ipaddr string) {
	sess := getSession()
	svc := ec2.New(sess)
	params := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("addresses.private-ip-address"),
				Values: []*string{
					aws.String(ipaddr),
				},
			},
		},
	}
	resp, err := svc.DescribeNetworkInterfaces(params)
	if err != nil {
		panic(err)
	}
	ec2params := &ec2.DescribeInstancesInput{

		InstanceIds: []*string{
			aws.String(*resp.NetworkInterfaces[0].Attachment.InstanceId),
		},
	}
	instanceInfo, err := svc.DescribeInstances(ec2params)
	fmt.Println(instanceInfo.Reservations[0].Instances[0].Tags)
}

func main() {
	var ipad string

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "ipaddr,i",
			Usage:       "find_instance.go --ipaddr <IP_ADDR>",
			Destination: &ipad,
		},
	}

	app.Action = func(c *cli.Context) error {
		if len(ipad) > 0 {
			findInstance(ipad)
		}
		return nil
	}

	app.Run(os.Args)

	//fmt.Println(*resp.NetworkInterfaces[0].Attachment.InstanceId)
}
