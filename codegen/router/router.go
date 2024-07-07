package router

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenRouterPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenRouterPlugin) Init() error {
	// Initialize Files
	g.Files = map[string]func() string{
		"app/router/router.go": func() string {
			return `package router

import (
    "{{.AppName}}/app/handlers"
    "{{.AppName}}/app/middleware"
    "{{.AppName}}/app/types/core"
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
	}

	return nil
}

func (g *GenRouterPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenRouterPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenRouterPlugin) Name() string {
	return "GenRouterPlugin"
}

func (g *GenRouterPlugin) Version() string {
	return "1.0.0"
}

func (g *GenRouterPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenRouterPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenRouterPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenRouterPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenRouterPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/router"
}

func (g *GenRouterPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenRouterPlugin(data config.ProjectData) *GenRouterPlugin {
	return &GenRouterPlugin{
		Data: data,
	}
}
