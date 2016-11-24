package main

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"os"
	"strings"
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

	svc := route53.New(sess)
	params := &route53.ListHostedZonesInput{}
	resp, err := svc.ListHostedZones(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	f, err := os.Create("route53_dump.txt")
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	for i := 0; i < len(resp.HostedZones); i++ {
		params := &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(strings.Split(*resp.HostedZones[i].Id, "/")[2]),
		}
		rr, err := svc.ListResourceRecordSets(params)
		check(err)
		for i := 2; i < len(rr.ResourceRecordSets); i++ {
			if rr.ResourceRecordSets[i].ResourceRecords != nil {
				_, err := fmt.Fprintf(w, "%s =>  %s\n", *rr.ResourceRecordSets[i].Name, *rr.ResourceRecordSets[i].ResourceRecords[0].Value)
				check(err)
			} else {
				_, err := fmt.Fprintf(w, "%s =>  %s\n", *rr.ResourceRecordSets[i].Name, *rr.ResourceRecordSets[i].AliasTarget.DNSName)
				check(err)
			}
		}
	}

}
