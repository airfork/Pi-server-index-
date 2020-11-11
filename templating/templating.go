package templating

import (
    "fmt"
    "html/template"
    "net/http"
    "pi-server-manager/config"
)

type T struct {
    Services *config.ConfigTemplate
    Host string
}

func (t T) Index(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles(
        "templates/layout.html",
            "templates/index.html",
            "templates/navbar.html",
        ))

    data := make(map[string]interface{})
    data["Config"] = t.Services
    data["ws"] = "ws://"+t.Host+"/socket"
    data["Nav"] = true

    err := tmpl.ExecuteTemplate(w, "layout", data)
    if err != nil {
        fmt.Println(err)
        _, _ = w.Write([]byte("Error parsing template"))
    }
}

func (t T) Info(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles(
    "templates/layout.html",
            "templates/info.html",
            "templates/navbar.html",
    ))

    data := make(map[string]bool)
    data["Nav"] = true
    err := tmpl.ExecuteTemplate(w, "layout", data)
    if err != nil {
        _, _ = w.Write([]byte("Error parsing template"))
    }
}