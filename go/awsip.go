package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	IP_RANGES_URL     = "https://ip-ranges.amazonaws.com/ip-ranges.json"
	IP_RANGE_FILENAME = "ip-range.json"
)

type Prefix struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}

type IpRanges struct {
	Prefixes []Prefix
}

func main() {
	ip := flag.String("ip", "", "A comma-separated list of IP")
	flag.Parse()
	if *ip == "" {
		log.Fatalln("Need at least one IP address")
	}
	switch {
	// multiple IP addresses provided:
	case strings.Contains(*ip, ","):
		ips := strings.Split(*ip, ",")
		for _, ip := range ips {
			err := searchIp(ip)
			if err != nil {
				log.Fatalf("Can't find the mentioned ip addresses : %v", err)
			}
		}
		// a single ip address provided:
	default:
		err := searchIp(*ip)
		if err != nil {
			log.Fatalf("Can't find the mentioned ip address : %v", err)
		}
	}

}

func searchIp(ip string) error {
	var found bool
	if _, err := os.Stat(filepath.Join("/tmp/", IP_RANGE_FILENAME)); os.IsNotExist(err) {
		fmt.Println("IP ranges file doesnt exists downloading now.....")
		err = downloadIpRanges()
		if err != nil {
			log.Fatalf("Not able to fetch ip ranges : %v", err)
		}
	}
	fmt.Printf("Found Ip ranges file : %s\n\n", IP_RANGE_FILENAME)
	contents, err := ioutil.ReadFile(filepath.Join("/tmp/", IP_RANGE_FILENAME))
	if err != nil {
		return errors.Wrap(err, "Not able to read file")
	}
	var ipranges IpRanges
	_ = json.Unmarshal([]byte(contents), &ipranges)

	for _, iprange := range ipranges.Prefixes {
		_, prefix, err := net.ParseCIDR(iprange.IPPrefix)
		if err != nil {
			return errors.Wrapf(err, "Not able to parse IP Prefix : %v", iprange.IPPrefix)
		}
		if prefix.Contains(net.ParseIP(ip)) {
			found = true
			fmt.Println("IP Address bleongs to : ")
			fmt.Println("-------------------------------")
			fmt.Println("IP Prefix :\t", iprange.IPPrefix)
			fmt.Println("Service Name :\t", iprange.Service)
			fmt.Println("Region Name :\t", iprange.Region)
			fmt.Println("-------------------------------")
		}
	}
	if found != true {
		fmt.Printf("IP Address : %s not found in Ranges Published by AWS", ip)
	}
	return nil
}

func downloadIpRanges() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	req, err := http.NewRequestWithContext(ctx, "GET", IP_RANGES_URL, nil)
	client := &http.Client{}
	go func() {
		time.Sleep(time.Second * 60)
		println("Cancel")
		cancel()
	}()
	fmt.Println("Fetching the ip ranges")
	resp, err := client.Do(req)
	fmt.Println("Fetched ip ranges")
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		err := saveFile(*resp)
		if err != nil {
			return errors.Wrap(err, "Not able to save the file")
		}
	}
	return nil
}

func saveFile(resp http.Response) error {
	fmt.Println("saving the file")
	filename := filepath.Join("/tmp/", IP_RANGE_FILENAME)
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Could not able to create the file")
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.Wrap(err, "Could not able to save the file")
	}
	return nil

}
