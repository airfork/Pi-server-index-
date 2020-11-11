package helpers

import (
	"github.com/go-chi/chi"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
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

func StopSite(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func SetStopCode() (string, error) {
	sc := GetStopCode()
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

func GetStopCode() string {
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
