package types

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenTypesPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenTypesPlugin) Init() error {
	// Initialize Files
	g.Files = map[string]func() string{
		"app/types/core/gost.go": func() string {
			return `package core

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/a-h/templ"
    "github.com/gorilla/sessions"
)

var store *sessions.CookieStore

type HandlerFunc func(*Gost) error
type ErrorHandlerFunc func(*Gost, error)

type AuthKey struct{}

type Auth interface {
    Check() bool
}

type DefaultAuth struct{}

func (DefaultAuth) Check() bool { return false }

type Gost struct {
    Response http.ResponseWriter
    Request  *http.Request
    Router   Router
}

var errorHandler ErrorHandlerFunc = func(g *Gost, err error) {
    err = g.Text(http.StatusInternalServerError, err.Error())
    if err != nil {
        log.Println(err)
    }
}

func SetErrorHandler(h ErrorHandlerFunc) {
    errorHandler = h
}

func (g *Gost) Auth() Auth {
    if auth, ok := g.Request.Context().Value(AuthKey{}).(Auth); ok {
        return auth
    }
    log.Println("Warning: Authentication not set")
    return DefaultAuth{}
}

func (g *Gost) GetSession(name string) (*sessions.Session, error) {
    return store.Get(g.Request, name)
}

func (g *Gost) Redirect(status int, url string) error {
    if g.Request.Header.Get("HX-Request") != "" {
        g.Response.Header().Set("HX-Redirect", url)
        g.Response.WriteHeader(http.StatusSeeOther)
        return nil
    }
    http.Redirect(g.Response, g.Request, url, status)
    return nil
}

func (g *Gost) FormValue(name string) string {
    return g.Request.PostFormValue(name)
}

func (g *Gost) JSON(status int, v interface{}) error {
    g.Response.Header().Set("Content-Type", "application/json")
    g.Response.WriteHeader(status)
    return json.NewEncoder(g.Response).Encode(v)
}

func (g *Gost) Text(status int, msg string) error {
    g.Response.Header().Set("Content-Type", "text/plain")
    g.Response.WriteHeader(status)
    _, err := g.Response.Write([]byte(msg))
    return err
}

func (g *Gost) Bytes(status int, b []byte) error {
    g.Response.Header().Set("Content-Type", "application/octet-stream")
    g.Response.WriteHeader(status)
    _, err := g.Response.Write(b)
    return err
}

func (g *Gost) Render(c templ.Component) error {
    return c.Render(g.Request.Context(), g.Response)
}

func (g *Gost) GetEnv(name, def string) string {
    if value := os.Getenv(name); value != "" {
        return value
    }
    return def
}

func GetEnv(name, def string) string {
    if value := os.Getenv(name); value != "" {
        return value
    }
    return def
}

func IsDevelopment() bool {
    return os.Getenv("GOST_ENV") == "development"
}

func IsProduction() bool {
    return os.Getenv("GOST_ENV") == "production"
}

func Env() string {
    return os.Getenv("GOST_ENV")
}

func init() {
    appSecret := os.Getenv("GOST_SECRET")
    if len(appSecret) < 32 {
        log.Fatalf("Invalid GOST_SECRET variable. Ensure it is set in your .env file.")
    }
    store = sessions.NewCookieStore([]byte(appSecret))
}

// Router interface to abstract the underlying server implementation
type Router interface {
    Use(middleware ...interface{})
    Get(path string, handler HandlerFunc)
    Post(path string, handler HandlerFunc)
    NotFound(handler HandlerFunc)
}

{{if eq .BackendPkg "chi"}}
type chiRouter struct {
    router *chi.Mux
}

func (c *chiRouter) Use(middleware ...interface{}) {
    for _, m := range middleware {
        c.router.Use(m.(func(http.Handler) http.Handler))
    }
}

func (c *chiRouter) Get(path string, handler HandlerFunc) {
    c.router.Get(path, func(w http.ResponseWriter, r *http.Request) {
        g := &Gost{Response: w, Request: r, Router: c}
        if err := handler(g); err != nil {
            errorHandler(g, err)
        }
    })
}

func (c *chiRouter) Post(path string, handler HandlerFunc) {
    c.router.Post(path, func(w http.ResponseWriter, r *http.Request) {
        g := &Gost{Response: w, Request: r, Router: c}
        if err := handler(g); err != nil {
            errorHandler(g, err)
        }
    })
}

func (c *chiRouter) NotFound(handler HandlerFunc) {
    c.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
        g := &Gost{Response: w, Request: r, Router: c}
        if err := handler(g); err != nil {
            errorHandler(g, err)
        }
    })
}
{{end}}

{{if eq .BackendPkg "echo"}}
type echoRouter struct {
    router *echo.Echo
}

func (e *echoRouter) Use(middleware ...interface{}) {
    for _, m := range middleware {
        e.router.Use(m.(echo.MiddlewareFunc))
    }
}

func (e *echoRouter) Get(path string, handler HandlerFunc) {
    e.router.GET(path, func(c echo.Context) error {
        g := &Gost{Response: c.Response(), Request: c.Request(), Router: e}
        if err := handler(g); err != nil {
            errorHandler(g, err)
        }
        return nil
    })
}

func (e *echoRouter) Post(path string, handler HandlerFunc) {
    e.router.POST(path, func(c echo.Context) error {
        g := &Gost{Response: c.Response(), Request: c.Request(), Router: e}
        if err := handler(g); err != nil {
            errorHandler(g, err)
        }
        return nil
    })
}

func (e *echoRouter) NotFound(handler HandlerFunc) {
    e.router.HTTPErrorHandler = func(err error, c echo.Context) {
        if c.Response().Committed {
            return
        }
        g := &Gost{Response: c.Response(), Request: c.Request(), Router: e}
        if err := handler(g); err != nil {
            errorHandler(g, err)
        }
    }
}
{{end}}

{{if eq .BackendPkg "gin"}}
type ginRouter struct {
    router *gin.Engine
}

func (g *ginRouter) Use(middleware ...interface{}) {
    for _, m := range middleware {
        g.router.Use(m.(gin.HandlerFunc))
    }
}

func (g *ginRouter) Get(path string, handler HandlerFunc) {
    g.router.GET(path, func(c *gin.Context) {
        gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
        if err := handler(gost); err != nil {
            errorHandler(gost, err)
        }
    })
}

func (g *ginRouter) Post(path string, handler HandlerFunc) {
    g.router.POST(path, func(c *gin.Context) {
        gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
        if err := handler(gost); err != nil {
            errorHandler(gost, err)
        }
    })
}

func (g *ginRouter) NotFound(handler HandlerFunc) {
    g.router.NoRoute(func(c *gin.Context) {
        gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
        if err := handler(gost); err != nil {
            errorHandler(gost, err)
        }
    })
}
{{end}}

func NewRouter(backend string) Router {
    switch backend {
    case "chi":
        return &chiRouter{router: chi.NewRouter()}
    case "echo":
        return &echoRouter{router: echo.New()}
    case "gin":
        return &ginRouter{router: gin.Default()}
    default:
        log.Fatalf("Unsupported backend: %s", backend)
        return nil
    }
}
`
		},
	}

	return nil
}

func (g *GenTypesPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenTypesPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenTypesPlugin) Name() string {
	return "GenTypesPlugin"
}

func (g *GenTypesPlugin) Version() string {
	return "1.0.0"
}

func (g *GenTypesPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenTypesPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenTypesPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenTypesPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenTypesPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/types"
}

func (g *GenTypesPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenTypesPlugin(data config.ProjectData) *GenTypesPlugin {
	return &GenTypesPlugin{
		Data: data,
	}
}
