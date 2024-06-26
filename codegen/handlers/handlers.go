package handlers

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenHandlersPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenHandlersPlugin) Init() error {
	// Initialize Files
	g.Files = map[string]func() string{
		"app/handlers/frontend/auth.go": func() string {
			return `package handlers

import (
    "net/http"
)

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
		"app/handlers/frontend/views.go": func() string {
			return `package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))
`
		},
		"app/handlers/frontend/landing.go": func() string {
			return `package handlers

import (
    "net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "home.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
`
		},
		"app/handlers/frontend/about.go": func() string {
			return `package handlers

import (
    "net/http"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "about.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
`
		},

		"app/handlers/backend/auth.go": func() string {
			return `package handlers

import (
    "net/http"
)

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
		"app/handlers/backend/views.go": func() string {
			return `package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))
`
		},
		"app/handlers/backend/landing.go": func() string {
			return `package handlers

import (
    "net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "home.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
`
		},
		"app/handlers/backend/about.go": func() string {
			return `package handlers

import (
    "net/http"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "about.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
`
		},

		"app/handlers/api/api.go": func() string { return `` },
	}
	return nil

}

func (g *GenHandlersPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenHandlersPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenHandlersPlugin) Name() string {
	return "GenHandlersPlugin"
}

func (g *GenHandlersPlugin) Version() string {
	return "1.0.0"
}

func (g *GenHandlersPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenHandlersPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenHandlersPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenHandlersPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenHandlersPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/handlers"
}

func (g *GenHandlersPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenHandlersPlugin(data config.ProjectData) *GenHandlersPlugin {
	return &GenHandlersPlugin{
		Data: data,
	}
}
