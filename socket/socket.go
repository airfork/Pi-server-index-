package socket

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/events"
    "github.com/docker/docker/client"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    config2 "pi-server-manager/config"
)

type Con struct {
    Cli *client.Client
    Services *config2.ConfigTemplate
}

type actionType string

const (
    responseErr actionType = "err"
    die   actionType = "die"
    start actionType = "start"
)

type errResponse struct {
    Action     actionType
    ErrMessage string
}

type dieResponse struct {
    Action    actionType
    Container string
}

type startResponse struct {
    Action    actionType
    Container string
    Info      config2.Service
}

func (con Con) ContainerListener(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{}
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()
    cli, err := client.NewEnvClient()
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return;
    }
    ctx, _ := context.WithCancel(context.Background())
    msgs, errs := cli.Events(ctx, types.EventsOptions{})
    for {
        select {
        case errData := <-errs:
            err = c.WriteMessage(websocket.TextMessage, []byte(errData.Error()))
            if err != nil {
                log.Println("write:", err)
                break
            }
        case msg := <-msgs:
            data, err := con.updatePage(&msg)
            if err != nil {
                j, err := json.Marshal(errResponse{Action: responseErr, ErrMessage: err.Error()})
                if err != nil {
                    log.Println("Error marshalling err response:", err)
                    break
                }

                err = c.WriteMessage(websocket.TextMessage, j)
                if err != nil {
                    log.Println("Error writing err msg error response", err)
                    break
                }
            }

            if data == nil {
                continue
            }

            err = c.WriteMessage(websocket.TextMessage, data)
            if err != nil {
                log.Println("write:", err)
                break
            }
        }
    }
}

func (con Con) updatePage(msg *events.Message) ([]byte, error) {

    container := msg.Actor.Attributes["name"]
    s, ok := con.Services.Services[container]
    if !ok {
       return nil, nil
    }

    switch msg.Action {
    case "die":
        s.Running = false
        j, err := json.Marshal(dieResponse{Action: die, Container: container})
        if err != nil {
            return nil, err
        }

        return j, nil
    case "start":
        s.Running = true
        service := config2.Service{
            Description: s.Description,
            Name: s.Name,
            Url: string(s.Url),
        }
        j, err := json.Marshal(startResponse{
            Action:    start,
            Container: container,
            Info:      service,
        })

        if err != nil {
            return nil, err
        }
        return j, nil
    default:
        return nil, nil
    }
}
