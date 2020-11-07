package main

import (
    "fmt"
    "log"
    "net/http"
    "pi-server-manager/config"
    "pi-server-manager/templating"
)

const PORT = ":8000"

func main() {
    c, _ := config.ReadConfig("config.yaml")
    for _, service := range c.Services {
        fmt.Println("Url:", service.Url, "Name:", service.Name)
    }

    http.HandleFunc("/", templating.Index)

    fmt.Println("Server starting on port", PORT)
    log.Fatal(http.ListenAndServe(PORT, nil))
}
