package web

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
		"app/web/frontend/embed.go": func() string {
			return `package web

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
		"app/web/frontend/package.json": func() string {
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
    "printWidth": 110
  },
  "devDependencies": {
    "sass": "^1.45.0",
    "vite": "^5.0.11"
  }
}`
		},
		"app/web/frontend/.env":           func() string { return `` },
		"app/web/frontend/.env.dev":       func() string { return `` },
		"app/web/frontend/robots.txt":     func() string { return `` },
		"app/web/frontend/vite.config.js": func() string { return `` },
		"app/web/public/index.html":       func() string { return `` },
		"app/web/frontend/README.md":      func() string { return `` },
		"app/web/README.md":               func() string { return `` },
		"app/web/frontend/index.html":     func() string { return `` },
		"app/web/frontend/pages/signin.templ": func() string {
			return `package signin
		templ Signin(){
		}
		`
		},
		"app/web/frontend/pages/signup.templ": func() string {
			return `package signup
		templ Signup(){
		}
		`
		},
		"app/web/frontend/assets/css/style.css": func() string {
			return ``
		},
		"app/web/frontend/assets/js/index.js": func() string {
			return ``
		},
		"app/web/frontend/components/index.js": func() string { return `` },
		"app/web/frontend/pages/index.js":      func() string { return `` },
		"app/web/frontend/store/index.js":      func() string { return `` },
		"app/web/backend/package.json": func() string {
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
    "printWidth": 110
  },
  "devDependencies": {
    "sass": "^1.45.0",
    "vite": "^5.0.11"
  }
}`
		},
		"app/web/backend/vite.config.js": func() string { return `` },
		"app/web/backend/.env":           func() string { return `` },
		"app/web/backend/.env.dev":       func() string { return `` },
		"app/web/backend/robots.txt":     func() string { return `` },
		"app/web/backend/README.md":      func() string { return `` },
		"app/web/backend/assets/css/style.css": func() string {
			return ``
		},
		"app/web/backend/assets/js/index.js": func() string {
			return ``
		},
		"app/web/backend/components/index.js": func() string { return `` },
		"app/web/backend/pages/index.js":      func() string { return `` },
		"app/web/backend/store/index.js":      func() string { return `` },
		"app/web/backend/index.html":          func() string { return `` },
		"app/web/backend/signin.html":         func() string { return `` },
		"app/web/backend/signup.html":         func() string { return `` },
		"app/web/components/head.templ": func() string {
			return `package head

templ Head(title, css, js){
    <head>
		<title>{ title }</title>
		<link rel="icon" type="image/x-icon" href="/public/favicon.ico"/>
		<meta charset="UTF-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<link href="./app/assets/static/css/tailwind.css" rel="stylesheet">
		<link rel="stylesheet" href={ css }/>
		<script src={ js }></script>
		<!-- Alpine Plugins -->
		<script defer src="https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js"></script>
		<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
		<!-- HTMX -->
		<script src="./app/assets/static/js/htmx.min.js"></script>
	</head>
}
`
		},
		"app/web/layouts/base.templ": func() string {
			return `package layouts

import "{{.AppName}}/app/views"

templ Base(title, css, js string){
 	<!DOCTYPE html>
	<html lang="en">
		@components.Head(title, css, js)
		<body x-data="{theme: 'dark'}" :class="theme" lang="en">
			{ children... }
			@components.Footer()
		</body>
	</html>
}
`
		},
		"app/web/layouts/app.templ": func() string {
			return `package layouts

var (
	title = "{{.AppName}}"
)

templ App() {
	@BaseLayout() {
		@components.navigation.Sidebar()
		<div class="max-w-7xl mx-auto">
			{ children... }
		</div>
	}
}
`
		},
		"app/web/components/header/header.templ": func() string {
			return `package components

templ Header(){
	<header>
    	<h1>Welcome to {{.AppName}}</h1>
    </header>
}
`
		},
		"app/web/components/footer/footer.templ": func() string {
			return `package components

templ Footer(){
	<footer>
   		<p>© {{.CurrentYear}} {{.AppName}}</p>
    </footer>
}
`
		},
		"app/web/pages/home.templ": func() string {
			return `package pages

templ Home(){
	<h2>Home Page</h2>
	<p>This is the home page.</p>
}
`
		},
		"app/web/pages/about.templ": func() string {
			return `package pages

templ About(){
	<h2>About Page</h2>
	<p>This is the about page.</p>
}
`
		},
		"app/web/components/navigation/sidebar.templ": func() string {
			return `package navigation

templ Sidebar(){
	<div>
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
			<li>Item 3</li>
		</ul>
	</div>
}
`
		},
		"app/web/views.go": func() string {
			return `package views

import (
	"fmt"
	"os"
	"path/filepath"
)

// Asset retrieves the content of a file from the current working directory under app/assets/{any_folder}/{any_asset}
func Asset(fileName string) ([]byte, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf(">>Gost>>error getting current working directory: %v", err)
	}

	// Construct the full path to the asset
	assetPath := filepath.Join(cwd, "app", "assets", fileName)

	// Read the file content
	content, err := os.ReadFile(assetPath)
	if err != nil {
		return nil, fmt.Errorf(">>Gost>> error reading file %s: %v", assetPath, err)
	}

	return content, nil
}
`
		},
		"app/web/errors/404.templ": func() string {
			return `package errors

templ _404(){
	<div>404 Page Not Found</div>
}
`
		},
		"app/web/errors/500.templ": func() string {
			return `package errors

templ _500(){
		<div>500 Internal Server Error</div>
}
`
		},
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
