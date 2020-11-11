package config

import (
    "context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "gopkg.in/yaml.v3"
    "html/template"
    "io/ioutil"
    "strings"
)

type Service struct {
    Name string `yaml:"name"`
    Url  string `yaml:"url"`
    Description string `yaml:"description"`
}

type serviceTemplate struct {
    Name string
    Url template.URL
    Description string
    Running bool
}

// config holds the result of parsing our config.yaml file
type config struct {
    Services map[string]*Service `yaml:"services"`
}

type ConfigTemplate struct {
    Services map[string]*serviceTemplate
}

func containerRunning(c string, containers []types.Container) bool {
    for _, container := range containers {
        if c == strings.Trim(container.Names[0], "/") {
            return true
        }
    }
    return false
}

func ReadConfig(f string) (*ConfigTemplate, error) {
    d, err := ioutil.ReadFile(f)
    if err != nil {
        return nil, err
    }

    c := &config{}
    err = yaml.Unmarshal(d, &c)
    if err != nil {
        return nil, err
    }

    ct := &ConfigTemplate{
        Services: make(map[string]*serviceTemplate),
    }

    rcs, err := getRunningContainers()
    if err != nil {
        return nil, err
    }

    for containerName := range c.Services {
        name := strings.TrimSpace(c.Services[containerName].Name)
        url := c.Services[containerName].Url
        desc := strings.TrimSpace(c.Services[containerName].Description)
        ct.Services[containerName] = &serviceTemplate{
            Name: name,
            Url: template.URL(url),
            Description: desc,
            Running: containerRunning(containerName, rcs),
        }
    }

    return ct, nil
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
