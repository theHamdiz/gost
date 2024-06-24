package handlers

import (
	"github.com/theHamdiz/gost/config"
	"github.com/theHamdiz/gost/gen/general"
)

type Generator struct {
	Files map[string]func() string
}

func (g *Generator) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenerator() *Generator {
	return &Generator{
		Files: map[string]func() string{
			"app/handlers/auth.go": func() string {
				return `package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "signup.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "signin.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
`
			},
			"app/handlers/landing.go": func() string {
				return `package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "home.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

`
			},
			"app/handlers/about.go": func() string {
				return `package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))

func AboutHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "about.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
`
			},
		},
	}
}
