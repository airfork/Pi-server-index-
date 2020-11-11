package config

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

type Service struct {
	Name        string `yaml:"name"`
	Url         string `yaml:"url"`
	Description string `yaml:"description"`
	Running bool
}


// Config holds the result of parsing our Config.yaml file
type Config struct {
	Services map[string]*Service `yaml:"services"`
}

func containerRunning(c string, containers []types.Container) bool {
	for _, container := range containers {
		if c == strings.Trim(container.Names[0], "/") {
			return true
		}
	}
	return false
}

func ReadConfig(f string) (*Config, error) {
	d, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	rcs, err := getRunningContainers()
	if err != nil {
		return nil, err
	}

	for containerName, service := range c.Services {
		service.Name = strings.TrimSpace(service.Name)
		service.Url = strings.TrimSpace(service.Url)
		service.Description = strings.TrimSpace(service.Description)
		service.Running = containerRunning(containerName, rcs)
	}

	return c, nil
}

func getRunningContainers() ([]types.Container, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}
