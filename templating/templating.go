package templating

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"pi-server-manager/config"
	"sort"
)

type T struct {
	Services *config.Config
}

type serviceOutput struct {
	Name          string
	Url           template.URL
	Description   string
	Running       bool
	ContainerName string
}

func (t T) Index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/layout.html",
		"templates/index.html",
		"templates/navbar.html",
	))

	data := make(map[string]interface{})
	data["Config"] = getServiceOutputs(t.Services)
	if os.Getenv("PI_DEV") == "" {
		data["ws"] = "wss://tunjicus.com/socket"
	} else {
		data["ws"] = "ws://" + r.Host + "/socket"
	}
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

func getServiceOutputs(c *config.Config) []serviceOutput {
	out := make([]serviceOutput, 0)
	for key, value := range c.Services {
		so := serviceOutput{
			Name:          value.Name,
			Description:   value.Description,
			Url:           template.URL(value.Url),
			Running:       value.Running,
			ContainerName: key,
		}

		out = append(out, so)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})

	return out
}
