package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type s3Proxy struct {
	bucket  string
	config  string
	timeout time.Duration
	port    string
	region  string
}

func region() string {
	envRegion := os.Getenv("REGION")
	if len(envRegion) != 0 {
		return envRegion
	}
	log.Println("No region name specified using default: us-east-1")
	return "us-east-1"
}

func port() string {
	envPort := os.Getenv("PORT")
	if len(envPort) != 0 {
		return envPort
	}
	log.Println("No port specified using default: 8080")
	return "8080"
}

func bucketName() (string, error) {
	envBucket := os.Getenv("BUCKET")
	if len(envBucket) != 0 {
		return envBucket, nil
	} else {
		return "", errors.New("No bucket name provided")
	}
}

func configFileName() (string, error) {
	envConfigFile := os.Getenv("CONFIGFILE")
	if len(envConfigFile) != 0 {
		return envConfigFile, nil
	} else {
		return "", errors.New("No config file name provided")
	}
}

func timeout() time.Duration {
	envTimeout := os.Getenv("TIMEOUT")
	if len(envTimeout) != 0 {
		etimeout, _ := strconv.Atoi(envTimeout)
		return time.Duration(etimeout)
	} else {
		log.Println("No timeout provided using default : 10s")
		return time.Second * 10
	}
}
func check(err error) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Get error details
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		}
	}
}

func getIP(r *http.Request) string {
	tryHeader := func(key string) (string, bool) {
		if headerVal := r.Header.Get(key); len(headerVal) > 0 {
			if !strings.ContainsRune(headerVal, ',') {
				return headerVal, true
			}
			return strings.SplitN(headerVal, ",", 2)[0], true
		}
		return "", false
	}

	for _, header := range []string{"X-FORWARDED-FOR", "X-REAL-IP"} {
		if headerVal, ok := tryHeader(header); ok {
			return headerVal
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func getSession(region string) *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	check(err)
	return sess
}

func getConfig(w http.ResponseWriter, r *http.Request, s3p s3Proxy) {

	sess := getSession(s3p.region)
	svc := s3.New(sess)

	// Create a context with a timeout that will abort the download
	// if it takes more than the passed in timeout
	ctx := context.Background()

	var cancelFn func()

	if s3p.timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, s3p.timeout)
	}

	// Ensure the context is canceled to prevent leaking
	defer cancelFn()

	// Downloads the object to s3. The context will interrupt the request if
	// the timeout expires

	start := time.Now().UTC()
	result, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s3p.bucket),
		Key:    aws.String(s3p.config),
	})
	end := time.Now().UTC()
	delta := end.Sub(start)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			//if the sdk can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned
			fmt.Fprintf(os.Stderr, "download canceled due to timeout, %v\n", err)
			w.WriteHeader(http.StatusRequestTimeout)
			return
		} else {
			fmt.Fprintf(os.Stderr, "failed to download the object, %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Println("Error redaing http response body : ", err)
	}

	ip := getIP(r)

	contentLength := w.Header().Get("Content-Length")

	if len(contentLength) == 0 {
		contentLength = "cached"
	} else {
		contentLength = fmt.Sprintf("%s bytes", contentLength)
	}

	w.Write(body)

	log.Printf("%s - %s - %s - %s - %v - %s", start.Format(time.RFC3339), ip, r.Method, r.URL.Path, delta, contentLength)

}

func main() {
	var s3p s3Proxy

	if bucketName, berr := bucketName(); berr == nil {
		log.Println("Bucket Name : ", bucketName)
		s3p.bucket = bucketName
	} else {
		log.Println(berr)
		os.Exit(1)
	}

	if configName, cerr := configFileName(); cerr == nil {
		log.Println("Config File Name :", configName)
		s3p.config = configName
	} else {
		log.Println(cerr)
		os.Exit(1)
	}

	s3p.port = port()

	s3p.timeout = timeout()

	s3p.region = region()

	log.Println("Timeout Configured : ", s3p.timeout)
	log.Println("AWS Region Name : ", s3p.region)

	fmt.Printf("Config File server Listening on: %s\n", s3p.port)

	http.HandleFunc("/getConfig", func(w http.ResponseWriter, r *http.Request) {
		getConfig(w, r, s3p)
	})

	log.Fatal(http.ListenAndServe(":"+s3p.port, nil))
}
