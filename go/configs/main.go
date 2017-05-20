package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Cluster struct {
	Name       string   `yaml:"name"`
	DataCentre string   `yaml:"datacenre"`
	Nodes      []string `yaml:"nodes"`
}

type Configuration struct {
	Clusters    []Cluster `yaml:"clusters"`
	MinReplicas int       `yaml:"min_replicas"`
	MaxReplicas int       `yaml:"max_replicas"`
}

func saveConfig(c Configuration, filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0644)

}

func loadConfig(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}

	var c Configuration
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}

	return c, nil

}

func main() {
	err := saveConfig(createMockConfig(), "config.yaml")
	if err != nil {
		panic(err)
	}
	c, err := loadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("+%v\n", c)

}

func createMockConfig() Configuration {
	return Configuration{
		Clusters: []Cluster{
			Cluster{
				Name:       "Dev",
				DataCentre: "Local",
				Nodes:      []string{"dev1.company.com", "dev2.company.com"},
			},
			Cluster{
				Name:       "Prod",
				DataCentre: "Amazon",
				Nodes:      []string{"prod1.company.com", "prod2.company.com", "prod3.company.com"},
			},
		},
		MinReplicas: 1,
		MaxReplicas: 5,
	}
}
