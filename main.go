package main

import (
    "fmt"
    "github.com/docker/docker/client"
    "github.com/go-chi/chi"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "os"
    "path/filepath"
    config2 "pi-server-manager/config"
    "pi-server-manager/helpers"
    "pi-server-manager/socket"
    "pi-server-manager/templating"
)

type dClient struct {
    cli *client.Client
}

const PORT = ":8000"
var upgrader = websocket.Upgrader{}

func main() {
    cli, err := client.NewEnvClient()
    if err != nil {
        panic("Failed to connect to docker client: " + err.Error())
    }

    services, err := config2.ReadConfig("config.yaml")
    if err != nil {
        panic("Failed to parse yaml config: " + err.Error())
    }

    con := socket.Con{Cli: cli, Services: services}

    tmpls := templating.T{Services: services, Host: "localhost" + PORT}
    if os.Getenv("PI_PROD") == "" {
        tmpls.Host = "tunjicus.com"
    }
    sc, err := helpers.SetStopCode()
    if err != nil {
        panic("Failed to set stop code: " + err.Error())
    }

    r := chi.NewRouter()

    r.Get("/", tmpls.Index)
    r.Get("/info", tmpls.Info)
    r.Get("/socket", con.ContainerListener)
    r.Get("/" + sc, helpers.StopSite)

    workDir, _ := os.Getwd()
    filesDir := http.Dir(filepath.Join(workDir, "assets"))
    helpers.FileServer(r, "/assets", filesDir)

    fmt.Println("Server starting on port", PORT)
    log.Fatal(http.ListenAndServe(PORT, r))
}
