package system

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Application stores application resources
type Application struct {
	Configuration *Configuration
	Template      *template.Template
	SockJS        *http.Handler
}

// CurrentApplication contains data about current application
var CurrentApplication = &Application{}

// Init reads and parses configuration file
func Init(filename *string) (err error) {
	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		return
	}

	CurrentApplication.Configuration = &Configuration{}

	err = json.Unmarshal(data, &CurrentApplication.Configuration)
	if err != nil {
		return
	}

	return
}

// LoadTemplates loads templates
func LoadTemplates() (err error) {
	var templates []string

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			templates = append(templates, path)
		}
		return nil
	}

	err = filepath.Walk(CurrentApplication.Configuration.TemplatePath, fn)
	if err != nil {
		return err
	}

	CurrentApplication.Template = template.Must(template.ParseFiles(templates...))
	return nil
}
