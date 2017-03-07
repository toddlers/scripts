// Given an instnace id , it will return an array of lb name 
// instance registered with

// First Attempt on recursion in golang :)

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getElbSession() *elb.ELB {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("ap-southeast-1")})
	check(err)
	return elb.New(sess)
}

func searchInstance(elbs *elb.DescribeLoadBalancersOutput) ([]string, bool) {
	var lbname []string
	found := false
	iid := "i-1a2b8fd5"
	for _, l := range elbs.LoadBalancerDescriptions {
		for _, i := range l.Instances {
			if *i.InstanceId == iid {
				lbname = append(lbname, *l.LoadBalancerName)
				found = true
			}
		}
	}
	return lbname, found
}

var lbNames []string

func elbInfo(marker string, first_call bool) []string {
	var params *elb.DescribeLoadBalancersInput

	svc := getElbSession()

	if len(marker) == 0 && !first_call {
		return lbNames
	}

	if len(marker) == 0 {
		params = &elb.DescribeLoadBalancersInput{}
	} else {
		params = &elb.DescribeLoadBalancersInput{
			Marker: aws.String(marker),
		}
	}

	resp, err := svc.DescribeLoadBalancers(params)
	check(err)
	lbname, found := searchInstance(resp)
	if found {
		lbNames = append(lbNames, lbname...)
	}
	if resp.NextMarker != nil {
		marker = *resp.NextMarker
	} else {
		marker = ""
	}
	return elbInfo(marker, false)
}

func main() {
	lbs := elbInfo("", true)
	fmt.Println(strings.Join(lbs, ","))
}
