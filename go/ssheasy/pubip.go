package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
)

//query the OpenDNS servers for public ip

func getIP() (string, error) {
	config := dns.ClientConfig{Servers: []string{"208.67.220.220", "208.67.222.222"}, Port: "53"}
	dnsClient := new(dns.Client)
	message := new(dns.Msg)
	message.SetQuestion("myip.opendns.com.", dns.TypeA)
	message.RecursionDesired = false
	return doDNSLookup(config, dnsClient, message)
}

func doDNSLookup(config dns.ClientConfig, client *dns.Client, message *dns.Msg) (string, error) {
	var err error
	for _, server := range config.Servers {
		serverAddr := fmt.Sprintf("%s:%s", server, config.Port)
		response, _, cliErr := client.Exchange(message, serverAddr)
		if cliErr != nil {
			log.Printf("Error on DNS lookup: %s", cliErr)
			return "", cliErr
		}
		if response.Rcode != dns.RcodeSuccess {
			err = fmt.Errorf("DNS call not successful, Response code :%d", response.Rcode)
			log.Printf(err.Error())
		} else {
			for _, answer := range response.Answer {
				if aRecord, ok := answer.(*dns.A); ok {
					return aRecord.A.String(), nil
				}
			}
		}
	}
	return "", err
}
