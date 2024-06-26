package resources

import (
	"strings"

	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type ResourcePlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (r *ResourcePlugin) Init() error {
	backendPkg := strings.ToLower(r.Data.BackendPkg)
	// Initialize Files based on the backend package
	r.Files = map[string]func() string{
		"app/models/{{ .ResourceName }}.go": func() string {
			return `package models

type {{ .ResourceName }} struct {
	ID   uint   ` + "`json:\"id\" gorm:\"primaryKey\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
		},
		"app/ui/frontend/router/{{ .ResourceName }}.js": func() string {
			return `import React from 'react';
import { Route, Switch } from 'react-router-dom';
import {{ .ResourceName }}Page from '../pages/{{ .ResourceName }}Page';

const {{ .ResourceName }}Router = () => (
  <Switch>
    <Route exact path="/{{ .ResourceName }}" component={{ .ResourceName }}Page} />
  </Switch>
);

export default {{ .ResourceName }}Router;
`
		},
		"backend/router/{{ .ResourceName }}.go": func() string {
			if backendPkg == "echo" {
				return `package router

import (
	"{{ .AppName }}/handlers"
	"github.com/labstack/echo/v4"
)

func Register{{ .ResourceName }}Routes(e *echo.Echo) {
	e.GET("/{{ .ResourceName }}s", handlers.Get{{ .ResourceName }}s)
	e.GET("/{{ .ResourceName }}s/:id", handlers.Get{{ .ResourceName }})
	e.POST("/{{ .ResourceName }}s", handlers.Create{{ .ResourceName }})
	e.PUT("/{{ .ResourceName }}s/:id", handlers.Update{{ .ResourceName }})
	e.DELETE("/{{ .ResourceName }}s/:id", handlers.Delete{{ .ResourceName }})
}
`
			}
			// Add more backend packages as needed
			return ""
		},
		"api/{{ .ResourceName }}.go": func() string {
			return `package api

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func Get{{ .ResourceName }}s(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Get all {{ .ResourceName }}s")
}

func Get{{ .ResourceName }}(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Get single {{ .ResourceName }}")
}

func Create{{ .ResourceName }}(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Create {{ .ResourceName }}")
}

func Update{{ .ResourceName }}(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Update {{ .ResourceName }}")
}

func Delete{{ .ResourceName }}(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Delete {{ .ResourceName }}")
}
`
		},
		"events/{{ .ResourceName }}.go": func() string {
			return `package events

import "fmt"

func On{{ .ResourceName }}Created() {
	fmt.Println("{{ .ResourceName }} created")
}

func On{{ .ResourceName }}Updated() {
	fmt.Println("{{ .ResourceName }} updated")
}

func On{{ .ResourceName }}Deleted() {
	fmt.Println("{{ .ResourceName }} deleted")
}
`
		},
		"handlers/{{ .ResourceName }}.go": func() string {
			return `package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func Get{{ .ResourceName }}s(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Get all {{ .ResourceName }}s")
}

func Get{{ .ResourceName }}(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Get single {{ .ResourceName }}")
}

func Create{{ .ResourceName }}(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Create {{ .ResourceName }}")
}

func Update{{ .ResourceName }}(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Update {{ .ResourceName }}")
}

func Delete{{ .ResourceName }}(c echo.Context) error {
	// Implement logic
	return c.JSON(http.StatusOK, "Delete {{ .ResourceName }}")
}
`
		},
		"error_handling/{{ .ResourceName }}_error.go": func() string {
			return `package errorhandling

import "github.com/labstack/echo/v4"

func Handle{{ .ResourceName }}Error(err error, c echo.Context) {
	// Implement error handling logic
}
`
		},
		"templ_views/{{ .ResourceName }}_list.html": func() string {
			return `<html>
<head>
    <title>{{ .ResourceName }} List</title>
</head>
<body>
    <h1>{{ .ResourceName }} List</h1>
    <!-- List view -->
</body>
</html>
`
		},
		"templ_views/{{ .ResourceName }}_view.html": func() string {
			return `<html>
<head>
    <title>{{ .ResourceName }} View</title>
</head>
<body>
    <h1>{{ .ResourceName }} View</h1>
    <!-- View details -->
</body>
</html>
`
		},
		"templ_views/{{ .ResourceName }}_create.html": func() string {
			return `<html>
<head>
    <title>Create {{ .ResourceName }}</title>
</head>
<body>
    <h1>Create {{ .ResourceName }}</h1>
    <!-- Create form -->
</body>
</html>
`
		},
		"templ_views/{{ .ResourceName }}_edit.html": func() string {
			return `<html>
<head>
    <title>Edit {{ .ResourceName }}</title>
</head>
<body>
    <h1>Edit {{ .ResourceName }}</h1>
    <!-- Edit form -->
</body>
</html>
`
		},
	}
	return nil
}

func (r *ResourcePlugin) Execute() error {
	return r.Generate(r.Data)
}

func (r *ResourcePlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (r *ResourcePlugin) Name() string {
	return "Resources Plugin"
}

func (r *ResourcePlugin) Version() string {
	return "1.0.0"
}

func (r *ResourcePlugin) Dependencies() []string {
	return []string{}
}

func (r *ResourcePlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (r *ResourcePlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (r *ResourcePlugin) Website() string {
	return "https://theHamdiz.me"
}

func (r *ResourcePlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/plugins/resources"
}

func (r *ResourcePlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, r.Files)
}

func NewResourcePlugin(data config.ProjectData) *ResourcePlugin {
	return &ResourcePlugin{
		Data: data,
	}
}
