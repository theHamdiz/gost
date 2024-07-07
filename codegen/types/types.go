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
		"app/types/gost/gost.go": func() string {
			return `package core

import (
	"{{.AppName}}/app/types/core"
	"encoding/json"
	"reflect"
    {{- if eq .BackendPkg "echo"}}
    "github.com/labstack/echo/v5"
    {{- else if eq .BackendPkg "gin"}}
    "github.com/gin-gonic/gin"
    {{- else if eq .BackendPkg "chi"}}
    "github.com/go-chi/chi/v5"
    {{- else if eq .BackendPkg "stdlib"}}
    "net/http"
    {{- end }}

    "github.com/a-h/templ"
    "github.com/gorilla/sessions"
)

type Configurable interface {
	{{- if eq .PreferredConfigFormat "env"}}
	SaveAsEnv(filePath string) error
	{{- else if eq .PreferredConfigFormat "json"}}
	SaveAsJSON(filePath string) error
	{{- else if eq .PreferredConfigFormat "toml"}}
	SaveAsTOML(filePath string) error
	{{- else }}
	SaveAsYAML(filePath string) error
	{{- end }}
}

// Route -> Define the Route and ResourceRoutes structs
type Route struct {
	Handler     interface{}
	Method      string
	Middlewares []interface{}
	Path        string
}

// RouteGroup ->    A collection of routes prefixed with certain prefix.
type RouteGroup struct {
	Middlewares []interface{}
	Prefix      string
	Routes      []Route
}

// ResourceRoutes -> Define the Routes associated with a resource.
type ResourceRoutes struct {
	Controller  interface{}
	Middlewares []interface{}
	Resource    string
}

func New() *Gost {
    return &Gost{
        {{- if eq .BackendPkg "echo"}}
        router: echo.New(),
        {{- else if eq .BackendPkg "gin"}}
        router: gin.Default(),
        {{- else if eq .BackendPkg "chi"}}
        router: chi.NewRouter(),
        {{- else if eq .BackendPkg "stdlib"}}
        router: http.NewServeMux(),
        {{- end }}
    }
}

{{- if eq .BackendPkg "echo"}}
func (g *Gost) AddResource(resource string, controller interface{}, middlewares ...echo.MiddlewareFunc) {
    g.resourceRoutes = append(g.resourceRoutes, ResourceRoutes{
        Resource:    resource,
        Controller:  controller,
        Middlewares: middlewares,
    })
    
    g.registerResourceRoutes(resource, controller, middlewares...)
}

func (g *Gost) registerResourceRoutes(resource string, controller interface{}, middlewares ...echo.MiddlewareFunc) {
    basePath := "/" + resource
    g.router.GET(basePath, wrapHandler(controller, "Index", middlewares...))
    g.router.GET(basePath+"/:id", wrapHandler(controller, "Show", middlewares...))
    g.router.POST(basePath, wrapHandler(controller, "Create", middlewares...))
    g.router.PUT(basePath+"/:id", wrapHandler(controller, "Update", middlewares...))
    g.router.DELETE(basePath+"/:id", wrapHandler(controller, "Delete", middlewares...))
}

func wrapHandler(controller interface{}, methodName string, middlewares ...echo.MiddlewareFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        method := reflect.ValueOf(controller).MethodByName(methodName)
        if !method.IsValid() {
            return echo.NewHTTPError(http.StatusNotFound, "method not found")
        }
        args := []reflect.Value{reflect.ValueOf(c)}
        results := method.Call(args)
        if len(results) == 2 && results[1].Interface() != nil {
            return results[1].Interface().(error)
        }
        return results[0].Interface().(error)
    }
}
{{- else if eq .BackendPkg "gin"}}
func (g *Gost) AddResource(resource string, controller interface{}, middlewares ...gin.HandlerFunc) {
    g.resourceRoutes = append(g.resourceRoutes, ResourceRoutes{
        Resource:    resource,
        Controller:  controller,
        Middlewares: middlewares,
    })

    g.registerResourceRoutes(resource, controller, middlewares...)
}

func (g *Gost) registerResourceRoutes(resource string, controller interface{}, middlewares ...gin.HandlerFunc) {
    basePath := "/" + resource
    g.router.GET(basePath, wrapHandler(controller, "Index", middlewares...))
    g.router.GET(basePath+"/:id", wrapHandler(controller, "Show", middlewares...))
    g.router.POST(basePath, wrapHandler(controller, "Create", middlewares...))
    g.router.PUT(basePath+"/:id", wrapHandler(controller, "Update", middlewares...))
    g.router.DELETE(basePath+"/:id", wrapHandler(controller, "Delete", middlewares...))
}

func wrapHandler(controller interface{}, methodName string, middlewares ...gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        method := reflect.ValueOf(controller).MethodByName(methodName)
        if !method.IsValid() {
            c.JSON(http.StatusNotFound, gin.H{"error": "method not found"})
            return
        }
        args := []reflect.Value{reflect.ValueOf(c)}
        results := method.Call(args)
        if len(results) == 2 && results[1].Interface() != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": results[1].Interface().(error).Error()})
            return
        }
        c.JSON(http.StatusOK, results[0].Interface())
    }
}
{{- else if eq .BackendPkg "chi"}}
func (g *Gost) AddResource(resource string, controller interface{}, middlewares ...func(http.Handler) http.Handler) {
    g.resourceRoutes = append(g.resourceRoutes, ResourceRoutes{
        Resource:    resource,
        Controller:  controller,
        Middlewares: middlewares,
    })

    g.registerResourceRoutes(resource, controller, middlewares...)
}

func (g *Gost) registerResourceRoutes(resource string, controller interface{}, middlewares ...func(http.Handler) http.Handler) {
    basePath := "/" + resource
    g.router.Method("GET", basePath, wrapHandler(controller, "Index", middlewares...))
    g.router.Method("GET", basePath+"/:id", wrapHandler(controller, "Show", middlewares...))
    g.router.Method("POST", basePath, wrapHandler(controller, "Create", middlewares...))
    g.router.Method("PUT", basePath+"/:id", wrapHandler(controller, "Update", middlewares...))
    g.router.Method("DELETE", basePath+"/:id", wrapHandler(controller, "Delete", middlewares...))
}

func wrapHandler(controller interface{}, methodName string, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        method := reflect.ValueOf(controller).MethodByName(methodName)
        if !method.IsValid() {
            http.Error(w, "method not found", http.StatusNotFound)
            return
        }
        args := []reflect.Value{reflect.ValueOf(r)}
        results := method.Call(args)
        if len(results) == 2 && results[1].Interface() != nil {
            http.Error(w, results[1].Interface().(error).Error(), http.StatusInternalServerError)
            return
        }
        // Write the response
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(results[0].Interface().(string)))
    }
}
{{- else if eq .BackendPkg "stdlib"}}
func (g *Gost) AddResource(resource string, controller interface{}, middlewares ...func(http.Handler) http.Handler) {
    g.resourceRoutes = append(g.resourceRoutes, ResourceRoutes{
        Resource:    resource,
        Controller:  controller,
        Middlewares: middlewares,
    })

    g.registerResourceRoutes(resource, controller, middlewares...)
}

func (g *Gost) registerResourceRoutes(resource string, controller interface{}, middlewares ...func(http.Handler) http.Handler) {
    basePath := "/" + resource
    g.router.Handle(basePath, wrapHandler(controller, "Index", middlewares...))
    g.router.Handle(basePath+"/:id", wrapHandler(controller, "Show", middlewares...))
    g.router.Handle(basePath, wrapHandler(controller, "Create", middlewares...))
    g.router.Handle(basePath+"/:id", wrapHandler(controller, "Update", middlewares...))
    g.router.Handle(basePath+"/:id", wrapHandler(controller, "Delete", middlewares...))
}

func wrapHandler(controller interface{}, methodName string, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        method := reflect.ValueOf(controller).MethodByName(methodName)
        if !method.IsValid() {
            http.Error(w, "method not found", http.StatusNotFound)
            return
        }
        args := []reflect.Value{reflect.ValueOf(r)}
        results := method.Call(args)
        if len(results) == 2 && results[1].Interface() != nil {
            http.Error(w, results[1].Interface().(error).Error(), http.StatusInternalServerError)
            return
        }
        // Write the response
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(results[0].Interface().(string)))
    }
}
{{- end }}


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
	routeGroups    []RouteGroup
	resourceRoutes []ResourceRoutes
	eventHooks     []func(evt *core.ServeEvent) error
	request        http.Request
	response       http.ResponseWriter
}

func (g *Gost) Auth() Auth {
    if auth, ok := g.request.Context().Value(AuthKey{}).(Auth); ok {
        return auth
    }
    log.Println("Warning: Authentication not set")
    return DefaultAuth{}
}

func (g *Gost) GetSession(name string) (*sessions.Session, error) {
    return store.Get(g.request, name)
}

func (g *Gost) Redirect(status int, url string) error {
    if g.request.Header.Get("HX-Request") != "" {
        g.response.Header().Set("HX-Redirect", url)
        g.response.WriteHeader(http.StatusSeeOther)
        return nil
    }
    http.Redirect(g.response, &g.request, url, status)
    return nil
}

func (g *Gost) FormValue(name string) string {
	return g.request.PostFormValue(name)
}

func (g *Gost) JSON(status int, v interface{}) error {
    g.response.Header().Set("Content-Type", "application/json")
    g.response.WriteHeader(status)
    return json.NewEncoder(g.response).Encode(v)
}

func (g *Gost) Text(status int, msg string) error {
    g.response.Header().Set("Content-Type", "text/plain")
    g.response.WriteHeader(status)
    _, err := g.response.Write([]byte(msg))
    return err
}

func (g *Gost) Bytes(status int, b []byte) error {
    g.response.Header().Set("Content-Type", "application/octet-stream")
    g.response.WriteHeader(status)
    _, err := g.response.Write(b)
    return err
}

func (g *Gost) Render(c templ.Component) error {
    return c.Render(g.request.Context(), g.response)
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
	Put(path string, handler HandlerFunc)
	Delete(path string, handler HandlerFunc)
	Patch(path string, handler HandlerFunc)
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

func (g *ginRouter) Put(path string, handler HandlerFunc) {
	g.router.Put(path, func(c *gin.Context) {
		gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
		if err := handler(gost); err != nil {
			errorHandler(gost, err)
		}
	})
}

func (g *ginRouter) Patch(path string, handler HandlerFunc) {
	g.router.Patch(path, func(c *gin.Context) {
		gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
		if err := handler(gost); err != nil {
			errorHandler(gost, err)
		}
	})
}

func (g *ginRouter) Delete(path string, handler HandlerFunc) {
	g.router.Delete(path, func(c *gin.Context) {
		gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
		if err := handler(gost); err != nil {
			errorHandler(gost, err)
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


func NewRouter() Router {
	return &chiRouter{router: chi.NewRouter()}
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

func (g *ginRouter) Put(path string, handler HandlerFunc) {
	g.router.PUT(path, func(c *gin.Context) {
		gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
		if err := handler(gost); err != nil {
			errorHandler(gost, err)
		}
	})
}

func (g *ginRouter) Patch(path string, handler HandlerFunc) {
	g.router.PATCH(path, func(c *gin.Context) {
		gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
		if err := handler(gost); err != nil {
			errorHandler(gost, err)
		}
	})
}

func (g *ginRouter) Delete(path string, handler HandlerFunc) {
	g.router.DELETE(path, func(c *gin.Context) {
		gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
		if err := handler(gost); err != nil {
			errorHandler(gost, err)
		}
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

func NewRouter() Router {
	return return &echoRouter{router: echo.New()}
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

func (g *ginRouter) Put(path string, handler HandlerFunc) {
	g.router.PUT(path, func(c *gin.Context) {
		gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
		if err := handler(gost); err != nil {
			errorHandler(gost, err)
		}
	})
}

func (g *ginRouter) Patch(path string, handler HandlerFunc) {
	g.router.PATCH(path, func(c *gin.Context) {
		gost := &Gost{Response: c.Writer, Request: c.Request, Router: g}
		if err := handler(gost); err != nil {
			errorHandler(gost, err)
		}
	})
}

func (g *ginRouter) Delete(path string, handler HandlerFunc) {
	g.router.DELETE(path, func(c *gin.Context) {
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

func NewRouter() Router {
	return &ginRouter{router: gin.Default()}
}

{{end}}
`
		},
		"app/types/mailer/mailer.go": func() string {
			return `package mailer

            import (
                "context"
                {{- if eq .BackendPkg "echo"}}
                "github.com/labstack/echo/v5"
                {{- else if eq .BackendPkg "gin"}}
                "github.com/gin-gonic/gin"
                {{- else if eq .BackendPkg "chi"}}
                "github.com/go-chi/chi/v5"
                {{- else if eq .BackendPkg "stdlib"}}
                "net/http"
                {{- end }}
            )
            
            // Mailer interface defines the methods for sending emails
            type Mailer interface {
                {{- if eq .BackendPkg "echo"}}
                SendAdminResetPasswordMail(ctx context.Context, email string, data interface{}) error
                SendRecordResetPasswordMail(ctx context.Context, email string, data interface{}) error
                SendRecordVerificationMail(ctx context.Context, email string, data interface{}) error
                SendRecordChangeEmailMail(ctx context.Context, email string, data interface{}) error
                {{- else if eq .BackendPkg "gin"}}
                SendAdminResetPasswordMail(ctx *gin.Context, email string, data interface{}) error
                SendRecordResetPasswordMail(ctx *gin.Context, email string, data interface{}) error
                SendRecordVerificationMail(ctx *gin.Context, email string, data interface{}) error
                SendRecordChangeEmailMail(ctx *gin.Context, email string, data interface{}) error
                {{- else if eq .BackendPkg "chi"}}
                SendAdminResetPasswordMail(r chi.Router, email string, data interface{}) error
                SendRecordResetPasswordMail(r chi.Router, email string, data interface{}) error
                SendRecordVerificationMail(r chi.Router, email string, data interface{}) error
                SendRecordChangeEmailMail(r chi.Router, email string, data interface{}) error
                {{- else if eq .BackendPkg "stdlib"}}
                SendAdminResetPasswordMail(r *http.Request, email string, data interface{}) error
                SendRecordResetPasswordMail(r *http.Request, email string, data interface{}) error
                SendRecordVerificationMail(r *http.Request, email string, data interface{}) error
                SendRecordChangeEmailMail(r *http.Request, email string, data interface{}) error
                {{- end }}
            }
            
            // MailClient is a struct that implements the Mailer interface
            type MailClient struct {
                // Add any necessary fields here, such as configuration or dependencies
            }
            
            // NewMailClient creates a new instance of MailClient
            func NewMailClient() *MailClient {
                return &MailClient{}
            }
            
            // Implement Mailer interface methods based on backend package choice
            {{- if eq .BackendPkg "echo"}}
            func (mc *MailClient) SendAdminResetPasswordMail(ctx context.Context, email string, data interface{}) error {
                // Implementation for sending admin reset password email using Echo context
                return nil
            }
            
            func (mc *MailClient) SendRecordResetPasswordMail(ctx context.Context, email string, data interface{}) error {
                // Implementation for sending record reset password email using Echo context
                return nil
            }
            
            func (mc *MailClient) SendRecordVerificationMail(ctx context.Context, email string, data interface{}) error {
                // Implementation for sending record verification email using Echo context
                return nil
            }
            
            func (mc *MailClient) SendRecordChangeEmailMail(ctx context.Context, email string, data interface{}) error {
                // Implementation for sending record change email mail using Echo context
                return nil
            }
            {{- else if eq .BackendPkg "gin"}}
            func (mc *MailClient) SendAdminResetPasswordMail(ctx *gin.Context, email string, data interface{}) error {
                // Implementation for sending admin reset password email using Gin context
                return nil
            }
            
            func (mc *MailClient) SendRecordResetPasswordMail(ctx *gin.Context, email string, data interface{}) error {
                // Implementation for sending record reset password email using Gin context
                return nil
            }
            
            func (mc *MailClient) SendRecordVerificationMail(ctx *gin.Context, email string, data interface{}) error {
                // Implementation for sending record verification email using Gin context
                return nil
            }
            
            func (mc *MailClient) SendRecordChangeEmailMail(ctx *gin.Context, email string, data interface{}) error {
                // Implementation for sending record change email mail using Gin context
                return nil
            }
            {{- else if eq .BackendPkg "chi"}}
            func (mc *MailClient) SendAdminResetPasswordMail(r chi.Router, email string, data interface{}) error {
                // Implementation for sending admin reset password email using Chi router
                return nil
            }
            
            func (mc *MailClient) SendRecordResetPasswordMail(r chi.Router, email string, data interface{}) error {
                // Implementation for sending record reset password email using Chi router
                return nil
            }
            
            func (mc *MailClient) SendRecordVerificationMail(r chi.Router, email string, data interface{}) error {
                // Implementation for sending record verification email using Chi router
                return nil
            }
            
            func (mc *MailClient) SendRecordChangeEmailMail(r chi.Router, email string, data interface{}) error {
                // Implementation for sending record change email mail using Chi router
                return nil
            }
            {{- else if eq .BackendPkg "stdlib"}}
            func (mc *MailClient) SendAdminResetPasswordMail(r *http.Request, email string, data interface{}) error {
                // Implementation for sending admin reset password email using stdlib http request
                return nil
            }
            
            func (mc *MailClient) SendRecordResetPasswordMail(r *http.Request, email string, data interface{}) error {
                // Implementation for sending record reset password email using stdlib http request
                return nil
            }
            
            func (mc *MailClient) SendRecordVerificationMail(r *http.Request, email string, data interface{}) error {
                // Implementation for sending record verification email using stdlib http request
                return nil
            }
            
            func (mc *MailClient) SendRecordChangeEmailMail(r *http.Request, email string, data interface{}) error {
                // Implementation for sending record change email mail using stdlib http request
                return nil
            }
            {{- end }}
            `
		},
		"app/types/dao/dao.go": func() string {
			return `package dao`
		},
		"app/types/sessions/sessions.go": func() string {
			return `package sessions
import (
	"encoding/base64"
    "encoding/gob"
	"github.com/gorilla/securecookie"
    "github.com/gorilla/sessions"
	"net/http"
	"time"
)

type Session struct {
	ID         string
	Values     map[interface{}]interface{}
	Options    *Options
	IsNew      bool
}

type Options struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

type CookieStore struct {
    Codecs  []securecookie.Codec
    Options *sessions.Options
}

func(store *CookieStore) Get(r *http.Request, name string) (*sessions.Session, error) {
    return sessions.Get(r, name)
}

func NewCookieStore(keyPairs ...[]byte) *CookieStore {
    return &CookieStore{
        Codecs: securecookie.CodecsFromPairs(keyPairs...),
        Options: &sessions.Options{
            Path:   "/",
            MaxAge: 86400 * 30,
        },
    }
}

func init() {
    gob.Register(map[interface{}]interface{}{})
}
`
		},
		"app/types/models/models.go": func() string {
			return `package models`
		},
		"app/types/events/event.go": func() string {
			return `package event

import (
    "context"
    "reflect"
    {{- if eq .BackendPkg "echo"}}
    "github.com/labstack/echo/v5"
    {{- else if eq .BackendPkg "gin"}}
    "github.com/gin-gonic/gin"
    {{- else if eq .BackendPkg "chi"}}
    "github.com/go-chi/chi/v5"
    {{- else if eq .BackendPkg "stdlib"}}
    "net/http"
    {{- end }}
)

// Define Event Types
type EventType string

const (
    BeforeBootstrap    EventType = "BeforeBootstrap"
    AfterBootstrap     EventType = "AfterBootstrap"
    BeforeServe        EventType = "BeforeServe"
    BeforeApiError     EventType = "BeforeApiError"
    AfterApiError      EventType = "AfterApiError"
    Terminate          EventType = "Terminate"
    OnModelBeforeCreate EventType = "OnModelBeforeCreate"
    OnModelAfterCreate  EventType = "OnModelAfterCreate"
    OnModelBeforeUpdate EventType = "OnModelBeforeUpdate"
    OnModelAfterUpdate  EventType = "OnModelAfterUpdate"
    OnModelBeforeDelete EventType = "OnModelBeforeDelete"
    OnModelAfterDelete  EventType = "OnModelAfterDelete"
    OnMailerBeforeAdminResetPasswordSend EventType = "OnMailerBeforeAdminResetPasswordSend"
    OnMailerAfterAdminResetPasswordSend  EventType = "OnMailerAfterAdminResetPasswordSend"
    OnMailerBeforeRecordResetPasswordSend EventType = "OnMailerBeforeRecordResetPasswordSend"
    OnMailerAfterRecordResetPasswordSend  EventType = "OnMailerAfterRecordResetPasswordSend"
    OnMailerBeforeRecordVerificationSend EventType = "OnMailerBeforeRecordVerificationSend"
    OnMailerAfterRecordVerificationSend  EventType = "OnMailerAfterRecordVerificationSend"
    OnMailerBeforeRecordChangeEmailSend EventType = "OnMailerBeforeRecordChangeEmailSend"
    OnMailerAfterRecordChangeEmailSend  EventType = "OnMailerAfterRecordChangeEmailSend"
    OnRealtimeConnectRequest EventType = "OnRealtimeConnectRequest"
    OnRealtimeDisconnectRequest EventType = "OnRealtimeDisconnectRequest"
    OnRealtimeBeforeMessageSend EventType = "OnRealtimeBeforeMessageSend"
    OnRealtimeAfterMessageSend EventType = "OnRealtimeAfterMessageSend"
    OnRealtimeBeforeSubscribeRequest EventType = "OnRealtimeBeforeSubscribeRequest"
    OnRealtimeAfterSubscribeRequest EventType = "OnRealtimeAfterSubscribeRequest"
    OnSettingsListRequest EventType = "OnSettingsListRequest"
    OnSettingsBeforeUpdateRequest EventType = "OnSettingsBeforeUpdateRequest"
    OnSettingsAfterUpdateRequest EventType = "OnSettingsAfterUpdateRequest"
    OnFileDownloadRequest EventType = "OnFileDownloadRequest"
    OnFileBeforeTokenRequest EventType = "OnFileBeforeTokenRequest"
    OnFileAfterTokenRequest EventType = "OnFileAfterTokenRequest"
    OnAdminsListRequest EventType = "OnAdminsListRequest"
    OnAdminViewRequest EventType = "OnAdminViewRequest"
    OnAdminBeforeCreateRequest EventType = "OnAdminBeforeCreateRequest"
    OnAdminAfterCreateRequest EventType = "OnAdminAfterCreateRequest"
    OnAdminBeforeUpdateRequest EventType = "OnAdminBeforeUpdateRequest"
    OnAdminAfterUpdateRequest EventType = "OnAdminAfterUpdateRequest"
    OnAdminBeforeDeleteRequest EventType = "OnAdminBeforeDeleteRequest"
    OnAdminAfterDeleteRequest EventType = "OnAdminAfterDeleteRequest"
    OnAdminAuthRequest EventType = "OnAdminAuthRequest"
    OnAdminBeforeAuthWithPasswordRequest EventType = "OnAdminBeforeAuthWithPasswordRequest"
    OnAdminAfterAuthWithPasswordRequest EventType = "OnAdminAfterAuthWithPasswordRequest"
    OnAdminBeforeAuthRefreshRequest EventType = "OnAdminBeforeAuthRefreshRequest"
    OnAdminAfterAuthRefreshRequest EventType = "OnAdminAfterAuthRefreshRequest"
    OnAdminBeforeRequestPasswordResetRequest EventType = "OnAdminBeforeRequestPasswordResetRequest"
    OnAdminAfterRequestPasswordResetRequest EventType = "OnAdminAfterRequestPasswordResetRequest"
    OnAdminBeforeConfirmPasswordResetRequest EventType = "OnAdminBeforeConfirmPasswordResetRequest"
    OnAdminAfterConfirmPasswordResetRequest EventType = "OnAdminAfterConfirmPasswordResetRequest"
    OnRecordAuthRequest EventType = "OnRecordAuthRequest"
    OnRecordBeforeAuthWithPasswordRequest EventType = "OnRecordBeforeAuthWithPasswordRequest"
    OnRecordAfterAuthWithPasswordRequest EventType = "OnRecordAfterAuthWithPasswordRequest"
    OnRecordBeforeAuthWithOAuth2Request EventType = "OnRecordBeforeAuthWithOAuth2Request"
    OnRecordAfterAuthWithOAuth2Request EventType = "OnRecordAfterAuthWithOAuth2Request"
    OnRecordBeforeAuthRefreshRequest EventType = "OnRecordBeforeAuthRefreshRequest"
    OnRecordAfterAuthRefreshRequest EventType = "OnRecordAfterAuthRefreshRequest"
    OnRecordListExternalAuthsRequest EventType = "OnRecordListExternalAuthsRequest"
    OnRecordBeforeUnlinkExternalAuthRequest EventType = "OnRecordBeforeUnlinkExternalAuthRequest"
    OnRecordAfterUnlinkExternalAuthRequest EventType = "OnRecordAfterUnlinkExternalAuthRequest"
    OnRecordBeforeRequestPasswordResetRequest EventType = "OnRecordBeforeRequestPasswordResetRequest"
    OnRecordAfterRequestPasswordResetRequest EventType = "OnRecordAfterRequestPasswordResetRequest"
    OnRecordBeforeConfirmPasswordResetRequest EventType = "OnRecordBeforeConfirmPasswordResetRequest"
    OnRecordAfterConfirmPasswordResetRequest EventType = "OnRecordAfterConfirmPasswordResetRequest"
    OnRecordBeforeRequestVerificationRequest EventType = "OnRecordBeforeRequestVerificationRequest"
    OnRecordAfterRequestVerificationRequest EventType = "OnRecordAfterRequestVerificationRequest"
    OnRecordBeforeConfirmVerificationRequest EventType = "OnRecordBeforeConfirmVerificationRequest"
    OnRecordAfterConfirmVerificationRequest EventType = "OnRecordAfterConfirmVerificationRequest"
    OnRecordBeforeRequestEmailChangeRequest EventType = "OnRecordBeforeRequestEmailChangeRequest"
    OnRecordAfterRequestEmailChangeRequest EventType = "OnRecordAfterRequestEmailChangeRequest"
    OnRecordBeforeConfirmEmailChangeRequest EventType = "OnRecordBeforeConfirmEmailChangeRequest"
    OnRecordAfterConfirmEmailChangeRequest EventType = "OnRecordAfterConfirmEmailChangeRequest"
    OnRecordsListRequest EventType = "OnRecordsListRequest"
    OnRecordViewRequest EventType = "OnRecordViewRequest"
    OnRecordBeforeCreateRequest EventType = "OnRecordBeforeCreateRequest"
    OnRecordAfterCreateRequest EventType = "OnRecordAfterCreateRequest"
    OnRecordBeforeUpdateRequest EventType = "OnRecordBeforeUpdateRequest"
    OnRecordAfterUpdateRequest EventType = "OnRecordAfterUpdateRequest"
    OnRecordBeforeDeleteRequest EventType = "OnRecordBeforeDeleteRequest"
    OnRecordAfterDeleteRequest EventType = "OnRecordAfterDeleteRequest"
)

type EventHandler func(evt interface{}) error

type EventRegistry struct {
    handlers map[EventType][]EventHandler
}

func NewEventRegistry() *EventRegistry {
    return &EventRegistry{
        handlers: make(map[EventType][]EventHandler),
    }
}

func (r *EventRegistry) Register(eventType EventType, handler EventHandler) {
    r.handlers[eventType] = append(r.handlers[eventType], handler)
}

func (r *EventRegistry) Invoke(eventType EventType, evt interface{}) error {
    if handlers, found := r.handlers[eventType]; found {
        for _, handler := range handlers {
            if err := handler(evt); err != nil {
                return err
            }
        }
    }
    return nil
}

type Event interface {
	{{- if eq .BackendPkg "echo"}}
	OnBeforeBootstrap(handler func(evt *BootstrapEvent) error) Event
	OnAfterBootstrap(handler func(evt *BootstrapEvent) error) Event
	OnBeforeServe(handler func(evt *ServeEvent) error) Event
	OnBeforeApiError(handler func(evt *ApiErrorEvent) error) Event
	OnAfterApiError(handler func(evt *ApiErrorEvent) error) Event
	OnTerminate(handler func(evt *TerminateEvent) error) Event
	OnModelBeforeCreate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterCreate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelBeforeUpdate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterUpdate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelBeforeDelete(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterDelete(tags []string, handler func(evt *ModelEvent) error) Event
	OnMailerBeforeAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) Event
	OnMailerAfterAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) Event
	OnMailerBeforeRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerBeforeRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerBeforeRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnRealtimeConnectRequest(handler func(evt *RealtimeConnectEvent) error) Event
	OnRealtimeDisconnectRequest(handler func(evt *RealtimeDisconnectEvent) error) Event
	OnRealtimeBeforeMessageSend(handler func(evt *RealtimeMessageEvent) error) Event
	OnRealtimeAfterMessageSend(handler func(evt *RealtimeMessageEvent) error) Event
	OnRealtimeBeforeSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) Event
	OnRealtimeAfterSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) Event
	OnSettingsListRequest(handler func(evt *SettingsListEvent) error) Event
	OnSettingsBeforeUpdateRequest(handler func(evt *SettingsUpdateEvent) error) Event
	OnSettingsAfterUpdateRequest(handler func(evt *SettingsUpdateEvent) error) Event
	OnFileDownloadRequest(tags []string, handler func(evt *FileDownloadEvent) error) Event
	OnFileBeforeTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) Event
	OnFileAfterTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) Event
	OnAdminsListRequest(handler func(evt *AdminsListEvent) error) Event
	OnAdminViewRequest(handler func(evt *AdminViewEvent) error) Event
	OnAdminBeforeCreateRequest(handler func(evt *AdminCreateEvent) error) Event
	OnAdminAfterCreateRequest(handler func(evt *AdminCreateEvent) error) Event
	OnAdminBeforeUpdateRequest(handler func(evt *AdminUpdateEvent) error) Event
	OnAdminAfterUpdateRequest(handler func(evt *AdminUpdateEvent) error) Event
	OnAdminBeforeDeleteRequest(handler func(evt *AdminDeleteEvent) error) Event
	OnAdminAfterDeleteRequest(handler func(evt *AdminDeleteEvent) error) Event
	OnAdminAuthRequest(handler func(evt *AdminAuthEvent) error) Event
	OnAdminBeforeAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) Event
	OnAdminAfterAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) Event
	OnAdminBeforeAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) Event
	OnAdminAfterAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) Event
	OnAdminBeforeRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) Event
	OnAdminAfterRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) Event
	OnAdminBeforeConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) Event
	OnAdminAfterConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) Event
	OnRecordAuthRequest(tags []string, handler func(evt *RecordAuthEvent) error) Event
	OnRecordBeforeAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) Event
	OnRecordAfterAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) Event
	OnRecordBeforeAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) Event
	OnRecordAfterAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) Event
	OnRecordBeforeAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) Event
	OnRecordAfterAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) Event
	OnRecordListExternalAuthsRequest(tags []string, handler func(evt *RecordListExternalAuthsEvent) error) Event
	OnRecordBeforeUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) Event
	OnRecordAfterUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) Event
	OnRecordBeforeRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) Event
	OnRecordAfterRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) Event
	OnRecordBeforeConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) Event
	OnRecordAfterConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) Event
	OnRecordBeforeRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) Event
	OnRecordAfterRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) Event
	OnRecordBeforeConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) Event
	OnRecordAfterConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) Event
	OnRecordBeforeRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) Event
	OnRecordAfterRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) Event
	OnRecordBeforeConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) Event
	OnRecordAfterConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) Event
	OnRecordsListRequest(tags []string, handler func(evt *RecordsListEvent) error) Event
	OnRecordViewRequest(tags []string, handler func(evt *RecordViewEvent) error) Event
	OnRecordBeforeCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) Event
	OnRecordAfterCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) Event
	OnRecordBeforeUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) Event
	OnRecordAfterUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) Event
	OnRecordBeforeDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) Event
	OnRecordAfterDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) Event
	{{- else if eq .BackendPkg "gin"}}
	OnBeforeBootstrap(handler func(evt *BootstrapEvent) error) Event
	OnAfterBootstrap(handler func(evt *BootstrapEvent) error) Event
	OnBeforeServe(handler func(evt *ServeEvent) error) Event
	OnBeforeApiError(handler func(evt *ApiErrorEvent) error) Event
	OnAfterApiError(handler func(evt *ApiErrorEvent) error) Event
	OnTerminate(handler func(evt *TerminateEvent) error) Event
	OnModelBeforeCreate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterCreate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelBeforeUpdate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterUpdate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelBeforeDelete(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterDelete(tags []string, handler func(evt *ModelEvent) error) Event
	OnMailerBeforeAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) Event
	OnMailerAfterAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) Event
	OnMailerBeforeRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerBeforeRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerBeforeRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnRealtimeConnectRequest(handler func(evt *RealtimeConnectEvent) error) Event
	OnRealtimeDisconnectRequest(handler func(evt *RealtimeDisconnectEvent) error) Event
	OnRealtimeBeforeMessageSend(handler func(evt *RealtimeMessageEvent) error) Event
	OnRealtimeAfterMessageSend(handler func(evt *RealtimeMessageEvent) error) Event
	OnRealtimeBeforeSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) Event
	OnRealtimeAfterSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) Event
	OnSettingsListRequest(handler func(evt *SettingsListEvent) error) Event
	OnSettingsBeforeUpdateRequest(handler func(evt *SettingsUpdateEvent) error) Event
	OnSettingsAfterUpdateRequest(handler func(evt *SettingsUpdateEvent) error) Event
	OnFileDownloadRequest(tags []string, handler func(evt *FileDownloadEvent) error) Event
	OnFileBeforeTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) Event
	OnFileAfterTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) Event
	OnAdminsListRequest(handler func(evt *AdminsListEvent) error) Event
	OnAdminViewRequest(handler func(evt *AdminViewEvent) error) Event
	OnAdminBeforeCreateRequest(handler func(evt *AdminCreateEvent) error) Event
	OnAdminAfterCreateRequest(handler func(evt *AdminCreateEvent) error) Event
	OnAdminBeforeUpdateRequest(handler func(evt *AdminUpdateEvent) error) Event
	OnAdminAfterUpdateRequest(handler func(evt *AdminUpdateEvent) error) Event
	OnAdminBeforeDeleteRequest(handler func(evt *AdminDeleteEvent) error) Event
	OnAdminAfterDeleteRequest(handler func(evt *AdminDeleteEvent) error) Event
	OnAdminAuthRequest(handler func(evt *AdminAuthEvent) error) Event
	OnAdminBeforeAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) Event
	OnAdminAfterAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) Event
	OnAdminBeforeAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) Event
	OnAdminAfterAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) Event
	OnAdminBeforeRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) Event
	OnAdminAfterRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) Event
	OnAdminBeforeConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) Event
	OnAdminAfterConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) Event
	OnRecordAuthRequest(tags []string, handler func(evt *RecordAuthEvent) error) Event
	OnRecordBeforeAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) Event
	OnRecordAfterAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) Event
	OnRecordBeforeAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) Event
	OnRecordAfterAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) Event
	OnRecordBeforeAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) Event
	OnRecordAfterAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) Event
	OnRecordListExternalAuthsRequest(tags []string, handler func(evt *RecordListExternalAuthsEvent) error) Event
	OnRecordBeforeUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) Event
	OnRecordAfterUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) Event
	OnRecordBeforeRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) Event
	OnRecordAfterRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) Event
	OnRecordBeforeConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) Event
	OnRecordAfterConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) Event
	OnRecordBeforeRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) Event
	OnRecordAfterRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) Event
	OnRecordBeforeConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) Event
	OnRecordAfterConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) Event
	OnRecordBeforeRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) Event
	OnRecordAfterRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) Event
	OnRecordBeforeConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) Event
	OnRecordAfterConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) Event
	OnRecordsListRequest(tags []string, handler func(evt *RecordsListEvent) error) Event
	OnRecordViewRequest(tags []string, handler func(evt *RecordViewEvent) error) Event
	OnRecordBeforeCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) Event
	OnRecordAfterCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) Event
	OnRecordBeforeUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) Event
	OnRecordAfterUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) Event
	OnRecordBeforeDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) Event
	OnRecordAfterDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) Event
	{{- else if eq .BackendPkg "chi"}}
	OnBeforeBootstrap(handler func(evt *BootstrapEvent) error) Event
	OnAfterBootstrap(handler func(evt *BootstrapEvent) error) Event
	OnBeforeServe(handler func(evt *ServeEvent) error) Event
	OnBeforeApiError(handler func(evt *ApiErrorEvent) error) Event
	OnAfterApiError(handler func(evt *ApiErrorEvent) error) Event
	OnTerminate(handler func(evt *TerminateEvent) error) Event
	OnModelBeforeCreate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterCreate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelBeforeUpdate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterUpdate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelBeforeDelete(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterDelete(tags []string, handler func(evt *ModelEvent) error) Event
	OnMailerBeforeAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) Event
	OnMailerAfterAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) Event
	OnMailerBeforeRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerBeforeRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerBeforeRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnRealtimeConnectRequest(handler func(evt *RealtimeConnectEvent) error) Event
	OnRealtimeDisconnectRequest(handler func(evt *RealtimeDisconnectEvent) error) Event
	OnRealtimeBeforeMessageSend(handler func(evt *RealtimeMessageEvent) error) Event
	OnRealtimeAfterMessageSend(handler func(evt *RealtimeMessageEvent) error) Event
	OnRealtimeBeforeSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) Event
	OnRealtimeAfterSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) Event
	OnSettingsListRequest(handler func(evt *SettingsListEvent) error) Event
	OnSettingsBeforeUpdateRequest(handler func(evt *SettingsUpdateEvent) error) Event
	OnSettingsAfterUpdateRequest(handler func(evt *SettingsUpdateEvent) error) Event
	OnFileDownloadRequest(tags []string, handler func(evt *FileDownloadEvent) error) Event
	OnFileBeforeTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) Event
	OnFileAfterTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) Event
	OnAdminsListRequest(handler func(evt *AdminsListEvent) error) Event
	OnAdminViewRequest(handler func(evt *AdminViewEvent) error) Event
	OnAdminBeforeCreateRequest(handler func(evt *AdminCreateEvent) error) Event
	OnAdminAfterCreateRequest(handler func(evt *AdminCreateEvent) error) Event
	OnAdminBeforeUpdateRequest(handler func(evt *AdminUpdateEvent) error) Event
	OnAdminAfterUpdateRequest(handler func(evt *AdminUpdateEvent) error) Event
	OnAdminBeforeDeleteRequest(handler func(evt *AdminDeleteEvent) error) Event
	OnAdminAfterDeleteRequest(handler func(evt *AdminDeleteEvent) error) Event
	OnAdminAuthRequest(handler func(evt *AdminAuthEvent) error) Event
	OnAdminBeforeAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) Event
	OnAdminAfterAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) Event
	OnAdminBeforeAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) Event
	OnAdminAfterAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) Event
	OnAdminBeforeRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) Event
	OnAdminAfterRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) Event
	OnAdminBeforeConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) Event
	OnAdminAfterConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) Event
	OnRecordAuthRequest(tags []string, handler func(evt *RecordAuthEvent) error) Event
	OnRecordBeforeAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) Event
	OnRecordAfterAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) Event
	OnRecordBeforeAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) Event
	OnRecordAfterAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) Event
	OnRecordBeforeAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) Event
	OnRecordAfterAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) Event
	OnRecordListExternalAuthsRequest(tags []string, handler func(evt *RecordListExternalAuthsEvent) error) Event
	OnRecordBeforeUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) Event
	OnRecordAfterUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) Event
	OnRecordBeforeRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) Event
	OnRecordAfterRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) Event
	OnRecordBeforeConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) Event
	OnRecordAfterConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) Event
	OnRecordBeforeRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) Event
	OnRecordAfterRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) Event
	OnRecordBeforeConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) Event
	OnRecordAfterConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) Event
	OnRecordBeforeRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) Event
	OnRecordAfterRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) Event
	OnRecordBeforeConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) Event
	OnRecordAfterConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) Event
	OnRecordsListRequest(tags []string, handler func(evt *RecordsListEvent) error) Event
	OnRecordViewRequest(tags []string, handler func(evt *RecordViewEvent) error) Event
	OnRecordBeforeCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) Event
	OnRecordAfterCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) Event
	OnRecordBeforeUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) Event
	OnRecordAfterUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) Event
	OnRecordBeforeDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) Event
	OnRecordAfterDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) Event
	{{- else if eq .BackendPkg "stdlib"}}
	OnBeforeBootstrap(handler func(evt *BootstrapEvent) error) Event
	OnAfterBootstrap(handler func(evt *BootstrapEvent) error) Event
	OnBeforeServe(handler func(evt *ServeEvent) error) Event
	OnBeforeApiError(handler func(evt *ApiErrorEvent) error) Event
	OnAfterApiError(handler func(evt *ApiErrorEvent) error) Event
	OnTerminate(handler func(evt *TerminateEvent) error) Event
	OnModelBeforeCreate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterCreate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelBeforeUpdate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterUpdate(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelBeforeDelete(tags []string, handler func(evt *ModelEvent) error) Event
	OnModelAfterDelete(tags []string, handler func(evt *ModelEvent) error) Event
	OnMailerBeforeAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) Event
	OnMailerAfterAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) Event
	OnMailerBeforeRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerBeforeRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerBeforeRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnMailerAfterRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) Event
	OnRealtimeConnectRequest(handler func(evt *RealtimeConnectEvent) error) Event
	OnRealtimeDisconnectRequest(handler func(evt *RealtimeDisconnectEvent) error) Event
	OnRealtimeBeforeMessageSend(handler func(evt *RealtimeMessageEvent) error) Event
	OnRealtimeAfterMessageSend(handler func(evt *RealtimeMessageEvent) error) Event
	OnRealtimeBeforeSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) Event
	OnRealtimeAfterSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) Event
	OnSettingsListRequest(handler func(evt *SettingsListEvent) error) Event
	OnSettingsBeforeUpdateRequest(handler func(evt *SettingsUpdateEvent) error) Event
	OnSettingsAfterUpdateRequest(handler func(evt *SettingsUpdateEvent) error) Event
	OnFileDownloadRequest(tags []string, handler func(evt *FileDownloadEvent) error) Event
	OnFileBeforeTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) Event
	OnFileAfterTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) Event
	OnAdminsListRequest(handler func(evt *AdminsListEvent) error) Event
	OnAdminViewRequest(handler func(evt *AdminViewEvent) error) Event
	OnAdminBeforeCreateRequest(handler func(evt *AdminCreateEvent) error) Event
	OnAdminAfterCreateRequest(handler func(evt *AdminCreateEvent) error) Event
	OnAdminBeforeUpdateRequest(handler func(evt *AdminUpdateEvent) error) Event
	OnAdminAfterUpdateRequest(handler func(evt *AdminUpdateEvent) error) Event
	OnAdminBeforeDeleteRequest(handler func(evt *AdminDeleteEvent) error) Event
	OnAdminAfterDeleteRequest(handler func(evt *AdminDeleteEvent) error) Event
	OnAdminAuthRequest(handler func(evt *AdminAuthEvent) error) Event
	OnAdminBeforeAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) Event
	OnAdminAfterAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) Event
	OnAdminBeforeAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) Event
	OnAdminAfterAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) Event
	OnAdminBeforeRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) Event
	OnAdminAfterRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) Event
	OnAdminBeforeConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) Event
	OnAdminAfterConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) Event
	OnRecordAuthRequest(tags []string, handler func(evt *RecordAuthEvent) error) Event
	OnRecordBeforeAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) Event
	OnRecordAfterAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) Event
	OnRecordBeforeAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) Event
	OnRecordAfterAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) Event
	OnRecordBeforeAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) Event
	OnRecordAfterAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) Event
	OnRecordListExternalAuthsRequest(tags []string, handler func(evt *RecordListExternalAuthsEvent) error) Event
	OnRecordBeforeUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) Event
	OnRecordAfterUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) Event
	OnRecordBeforeRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) Event
	OnRecordAfterRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) Event
	OnRecordBeforeConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) Event
	OnRecordAfterConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) Event
	OnRecordBeforeRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) Event
	OnRecordAfterRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) Event
	OnRecordBeforeConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) Event
	OnRecordAfterConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) Event
	OnRecordBeforeRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) Event
	OnRecordAfterRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) Event
	OnRecordBeforeConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) Event
	OnRecordAfterConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) Event
	OnRecordsListRequest(tags []string, handler func(evt *RecordsListEvent) error) Event
	OnRecordViewRequest(tags []string, handler func(evt *RecordViewEvent) error) Event
	OnRecordBeforeCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) Event
	OnRecordAfterCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) Event
	OnRecordBeforeUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) Event
	OnRecordAfterUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) Event
	OnRecordBeforeDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) Event
	OnRecordAfterDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) Event
	{{- end }}
}

            `
		},

		"app/types/core/app.go": func() string {
			return `package core
import (
	{{- if eq .BackendPkg "echo"}}
	"github.com/labstack/echo/v5"
	{{- else if eq .BackendPkg "gin"}}
	"github.com/gin-gonic/gin"
	{{- else if eq .BackendPkg "chi"}}
	"github.com/go-chi/chi/v5"
	{{- else if eq .BackendPkg "stdlib"}}
	"net/http"
	{{- end }}
)
type App interface {
// App-related methods
Dao() *daos.Dao
LogsDao() *daos.Dao
Logger() *slog.Logger
DataDir() string
EncryptionEnv() string
IsDev() bool
Settings() *settings.Settings
Store() *store.Store[any]
SubscriptionsBroker() *subscriptions.Broker
NewMailClient() mailer.Mailer
NewFilesystem() (*filesystem.System, error)
NewBackupsFilesystem() (*filesystem.System, error)
RefreshSettings() error
IsBootstrapped() bool
Bootstrap() error
ResetBootstrapState() error
CreateBackup(ctx context.Context, name string) error
RestoreBackup(ctx context.Context, name string) error
Restart() error

// Event hooks
OnBeforeBootstrap(handler func(evt *BootstrapEvent) error) App
OnAfterBootstrap(handler func(evt *BootstrapEvent) error) App
OnBeforeServe(handler func(evt *ServeEvent) error) App
OnBeforeApiError(handler func(evt *ApiErrorEvent) error) App
OnAfterApiError(handler func(evt *ApiErrorEvent) error) App
OnTerminate(handler func(evt *TerminateEvent) error) App

// DAO event hooks
OnModelBeforeCreate(tags []string, handler func(evt *ModelEvent) error) App
OnModelAfterCreate(tags []string, handler func(evt *ModelEvent) error) App
OnModelBeforeUpdate(tags []string, handler func(evt *ModelEvent) error) App
OnModelAfterUpdate(tags []string, handler func(evt *ModelEvent) error) App
OnModelBeforeDelete(tags []string, handler func(evt *ModelEvent) error) App
OnModelAfterDelete(tags []string, handler func(evt *ModelEvent) error) App

// Mailer event hooks
OnMailerBeforeAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) App
OnMailerAfterAdminResetPasswordSend(handler func(evt *MailerAdminEvent) error) App
OnMailerBeforeRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) App
OnMailerAfterRecordResetPasswordSend(tags []string, handler func(evt *MailerRecordEvent) error) App
OnMailerBeforeRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) App
OnMailerAfterRecordVerificationSend(tags []string, handler func(evt *MailerRecordEvent) error) App
OnMailerBeforeRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) App
OnMailerAfterRecordChangeEmailSend(tags []string, handler func(evt *MailerRecordEvent) error) App

// Realtime API event hooks
OnRealtimeConnectRequest(handler func(evt *RealtimeConnectEvent) error) App
OnRealtimeDisconnectRequest(handler func(evt *RealtimeDisconnectEvent) error) App
OnRealtimeBeforeMessageSend(handler func(evt *RealtimeMessageEvent) error) App
OnRealtimeAfterMessageSend(handler func(evt *RealtimeMessageEvent) error) App
OnRealtimeBeforeSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) App
OnRealtimeAfterSubscribeRequest(handler func(evt *RealtimeSubscribeEvent) error) App

// Settings API event hooks
OnSettingsListRequest(handler func(evt *SettingsListEvent) error) App
OnSettingsBeforeUpdateRequest(handler func(evt *SettingsUpdateEvent) error) App
OnSettingsAfterUpdateRequest(handler func(evt *SettingsUpdateEvent) error) App

// File API event hooks
OnFileDownloadRequest(tags []string, handler func(evt *FileDownloadEvent) error) App
OnFileBeforeTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) App
OnFileAfterTokenRequest(tags []string, handler func(evt *FileTokenEvent) error) App

// Admin API event hooks
OnAdminsListRequest(handler func(evt *AdminsListEvent) error) App
OnAdminViewRequest(handler func(evt *AdminViewEvent) error) App
OnAdminBeforeCreateRequest(handler func(evt *AdminCreateEvent) error) App
OnAdminAfterCreateRequest(handler func(evt *AdminCreateEvent) error) App
OnAdminBeforeUpdateRequest(handler func(evt *AdminUpdateEvent) error) App
OnAdminAfterUpdateRequest(handler func(evt *AdminUpdateEvent) error) App
OnAdminBeforeDeleteRequest(handler func(evt *AdminDeleteEvent) error) App
OnAdminAfterDeleteRequest(handler func(evt *AdminDeleteEvent) error) App
OnAdminAuthRequest(handler func(evt *AdminAuthEvent) error) App
OnAdminBeforeAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) App
OnAdminAfterAuthWithPasswordRequest(handler func(evt *AdminAuthWithPasswordEvent) error) App
OnAdminBeforeAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) App
OnAdminAfterAuthRefreshRequest(handler func(evt *AdminAuthRefreshEvent) error) App
OnAdminBeforeRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) App
OnAdminAfterRequestPasswordResetRequest(handler func(evt *AdminRequestPasswordResetEvent) error) App
OnAdminBeforeConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) App
OnAdminAfterConfirmPasswordResetRequest(handler func(evt *AdminConfirmPasswordResetEvent) error) App

// Record Auth API event hooks
OnRecordAuthRequest(tags []string, handler func(evt *RecordAuthEvent) error) App
OnRecordBeforeAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) App
OnRecordAfterAuthWithPasswordRequest(tags []string, handler func(evt *RecordAuthWithPasswordEvent) error) App
OnRecordBeforeAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) App
OnRecordAfterAuthWithOAuth2Request(tags []string, handler func(evt *RecordAuthWithOAuth2Event) error) App
OnRecordBeforeAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) App
OnRecordAfterAuthRefreshRequest(tags []string, handler func(evt *RecordAuthRefreshEvent) error) App
OnRecordListExternalAuthsRequest(tags []string, handler func(evt *RecordListExternalAuthsEvent) error) App
OnRecordBeforeUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) App
OnRecordAfterUnlinkExternalAuthRequest(tags []string, handler func(evt *RecordUnlinkExternalAuthEvent) error) App
OnRecordBeforeRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) App
OnRecordAfterRequestPasswordResetRequest(tags []string, handler func(evt *RecordRequestPasswordResetEvent) error) App
OnRecordBeforeConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) App
OnRecordAfterConfirmPasswordResetRequest(tags []string, handler func(evt *RecordConfirmPasswordResetEvent) error) App
OnRecordBeforeRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) App
OnRecordAfterRequestVerificationRequest(tags []string, handler func(evt *RecordRequestVerificationEvent) error) App
OnRecordBeforeConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) App
OnRecordAfterConfirmVerificationRequest(tags []string, handler func(evt *RecordConfirmVerificationEvent) error) App
OnRecordBeforeRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) App
OnRecordAfterRequestEmailChangeRequest(tags []string, handler func(evt *RecordRequestEmailChangeEvent) error) App
OnRecordBeforeConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) App
OnRecordAfterConfirmEmailChangeRequest(tags []string, handler func(evt *RecordConfirmEmailChangeEvent) error) App

// Record CRUD API event hooks
OnRecordsListRequest(tags []string, handler func(evt *RecordsListEvent) error) App
OnRecordViewRequest(tags []string, handler func(evt *RecordViewEvent) error) App
OnRecordBeforeCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) App
OnRecordAfterCreateRequest(tags []string, handler func(evt *RecordCreateEvent) error) App
OnRecordBeforeUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) App
OnRecordAfterUpdateRequest(tags []string, handler func(evt *RecordUpdateEvent) error) App
OnRecordBeforeDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) App
OnRecordAfterDeleteRequest(tags []string, handler func(evt *RecordDeleteEvent) error) App

// Route management methods
AddRoute(method, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) App
AddRouteGroup(prefix string, middlewares []echo.MiddlewareFunc, routes []Route) App
AddResource(resource string, controller interface{}, middlewares ...echo.MiddlewareFunc) App
Serve() App
} 

type ServeEvent struct {
    {{- if eq .BackendPkg "echo"}}
    Router *echo.Echo
    {{- else if eq .BackendPkg "gin"}}
    Router *gin.Engine
    {{- else if eq .BackendPkg "chi"}}
    Router chi.Router
    {{- else if eq .BackendPkg "stdlib"}}
    Router *http.ServeMux
    {{- end }}
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
