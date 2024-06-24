package router

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
			"app/router/router.go": func() string {
				return `package router

import (
    "{{.AppName}}/app/handlers"
    "{{.AppName}}/app/middleware"
    "prelude"
    "log"
    {{if eq .BackendPkg "chi"}}
    "github.com/go-chi/chi/v5/middleware"
    {{end}}
    {{if eq .BackendPkg "echo"}}
    "github.com/labstack/echo/v5/middleware"
    {{end}}
)

func InitializeMiddleware(router prelude.Router) {
    router.Use(middleware.Logger)
    router.Use(middleware.Recover)
    router.Use(middleware.WithRequestURL)			
}

func InitializeRoutes(router prelude.Router) {
    {{if .IncludeAuth}}
    authConfig := prelude.AuthConfig{
        AuthFunc:    authenticateUser,
        RedirectURL: "/login",
    }

    router.Group(func(r prelude.Router) {
        r.Use(prelude.WithAuth(authConfig, false))

        r.Get("/", handlers.HomeHandler)
        r.Get("/about", handlers.AboutHandler)
        r.Get("/signin", handlers.SigninHandler)
        r.Get("/signup", handlers.SignupHandler)
    })

    router.Group(func(r prelude.Router) {
        r.Use(prelude.WithAuth(authConfig, true))

        // r.Get("/path", handlers.SomeProtectedHandler)
    })
    {{else}}
    router.Get("/", handlers.HomeHandler)
    router.Get("/about", handlers.AboutHandler)
    router.Get("/signin", handlers.SigninHandler)
    router.Get("/signup", handlers.SignupHandler)
    {{end}}

    router.NotFound(handlers.NotFoundHandler)
}

{{if .IncludeAuth}}
func authenticateUser(g *prelude.Gost) (prelude.Auth, error) {
    return prelude.DefaultAuth{}, nil
}
{{end}}

func InitRoutes(backend string) prelude.Router {
    router := prelude.NewRouter(backend)
    InitializeMiddleware(router)
    InitializeRoutes(router)
    return router
}
`
			},
		},
	}
}
