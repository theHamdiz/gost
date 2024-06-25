package views

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
			"app/views/components/head.templ": func() string {
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
			"app/views/layouts/base.templ": func() string {
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
			"app/views/layouts/app.templ": func() string {
				return `package layouts

var (
	title = "gost project"
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
			"app/views/components/header/header.templ": func() string {
				return `package components

templ Header(){
	<header>
    	<h1>Welcome to {{.AppName}}</h1>
    </header>
}
`
			},
			"app/views/components/footer/footer.templ": func() string {
				return `package components

templ Footer(){
	<footer>
   		<p>© {{.CurrentYear}} {{.AppName}}</p>
    </footer>
}
`
			},
			"app/views/pages/home.templ": func() string {
				return `package pages

templ Home(){
	<h2>Home Page</h2>
	<p>This is the home page.</p>
}
`
			},
			"app/views/pages/about.templ": func() string {
				return `package pages

templ About(){
	<h2>About Page</h2>
	<p>This is the about page.</p>
}
`
			},
			"app/views/components/navigation/sidebar.templ": func() string {
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
			"app/views/views.go": func() string {
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
			"app/views/errors/404.templ": func() string {
				return `package errors

templ _404(){
	<div>404 Page Not Found</div>
}
`
			},
			"app/views/errors/500.templ": func() string {
				return `package errors

templ _500(){
		<div>500 Internal Server Error</div>
}
`
			},
		},
	}
}
