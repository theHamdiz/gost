package ui

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenUiPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenUiPlugin) Init() error {
	// Initialize Files
	g.Files = map[string]func() string{
		"app/ui/embed.go": func() string {
			return `package ui

import (
	"embed"
	"io/fs"
	"net/http"
	{{- if eq .BackendPkg "echo" }}
	"github.com/labstack/echo/v5"
	{{- else if eq .BackendPkg "gin" }}
	"github.com/gin-gonic/gin"
	{{- else if eq .BackendPkg "chi" }}
	"github.com/go-chi/chi/v5"
	{{- end }}
)

//go:embed front/dist/* back/dist/*
var uiFS embed.FS

// RegisterRoutes -> registers the embedded static files with the chosen router (Echo, Gin, Chi, or http.ServeMux)
func RegisterRoutes(mux interface{}) {
	frontendFS, _ := fs.Sub(uiFS, "frontend/dist")
	backendFS, _ := fs.Sub(uiFS, "backend/dist")

	switch m := mux.(type) {
	{{- if eq .BackendPkg "echo" }}
	case *echo.Echo:
		m.GET("/frontend/*", echo.WrapHandler(http.StripPrefix("/frontend", http.FileServer(http.FS(frontendFS)))))
		m.GET("/backend/*", echo.WrapHandler(http.StripPrefix("/backend", http.FileServer(http.FS(backendFS)))))
	{{- else if eq .BackendPkg "gin" }}
	case *gin.Engine:
		m.StaticFS("/frontend", http.FS(frontendFS))
		m.StaticFS("/backend", http.FS(backendFS))
	{{- else if eq .BackendPkg "chi" }}
	case chi.Router:
		m.Handle("/frontend/*", http.StripPrefix("/frontend", http.FileServer(http.FS(frontendFS))))
		m.Handle("/backend/*", http.StripPrefix("/backend", http.FileServer(http.FS(backendFS))))
	{{- else }}
	case *http.ServeMux:
		m.Handle("/frontend/", http.StripPrefix("/frontend", http.FileServer(http.FS(frontendFS))))
		m.Handle("/backend/", http.StripPrefix("/backend", http.FileServer(http.FS(backendFS))))
	{{- end }}
	}
}
			`
		},
		"app/ui/frontend/package.json": func() string {
			return `{
  "name": "front",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "type": "module",
  "prettier": {
    "tabWidth": 4,
    "printWidth": 110,
  },
  "devDependencies": {
    "sass": "^1.45.0",
    "vite": "^5.0.11"
  }
}`
		},
		"app/ui/frontend/.env":           func() string { return `` },
		"app/ui/frontend/.env.dev":       func() string { return `` },
		"app/ui/frontend/robots.txt":     func() string { return `` },
		"app/ui/frontend/vite.config.js": func() string { return `` },
		"app/ui/public/index.html":       func() string { return `` },
		"app/ui/frontend/README.md":      func() string { return `` },
		"app/ui/README.md":               func() string { return `` },
		"app/ui/frontend/index.html":     func() string { return `` },
		"app/ui/frontend/pages/signin.templ": func() string {
			return `package signin
		templ Signin(){
		}
		`
		},
		"app/ui/frontend/pages/signup.templ": func() string {
			return `package signup
		templ Signup(){
		}
		`
		},
		"app/ui/frontend/assets/css/style.css": func() string {
			return ``
		},
		"app/ui/frontend/assets/js/index.js": func() string {
			return ``
		},
		"app/ui/frontend/components/index.js": func() string { return `` },
		"app/ui/frontend/pages/index.js":      func() string { return `` },
		"app/ui/frontend/store/index.js":      func() string { return `` },

		"app/ui/backend/package.json": func() string {
			return `{
  "name": "front",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "type": "module",
  "prettier": {
    "tabWidth": 4,
    "printWidth": 110,
  },
  "devDependencies": {
    "sass": "^1.45.0",
    "vite": "^5.0.11"
  }
}`
		},
		"app/ui/backend/vite.config.js": func() string { return `` },
		"app/ui/backend/.env":           func() string { return `` },
		"app/ui/backend/.env.dev":       func() string { return `` },
		"app/ui/backend/robots.txt":     func() string { return `` },
		"app/ui/backend/README.md":      func() string { return `` },
		"app/ui/backend/assets/css/style.css": func() string {
			return ``
		},
		"app/ui/backend/assets/js/index.js": func() string {
			return ``
		},
		"app/ui/backend/components/index.js": func() string { return `` },
		"app/ui/backend/pages/index.js":      func() string { return `` },
		"app/ui/backend/store/index.js":      func() string { return `` },
		"app/ui/backend/index.html":          func() string { return `` },
		"app/ui/backend/signin.html":         func() string { return `` },
		"app/ui/backend/signup.html":         func() string { return `` },
	}

	return nil
}

func (g *GenUiPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenUiPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenUiPlugin) Name() string {
	return "GenUiPlugin"
}

func (g *GenUiPlugin) Version() string {
	return "1.0.0"
}

func (g *GenUiPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenUiPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenUiPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenUiPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenUiPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/ui"
}

func (g *GenUiPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenUiPlugin(data config.ProjectData) *GenUiPlugin {
	return &GenUiPlugin{
		Data: data,
	}
}
