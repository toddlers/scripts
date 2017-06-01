package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type mtConfig struct {
	marathonEndpoint string
}

type mesosConfig struct {
	mesosEndpoint string
}

type sandboxOffset struct {
	data   string `json:"data"`
	offset int    `json:"offset"`
}

type marathonAppData struct {
	App struct {
		ID    string `json:"id"`
		Tasks []struct {
			ID      string `json:"id"`
			SlaveID string `json:"slaveId"`
			Host    string `json:"host"`
		} `json:"tasks"`
	} `json:"app"`
}
type mesosAppData struct {
	Frameworks []struct {
		Executors []struct {
			Container string `json:"container"`
			Directory string `json:"directory"`
			ID        string `json:"id"`
			Tasks     []struct {
				ExecutorID  string `json:"executor_id"`
				FrameworkID string `json:"framework_id"`
				ID          string `json:"id"`
				Name        string `json:"name"`
				State       string `json:"state"`
			} `json:"tasks"`
		} `json:"executors"`
		Hostname string `json:"hostname"`
		ID       string `json:"id"`
	} `json:"frameworks"`
}

func getUrl(url string) (*http.Response, error) {
	httpClient := NewTimeoutClient(1000*time.Millisecond, 2*time.Second)
	return httpClient.Get(url)
}

func getSandboxOffset(url string) int {
	var sbox sandboxOffset
	sboxUrlResp, err := getUrl(url)
	if err != nil {
		log.Printf("Not able to fetch the log url (%s) : %s", url, err)
		os.Exit(1)
	}
	response, err := ioutil.ReadAll(sboxUrlResp.Body)
	sboxUrlResp.Body.Close()
	if err != nil {
		log.Println("Not able to read from response")
	}
	if err := json.Unmarshal(response, &sbox); err != nil {
		log.Println("Error while parsing response: ", err)
	}
	return sbox.offset

}

func (ms mesosConfig) getMesosSlaveApps(mad marathonAppData) []string {
	var mesosad mesosAppData
	var logUrls []string
	mesosSlaveResp, err := getUrl(ms.mesosEndpoint)
	if err != nil || mesosSlaveResp.StatusCode != 200 {
		log.Println("Error : ", err)
		log.Println("Not able to fetch mesos slave state")
		os.Exit(1)
	}
	response, err := ioutil.ReadAll(mesosSlaveResp.Body)
	mesosSlaveResp.Body.Close()
	if err != nil {
		fmt.Println("Not able to readg mseos slave response")
	}
	if err := json.Unmarshal(response, &mesosad); err != nil {
		log.Println("Error :", err)
	}
	//log.Println("stuff")
	slaveHostname := fmt.Sprintf("%s:5051", mad.App.Tasks[0].Host)
	//slaveHostname := "localhost:9090"
	//log.Println(mesosad.Frameworks[0].Hostname)
	for _, f := range mesosad.Frameworks {
		for _, e := range f.Executors {
			if e.ID == mad.App.Tasks[0].ID {
				logUrls = append(logUrls, fmt.Sprintf("http://%s/files/read?path=%s/stdout", slaveHostname, e.Directory))
				logUrls = append(logUrls, fmt.Sprintf("http://%s/files/read?path=%s/stderr", slaveHostname, e.Directory))
			}
		}
	}
	return logUrls
}

func (mt mtConfig) getMarathonApps() marathonAppData {
	var mad marathonAppData

	marathonResp, err := getUrl(mt.marathonEndpoint)
	if err != nil || marathonResp.StatusCode != 200 {
		log.Printf("Not able to fetch marathon app data : ", err)
	}
	response, err := ioutil.ReadAll(marathonResp.Body)
	marathonResp.Body.Close()
	if err != nil {
		fmt.Println("Not able to read marathon response data")
	}
	if err := json.Unmarshal(response, &mad); err != nil {
		log.Println("Error :", err)
	}
	return mad
}

func main() {
	appName := flag.String("appName", "", "App Name")
	marathonUrl := flag.String("marathonUrl", "http://marathon1-123456.us-east-1.elb.amazonaws.com/v2/apps/", "Marathon Endpoint")

	flag.Parse()

	if *appName == "" {
		fmt.Println("Please provide an app name")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	var initialOffset int
	var containerLogs interface{}

	marathonConfig := mtConfig{marathonEndpoint: fmt.Sprintf("%s%s", *marathonUrl, *appName)}
	mad := marathonConfig.getMarathonApps()

	mssConfig := mesosConfig{mesosEndpoint: fmt.Sprintf("http://%s:5051/state.json", mad.App.Tasks[0].Host)}
	sandboxLogUrls := mssConfig.getMesosSlaveApps(mad)
	for _, url := range sandboxLogUrls {
		sboxOffset := getSandboxOffset(url)
		if sboxOffset < 0 {
			initialOffset = 0
		} else {
			initialOffset = sboxOffset
		}
		sandboxLogUrl := fmt.Sprintf("%s&offset=%d", url, initialOffset)
		sboxLog, err := getUrl(sandboxLogUrl)
		if err != nil {
			log.Printf("Not able to fetch log URL : (%s) : %s", url, err)
			os.Exit(1)
		}
		logData, err := ioutil.ReadAll(sboxLog.Body)
		if err != nil {
			log.Println("Not able to read log response : ", err)
		}
		sboxLog.Body.Close()
		if err := json.Unmarshal(logData, &containerLogs); err != nil {
			log.Println("Not able to unmarshal logs : ", err)
		}
		fmt.Println(containerLogs)

	}
}
