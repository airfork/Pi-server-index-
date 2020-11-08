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

type service struct {
    Name string `yaml:"name"`
    Url  string `yaml:"url"`
    Description string `yaml:"description"`
}

type serviceTemplate struct {
    Name string
    Url template.URL
    Description string
}

// config holds the result of parsing our config.yaml file
type config struct {
    Services map[string]service `yaml:"services"`
}

type configTemplate struct {
    Services map[string]serviceTemplate
    Nav bool
}

func containerRunning(c string, containers []types.Container) bool {
    for _, container := range containers {
        if c == strings.Trim(container.Names[0], "/") {
            return true
        }
    }
    return false
}

func ReadConfig(f string) (*configTemplate, error) {
    d, err := ioutil.ReadFile(f)
    if err != nil {
        return nil, err
    }

    c := &config{}
    err = yaml.Unmarshal(d, &c)
    if err != nil {
        return nil, err
    }

    containers, err := getRunningContainers()
    if err != nil {
        return nil, err
    }

    ct := &configTemplate{
        Services: make(map[string]serviceTemplate),
        Nav: true,
    }

    for containerName := range c.Services {
        if containerRunning(containerName, containers) {
            name := strings.TrimSpace(c.Services[containerName].Name)
            url := c.Services[containerName].Url
            desc := strings.TrimSpace(c.Services[containerName].Description)
            ct.Services[containerName] = serviceTemplate{
                Name: name,
                Url: template.URL(url),
                Description: desc,
            }
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
