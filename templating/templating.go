package templating

import (
    "html/template"
    "net/http"
    "pi-server-manager/config"
)

func Index(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templating/templates/index.html"))
    c, err := config.ReadConfig("config.yaml")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        _, _ = w.Write([]byte(err.Error()))
        return
    }

    err = tmpl.Execute(w, c.Services)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        _, _ = w.Write([]byte("Error parsing template"))
    }
}
