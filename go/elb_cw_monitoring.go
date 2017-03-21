package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	if err != nil {
		fmt.Println("failed to create session", err)
	}
	return sess
}

func addMonitoring(elb string) {
	AName := elb + " HealthyHostCount"
	sess := getSession()
	svc := cloudwatch.New(sess)

	params := &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(AName),
		ComparisonOperator: aws.String("LessThanOrEqualToThreshold"),
		EvaluationPeriods:  aws.Int64(1),
		MetricName:         aws.String("HealthyHostCount"),
		Namespace:          aws.String(NS),
		Period:             aws.Int64(60),
		Threshold:          aws.Float64(0),
		ActionsEnabled:     aws.Bool(true),
		AlarmActions: []*string{
			aws.String(SNS),
		},
		AlarmDescription: aws.String("CQA Alerts Created From API"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("LoadBalancerName"),
				Value: aws.String(elb),
			},
		},
		OKActions: []*string{
			aws.String(SNS),
		},
		Statistic: aws.String("Minimum"),
	}
	_, err := svc.PutMetricAlarm(params)
	check(err)
	fmt.Println("Added monitoring on", elb)
}

var (
	filename string
	elbName  string
)

func init() {
	flag.StringVarP(&filename, "fname", "f", "", "File Name containing elb names")
	flag.StringVarP(&elbName, "elb", "e", "", "ELB name")
}

const (
	NS  = "AWS/ELB"
	SNS = "arn:aws:sns:us-east-1:123231312:abc"
)

func main() {
	flag.Parse()

	// safe checking
	if filename != "" && elbName != "" {
		fmt.Println("Both options can not be used simultaneously")
	} else if filename != "" {
		data, err := ioutil.ReadFile(filename)
		check(err)
		for _, elb := range strings.Split(string(data), "\n") {
			if elb == "" {
				continue
			} else {
				addMonitoring(elb)
			}
		}
	} else if elbName != "" {
		addMonitoring(elbName)
	} else {
		fmt.Println("Please provide at least one option")
		os.Exit(-1)
	}
}
