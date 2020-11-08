package main

import (
    "fmt"
    "github.com/go-chi/chi"
    "log"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
    "pi-server-manager/templating"
    "strings"
    "time"
)

const PORT = ":8000"

func main() {
    sc, err := setStopCode()
    if err != nil {
        panic("Failed to set stop code: " + err.Error())
    }
    r := chi.NewRouter()

    r.Get("/", templating.Index)
    r.Get("/" + sc, stopSite)

    workDir, _ := os.Getwd()
    filesDir := http.Dir(filepath.Join(workDir, "assets"))
    fileServer(r, "/assets", filesDir)

    fmt.Println("Server starting on port", PORT)
    log.Fatal(http.ListenAndServe(PORT, r))
}

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
    if strings.ContainsAny(path, "{}*") {
        panic("FileServer does not permit any URL parameters.")
    }

    if path != "/" && path[len(path)-1] != '/' {
        r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
        path += "/"
    }
    path += "*"

    r.Get(path, func(w http.ResponseWriter, r *http.Request) {
        rctx := chi.RouteContext(r.Context())
        pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
        fs := http.StripPrefix(pathPrefix, http.FileServer(root))
        fs.ServeHTTP(w, r)
    })
}

func stopSite(w http.ResponseWriter, r *http.Request) {
    os.Exit(0)
}

func setStopCode() (string, error) {
    sc := getStopCode()
    f, err := os.Create("stopCode.txt")
    if err != nil {
        return "", err
    }

    _, err = f.WriteString(sc)
    if err != nil {
        return "", err
    }
    return sc, nil
}

func getStopCode() string {
    rand.Seed(time.Now().Unix())

    //Only lowercase
    charSet := "abcdedfghijklmnopqrstuvwzyz0123456789"
    var output strings.Builder
    length := 72
    for i := 0; i < length; i++ {
        random := rand.Intn(len(charSet))
        randomChar := charSet[random]
        output.WriteString(string(randomChar))
    }
    return output.String()
}