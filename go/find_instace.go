package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
	"os"
)

func getSession() *ec2.EC2 {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("ap-southeast-1")})
	if err != nil {
		fmt.Println("failed to create session", err)
	}
	return ec2.New(sess)
}

func instanceId(ipadd string) string {
	svc := getSession()
	params := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("addresses.private-ip-address"),
				Values: []*string{
					aws.String(ipadd),
				},
			},
		},
	}

	resp, err := svc.DescribeNetworkInterfaces(params)
	if err != nil {
		panic(err)
	}
	return *resp.NetworkInterfaces[0].Attachment.InstanceId
}

func findInstance(ipaddr string) {
	svc := getSession()
	iid := instanceId(ipaddr)
	ec2params := &ec2.DescribeInstancesInput{

		InstanceIds: []*string{
			aws.String(iid),
		},
	}
	instanceInfo, err := svc.DescribeInstances(ec2params)
	if err != nil {
		panic(err)
	}
	fmt.Println(instanceInfo.Reservations[0].Instances[0].Tags)
}

//func find_elb()

func main() {
	var ipad string

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "ipaddr,i",
			Usage:       "IP address of the instance",
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

}
