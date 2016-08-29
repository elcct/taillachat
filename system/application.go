package system

import (
	"errors"
	"html/template"
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
func Init() (err error) {
	config := &Configuration{}

	publicPath := os.Getenv("TAILLA_PUBLIC_PATH")
	templatePath := os.Getenv("TAILLA_TEMPLATE_PATH")
	bind := os.Getenv("TAILLA_BIND")

	if publicPath == "" {
		return errors.New("Missing TAILLA_PUBLIC_PATH")
	}
	if templatePath == "" {
		return errors.New("Missing TAILLA_TEMPLATE_PATH")
	}
	if bind == "" {
		bind = "0.0.0.0:8000"
	}

	config.PublicPath = publicPath
	config.TemplatePath = templatePath
	config.Bind = bind

	CurrentApplication.Configuration = config

	loadTemplates()
	return
}

// loadTemplates loads templates
func loadTemplates() (err error) {
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
