package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
)

type IPInput struct {
	GroupId    *string
	IpProtocol *string
	FromPort   *int64
	ToPort     *int64
	CidrIp     *string
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

func getInstanceID(svc *ec2.EC2, ipadd string) string {
	params := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("association.public-ip"),
				Values: []*string{
					aws.String(ipadd),
				},
			},
		},
	}
	resp, err := svc.DescribeNetworkInterfaces(params)
	BackOffWaitIfError(err, "aws-api")
	log.Println("Instance ID: ", *resp.NetworkInterfaces[0].Attachment.InstanceId)
	return *resp.NetworkInterfaces[0].Attachment.InstanceId
}

func describeSG(svc *ec2.EC2, ipadd string) []string {
	iid := getInstanceID(svc, ipadd)
	var sgid []string
	ec2params := &ec2.DescribeInstancesInput{

		InstanceIds: []*string{
			aws.String(iid),
		},
	}
	instanceInfo, err := svc.DescribeInstances(ec2params)
	if err != nil {
		check(err)
	}

	for _, group := range instanceInfo.Reservations[0].Instances[0].SecurityGroups {
		sgid = append(sgid, *group.GroupId)
	}
	log.Println("Security groups attached to instance :", sgid)
	return sgid

}

func whitelistIP(client *ec2.EC2, input IPInput, jip string) {
	var groups []string
	groups = describeSG(client, jip)
	for _, group := range groups {
		log.Printf("Trying to whitelist IP: %s in SG: %s", *input.CidrIp, group)
		_, err := client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:    &group,
			IpProtocol: input.IpProtocol,
			FromPort:   input.FromPort,
			ToPort:     input.ToPort,
			CidrIp:     input.CidrIp,
		})

		// be idempotent, i.e. skip error if this permission already exists in group
		if err != nil {
			if err.(awserr.Error).Code() != "InvalidPermission.Duplicate" {
				check(err)
			} else {
				continue
			}
		} else {
			log.Printf("Whitelisted IP: %s in SG: %s", *input.CidrIp, group)
			return
		}

	}
}

func revokeWhitelisting(client *ec2.EC2, input IPInput, jip string) {
	var groups []string
	groups = describeSG(client, jip)
	for _, group := range groups {
		log.Println("Trying whitelisting in", group)
		_, err := client.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
			GroupId:    &group,
			IpProtocol: input.IpProtocol,
			FromPort:   input.FromPort,
			ToPort:     input.ToPort,
			CidrIp:     input.CidrIp,
		})

		// be idempotent, i.e. skip error if this permission already exists in group
		if err != nil {
			if err.(awserr.Error).Code() != "InvalidPermission.NotFound" {
				check(err)
			} else {
				continue
			}
		} else {
			log.Printf("Revoked IP: %s from SG: %s", *input.CidrIp, group)
			return
		}
	}
}

func main() {

	sip := flag.String("sip", "", "source IP address")
	port := flag.Int("port", 22, "Port number to allow")
	jip := flag.String("jip", "", "Jump Server IP address")
	clean := flag.Bool("clean", false, "Clean Security Groups")
	proto := flag.String("proto", "tcp", "Protocol for whitelising")

	sess := getSession()
	svc := ec2.New(sess)

	flag.Parse()

	var myIPAddr string

	if *sip == "" {
		ipaddr, _ := getIP()
		myIPAddr = ipaddr + "/32"
	}

	inp := IPInput{
		IpProtocol: proto,
		FromPort:   aws.Int64(int64(*port)),
		ToPort:     aws.Int64(int64(*port)),
		CidrIp:     &myIPAddr,
	}
	if *clean == true {
		log.Println("Revoking SSH for: ", myIPAddr)
		revokeWhitelisting(svc, inp, *jip)
	} else {
		log.Println("Adding SSH for: ", myIPAddr)
		whitelistIP(svc, inp, *jip)
	}
}
