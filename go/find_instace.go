package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/urfave/cli"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("ap-southeast-1")})
	check(err)
	return sess
}
func searchInstance(iid string) {
	elbs := elbInfo("")
	for _, l := range elbs.LoadBalancerDescriptions {
		for _, i := range l.Instances {
			if *i.InstanceId == iid {
				fmt.Println("Load Balancer Name : ", *l.LoadBalancerName)
			}
		}
	}
}

func elbInfo(marker string) elb.DescribeLoadBalancersOutput {
	var params *elb.DescribeLoadBalancersInput
	sess := getSession()
	svc := elb.New(sess)
	if len(marker) > 0 {
		params = &elb.DescribeLoadBalancersInput{
			Marker: aws.String(marker),
		}
	} else {
		params = &elb.DescribeLoadBalancersInput{}
	}
	resp, err := svc.DescribeLoadBalancers(params)
	check(err)
	return *resp
}

func instanceId(svc *ec2.EC2, ipadd string) string {
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
	check(err)
	return *resp.NetworkInterfaces[0].Attachment.InstanceId
}

func findInstance(ipaddr string) {
	sess := getSession()
	svc := ec2.New(sess)
	iid := instanceId(svc, ipaddr)
	ec2params := &ec2.DescribeInstancesInput{

		InstanceIds: []*string{
			aws.String(iid),
		},
	}
	instanceInfo, err := svc.DescribeInstances(ec2params)
	check(err)
	fmt.Println(instanceInfo.Reservations[0].Instances[0].Tags)
	searchInstance(iid)
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
