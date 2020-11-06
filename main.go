package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
)

type service struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type config struct {
	Services map[string]service `yaml:"services"`
}

func containerRunning(c string, containers []types.Container) bool {
	for _, container := range containers {
		if c == strings.Trim(container.Names[0], "/") {
			return true
		}
	}
	return false
}

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	c := config{}
	d, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println(c.Services)

	runningApps := make(map[string]service)
	for containerName := range c.Services {
		if containerRunning(containerName, containers) {
			runningApps[containerName] = c.Services[containerName]
		}
	}

	fmt.Println(runningApps)
}
