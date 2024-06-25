package ui

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
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
			"ui/embed.go": func() string {
				return `package ui

import (
	"embed"
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

// RegisterRoutes -> registers the embedded static files with the provided router (Echo, Gin, Chi, or http.ServeMux)
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
			"ui/front/.env":     func() string { return `` },
			"ui/back/.env":      func() string { return `` },
			"ui/front/.env.dev": func() string { return `` },
			"ui/back/.env.dev":  func() string { return `` },
			"ui/.gitignore": func() string {
				return `.DS_Store
node_modules/*
/dist

# local env files
.env.local
.env.*.local
.env.dev*


# Log files
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*

# Editor directories and files
.idea
.vscode
*.suo
*.ntvs*
*.njsproj
*.sln
*.sw?
`
			},
			"ui/front/README.md":   func() string { return `` },
			"ui/back/README.md":    func() string { return `` },
			"ui/README.md":         func() string { return `` },
			"ui/front/index.html":  func() string { return `` },
			"ui/back/index.html":   func() string { return `` },
			"ui/back/signin.html":  func() string { return `` },
			"ui/front/signin.html": func() string { return `` },
			"ui/front/signup.html": func() string { return `` },
			"ui/back/signup.html":  func() string { return `` },
			"ui/front/package.json": func() string {
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
			"ui/back/package.json": func() string {
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
			"ui/back/public/style.css": func() string {
				return ``
			},
			"ui/front/public/style.css": func() string {
				return ``
			},
			"ui/back/public/index.js": func() string {
				return ``
			},
			"ui/front/public/index.js": func() string {
				return ``
			},
		},
	}
}
