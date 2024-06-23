package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theHamdiz/gost/cleaner"
	"github.com/theHamdiz/gost/gen/config"
	"github.com/theHamdiz/gost/parser"
)

func Map() map[string]func() string {
	return map[string]func() string{
		"plugins/core.go": func() string {
			return `package plugins
type Plugin interface {
    Init() error
    Name() string
}
`
		},
		"plugins/orm.go": func() string {
			return `package plugins
import (
    "fmt"
    "strings"
)

// State represents the state of the query builder
type State int

const (
    Initial State = iota
    Selecting
    Froming
    Whereing
    Inserting
    Valuing
    Updating
    Setting
    Deleting
)

// PostgreSQLDialect is an implementation of Dialect for PostgreSQL
type PostgreSQLDialect struct{}

func (d *PostgreSQLDialect) Select(columns ...string) string {
    return fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
}

func (d *PostgreSQLDialect) From(table string) string {
    return fmt.Sprintf("FROM %s ", table)
}

func (d *PostgreSQLDialect) Where(condition string) string {
    return fmt.Sprintf("WHERE %s ", condition)
}

func (d *PostgreSQLDialect) InsertInto(table string, columns ...string) string {
    return fmt.Sprintf("INSERT INTO %s (%s) ", table, strings.Join(columns, ", "))
}

func (d *PostgreSQLDialect) Values(values ...interface{}) string {
    valStr := make([]string, len(values))
    for i, v := range values {
        valStr[i] = fmt.Sprintf("'%v'", v)
    }
    return fmt.Sprintf("VALUES (%s) ", strings.Join(valStr, ", "))
}

func (d *PostgreSQLDialect) Update(table string) string {
    return fmt.Sprintf("UPDATE %s ", table)
}

func (d *PostgreSQLDialect) Set(assignments ...string) string {
    return fmt.Sprintf("SET %s ", strings.Join(assignments, ", "))
}

func (d *PostgreSQLDialect) DeleteFrom(table string) string {
    return fmt.Sprintf("DELETE FROM %s ", table)
}

// MySQLDialect is an implementation of Dialect for MySQL
type MySQLDialect struct{}

func (d *MySQLDialect) Select(columns ...string) string {
    return fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
}

func (d *MySQLDialect) From(table string) string {
    return fmt.Sprintf("FROM %s ", table)
}

func (d *MySQLDialect) Where(condition string) string {
    return fmt.Sprintf("WHERE %s ", condition)
}

func (d *MySQLDialect) InsertInto(table string, columns ...string) string {
    return fmt.Sprintf("INSERT INTO %s (%s) ", table, strings.Join(columns, ", "))
}

func (d *MySQLDialect) Values(values ...interface{}) string {
    valStr := make([]string, len(values))
    for i, v := range values {
        valStr[i] = fmt.Sprintf("'%v'", v)
    }
    return fmt.Sprintf("VALUES (%s) ", strings.Join(valStr, ", "))
}

func (d *MySQLDialect) Update(table string) string {
    return fmt.Sprintf("UPDATE %s ", table)
}

func (d *MySQLDialect) Set(assignments ...string) string {
    return fmt.Sprintf("SET %s ", strings.Join(assignments, ", "))
}

func (d *MySQLDialect) DeleteFrom(table string) string {
    return fmt.Sprintf("DELETE FROM %s ", table)
}

type ORMPlugin interface {
    Plugin
    NewORMBuilder(dialect Dialect) *ORMBuilder
}

type ORM struct{}

func (o *ORM) Init() error {
    // Initialization logic if needed
    return nil
}

func (o *ORM) Name() string {
    return "ORM"
}

func (o *ORM) NewORMBuilder(dialect orm.Dialect) *orm.ORMBuilder {
    return orm.NewORMBuilder(dialect)
}

type ORMBuilder struct {
    dialect Dialect
    query   strings.Builder
    state   State
}

// NewORMBuilder creates a new ORMBuilder for a given dialect
func NewORMBuilder(dialect Dialect) *ORMBuilder {
    return &ORMBuilder{
        dialect: dialect,
        state:   Initial,
    }
}

// Select adds a SELECT clause to the query
func (b *ORMBuilder) Select(columns ...string) *ORMBuilder {
    if b.state != Initial {
        panic("Select can only be called at the beginning of the query")
    }
    b.query.WriteString(b.dialect.Select(columns...))
    b.state = Selecting
    return b
}

// From adds a FROM clause to the query
func (b *ORMBuilder) From(table string) *ORMBuilder {
    if b.state != Selecting {
        panic("From must be called after Select")
    }
    b.query.WriteString(b.dialect.From(table))
    b.state = Froming
    return b
}

// Where adds a WHERE clause to the query
func (b *ORMBuilder) Where(condition string) *ORMBuilder {
    if b.state != Froming && b.state != Updating {
        panic("Where must be called after From or Update")
    }
    b.query.WriteString(b.dialect.Where(condition))
    b.state = Whereing
    return b
}

// InsertInto adds an INSERT INTO clause to the query
func (b *ORMBuilder) InsertInto(table string, columns ...string) *ORMBuilder {
    if b.state != Initial {
        panic("InsertInto can only be called at the beginning of the query")
    }
    b.query.WriteString(b.dialect.InsertInto(table, columns...))
    b.state = Inserting
    return b
}

// Values adds a VALUES clause to the query
func (b *ORMBuilder) Values(values ...interface{}) *ORMBuilder {
    if b.state != Inserting {
        panic("Values must be called after InsertInto")
    }
    b.query.WriteString(b.dialect.Values(values...))
    b.state = Valuing
    return b
}

// Update adds an UPDATE clause to the query
func (b *ORMBuilder) Update(table string) *ORMBuilder {
    if b.state != Initial {
        panic("Update can only be called at the beginning of the query")
    }
    b.query.WriteString(b.dialect.Update(table))
    b.state = Updating
    return b
}

// Set adds a SET clause to the query
func (b *ORMBuilder) Set(assignments ...string) *ORMBuilder {
    if b.state != Updating {
        panic("Set must be called after Update")
    }
    b.query.WriteString(b.dialect.Set(assignments...))
    b.state = Setting
    return b
}

// DeleteFrom adds a DELETE FROM clause to the query
func (b *ORMBuilder) DeleteFrom(table string) *ORMBuilder {
    if b.state != Initial {
        panic("DeleteFrom can only be called at the beginning of the query")
    }
    b.query.WriteString(b.dialect.DeleteFrom(table))
    b.state = Deleting
    return b
}

// Build returns the final SQL query as a string
func (b *ORMBuilder) Build() string {
    return b.query.String()
}
`
		},
		"app/events/events.go": func() string {
			return `package events

import (
    "context"
    "log"
    "sync"
    "time"
)

type Event struct {
    Name string
    Data interface{}
}

type EventHandler func(context.Context, Event) error

type Subscription struct {
    CreatedAt int64
    EventName string
    Handler   EventHandler
}

type EventManager struct {
    mu        sync.RWMutex
    listeners map[string][]Subscription
    eventCh   chan Event
    quitCh    chan struct{}
}

func NewEventManager() *EventManager {
    em := &EventManager{
        listeners: make(map[string][]Subscription),
        eventCh:   make(chan Event, 128),
        quitCh:    make(chan struct{}),
    }
    go em.start()
    return em
}

func (em *EventManager) start() {
    ctx := context.Background()
    for {
        select {
        case {{ "<-" }} em.quitCh:
            return
        case event := {{ "<-" }}em.eventCh:
            if handlers, found := em.listeners[event.Name]; found {
                for _, sub := range handlers {
                    go func(sub Subscription, event Event) {
                        ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
                        defer cancel()
                        start := time.Now()
                        if err := sub.Handler(ctx, event); err != nil {
                            log.Printf("Error handling event %s: %v", event.Name, err)
                        }
                        log.Printf("Handled event %s in %v", event.Name, time.Since(start))
                    }(sub, event)
                }
            }
        }
    }
}

func (em *EventManager) Stop() {
    em.quitCh {{ "<-" }} struct{}{}
}

func (em *EventManager) RegisterListener(eventName string, handler EventHandler) Subscription {
    em.mu.Lock()
    defer em.mu.Unlock()

    sub := Subscription{
        CreatedAt: time.Now().UnixNano(),
        EventName: eventName,
        Handler:   handler,
    }

    em.listeners[eventName] = append(em.listeners[eventName], sub)

    return sub
}

func (em *EventManager) UnregisterListener(sub Subscription) {
    em.mu.Lock()
    defer em.mu.Unlock()

    if handlers, found := em.listeners[sub.EventName]; found {
        for i, s := range handlers {
            if s.CreatedAt == sub.CreatedAt {
                em.listeners[sub.EventName] = append(handlers[:i], handlers[i+1:]...)
                break
            }
        }
        if len(em.listeners[sub.EventName]) == 0 {
            delete(em.listeners, sub.EventName)
        }
    }
}

func (em *EventManager) Emit(event Event) {
    em.eventCh {{ "<-" }} event
}
`
		},
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

//go:embed frontend/dist/* backend/dist/*
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
		"gost/prelude.go": func() string {
			return `package prelude

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
    g.Text(http.StatusInternalServerError, err.Error())
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

		"cmd/app/main.go": func() string {
			return `package main

import (
	{{- if eq .BackendPkg "echo" }}
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	{{- end }}
    "{{.AppName}}/app/cfg"
    "{{.AppName}}/app/db"
    . "{{.AppName}}/app/middleware"
    {{- if eq .BackendPkg "chi" }}
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    {{- end }}
    {{- if eq .BackendPkg "gin" }}
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/middleware"
    {{- end }}
    "log"
	{{- if eq .BackendPkg "std" }}
    "net/http"
	{{- end }}
	"os"
    "os/signal"
    "syscall"
    "{{.AppName}}/app/router"
)

func waitForShutdown() {
    // Create a channel to receive OS signals
    sigs := make(chan os.Signal, 1)
    // Notify the channel of SIGINT and SIGTERM signals
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    // Block until a signal is received
    sig := {{"<-"}}sigs
    log.Printf("Received signal: %s", sig)
}

func main() {
    c := cfg.LoadConfig()
    logger.InitLogger()
    db.InitDB(c)
    defer db.CloseDB()

    server := router.InitRoutes()

    log.Println("Server starting on port", cfg.Port)
    go func() {
        {{- if eq .BackendPkg "gin" }}
        if err := server.(*gin.Engine).Run(":" + cfg.Port); err != nil {
            log.Fatal(err)
        }
        {{- end }}
        {{- if eq .BackendPkg "echo" }}
        if err := server.(*echo.Echo).Start(":" + cfg.Port); err != nil {
            log.Fatal(err)
        }
        {{- end }}
        {{- if eq .BackendPkg "chi" }}
        if err := http.ListenAndServe(":" + cfg.Port, server.(http.Handler)); err != nil {
            log.Fatal(err)
        }
        {{- end }}
        {{- if eq .BackendPkg "std" }}
        if err := http.ListenAndServe(":" + cfg.Port, server.(http.Handler)); err != nil {
            log.Fatal(err)
        }
        {{- end }}
    }()

    waitForShutdown()
}
`
		},
		"app/services/logger.go": func() string {
			return `package services

import (
    "log"
    "os"
)

var (
    InfoLogger  *log.Logger
    ErrorLogger *log.Logger
)

func InitLogger() {
    InfoLogger = log.New(os.Stdout, "ℹ INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(os.Stderr, "✗ ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
    InfoLogger.Println(message)
}

func Error(message string) {
    ErrorLogger.Println(message)
}
`
		},
		"app/services/rateLimiter.go": func() string {
			return `package services

import (
    "golang.org/x/time/rate"
)

var limiter *rate.Limiter

func InitRateLimiter(limit rate.Limit, burst int) {
    limiter = rate.NewLimiter(limit, burst)
}

func Allow() bool {
    return limiter.Allow()
}
`
		},
		"app/services/db.go": func() string {
			return `package db

import (
    "database/sql"
    "log"
    {{- if eq .DbDriver "sqlite3" }}
    _ "github.com/mattn/go-sqlite3"
    {{- end }}
)

var DB *sql.DB

func InitDB(c cfg.Config) {
    var err error
    DB, err = sql.Open(c.DbDriver, c.DbName)
    if err != nil {
        log.Fatal("Error opening database:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Error pinging database:", err)
    }
}

func CloseDB() {
    if err := DB.Close(); err != nil {
        log.Println("Error closing database:", err)
    }
}
`
		},
		"app/cfg/cfg.go": func() string {
			return `package cfg

	import (
    "os"
)

type Config struct {
    GostEnv                      string
    Port                         string
    DbDriver                     string
    DbUser                       string
    DbHost                       string
    DbPassword                   string
    DbName                       string
    DbOrm						 string
    MigrationsDir                string
    GostSecret                   string
    GostAuthRedirectAfterLogin   string
    GostAuthSessionExpiryInHours string
    GostAuthSkipVerify           bool
    GostBackend                  string
}

func (c *Config) IsDevelopment(){
	strings.ToLower(c.GostEnv) == "dev" || strings.ToLower(c.GostEnv) == "development"
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value, exists := os.LookupEnv(key); exists {
        return value == "true"
    }
    return defaultValue
}

func LoadFromEnv() *Config {
    return &Config{
        GostEnv:                      getEnv("GOST_ENV", "DEV"),
        Port:                         getEnv("PORT", ":9630"),
        DbDriver:                     getEnv("DB_DRIVER", "sqlite3"),
        DbUser:                       getEnv("DB_USER", ""),
        DbHost:                       getEnv("DB_HOST", ""),
        DbPassword:                   getEnv("DB_PASSWORD", ""),
        DbName:                       getEnv("DB_NAME", "data.db"),
        DbOrm:						  getEnv("DB_ORM", "entgo"),
        MigrationsDir:                getEnv("MIGRATIONS_DIR", "app/db/migrations"),
        GostSecret:                   getEnv("GOST_SECRET", "49cf26a7d274d62ad902ead6e69f5d71b4ffe703b4b07d25652c117cab74fcb1"),
        GostAuthRedirectAfterLogin:   getEnv("GOST_AUTH_REDIRECT_AFTER_LOGIN", "/profile"),
        GostAuthSessionExpiryInHours: getEnv("GOST_AUTH_SESSION_EXPIRY_IN_HOURS", "72"),
        GostAuthSkipVerify:           getEnvBool("GOST_AUTH_SKIP_VERIFY", true),
        Backend:                      getEnv("GOST_BACKEND", "gin"),
    }
}
`
		},
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

		"app/handlers/auth.go": func() string {
			return `package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))

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
		"app/handlers/landing.go": func() string {
			return `package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "home.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

`
		},
		"app/handlers/about.go": func() string {
			return `package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var views = template.Must(template.ParseGlob(filepath.Join("app", "views", "*.templ")))

func AboutHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "about.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
`
		},
		"app/db/db.go": func() string {
			return `package db

import (
		"log"
		"os"

		"github.com/theHamdiz/gost/cfg"

		{{- if eq .DbDriver "sqlite3" }}
		_ "github.com/mattn/go-sqlite3"
		{{- end }}
		"github.com/uptrace/bun"
		"github.com/uptrace/bun/dialect/sqlitedialect"
		"github.com/uptrace/bun/extra/bundebug"
)
var Query *bun.DB

func init() {
		cfg, err := cfg.LoadFromEnv()
		if err != nil {
			log.Fatal(err)
		}
		db, err := db.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		Query = bun.NewDB(db, sqlitedialect.New())
		if cfg.IsDevelopment() {
			Query.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		}
}
`
		},
		"app/middleware/auth.go": func() string {
			return `package middleware

import (
    "net/http"
)

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add your authentication logic here
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        // Validate the token
        // ...
        next.ServeHTTP(w, r)
    })
}
`
		},
		"app/middleware/cors.go": func() string {
			return `package middleware

import (
    "net/http"
)

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
`
		},
		"app/middleware/logger.go": func() string {
			return `package middleware

import (
    "log"
    "net/http"
    "time"
)

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
    })
}
`
		},
		"app/middleware/notifier.go": func() string {
			return `package middleware

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/go-mail/mail"
)

var (
    clients       = make(map[chan string]bool)
    notifyChannel = make(chan string)
)

// Notifier middleware that sends SSE notifications to connected clients
func Notifier(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        flusher, ok := w.(http.Flusher)
        if !ok {
            http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
            return
        }

        messageChan := make(chan string)
        clients[messageChan] = true

        defer func() {
            delete(clients, messageChan)
            close(messageChan)
        }()

        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")

        notifyChannel {{"<-"}} "System started"

        for {
            select {
            case {{"<-"}}r.Context().Done():
                return
            case msg := {{"<-"}}messageChan:
                fmt.Fprintf(w, "data: %s\n\n", msg)
                flusher.Flush()
            }
        }
    })
}

// Function to notify all connected SSE clients
func notifyClients(message string) {
    for client := range clients {
        client {{"<-"}} message
    }
}

// Function to notify users via email
func notifyByEmail(subject, body string) {
    m := mail.NewMessage()
    m.SetHeader("From", "your-email@example.com")
    m.SetHeader("To", "user@example.com") // add your recipient's email here
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    d := mail.NewDialer("smtp.example.com", 587, "your-username", "your-password")

    if err := d.DialAndSend(m); err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}
`
		},
		"app/middleware/recoverer.go": func() string {
			return `package middleware

import (
    "log"
    "net/http"
    "runtime/debug"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if err := recover(); err != nil {
            log.Printf("Recovered from panic: %v", err)
            notifyClients("System shutdown unexpectedly")
            notifyByEmail("System Shutdown", fmt.Sprintf("The system was shut down unexpectedly: %v", err))
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    }()
    next.ServeHTTP(w, r)
}
`
		},
		"app/middleware/rateLimiter.go": func() string {
			return `package middleware

import (
    "net/http"
    "time"

    "golang.org/x/time/rate"
)

func RateLimiter(limit rate.Limit, burst int) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(limit, burst)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
`
		},
		"app/middleware/requestId.go": func() string {
			return `package middleware

import (
    "context"
    "net/http"
    "github.com/google/uuid"
)

type key int

const requestIDKey key = 0

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := uuid.New().String()
        ctx := context.WithValue(r.Context(), requestIDKey, id)
        w.Header().Set("X-Request-ID", id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetRequestID(r *http.Request) string {
    if id, ok := r.Context().Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}
`
		},
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
		"public/public.go": func() string {
			return `package public

import "embed"

//go:embed assets
var AssetsFS embed.FS
`
		},
		"go.mod": func() string {
			return `module {{.AppName}}

go 1.22.4

require {{.VersionedBackendImport}}
`
		},
		"go.sum": func() string {
			return ``
		},
		"README.md": func() string {
			return "# {{ .AppName }}\n\n" +
				"A brief description of what your project does.\n\n" +
				"## Features\n\n" +
				"- Feature 1\n" +
				"- Feature 2\n" +
				"- Feature 3\n\n" +
				"## Installation\n\n" +
				"To install and run this project, follow these steps:\n\n" +
				"1. Clone the repository:\n\n" +
				"```sh\n" +
				"git clone https://github.com/yourusername/yourproject.git\n" +
				"cd yourproject\n" +
				"```\n\n" +
				"2. Install dependencies:\n\n" +
				"```sh\n" +
				"go mod tidy\n" +
				"```\n\n" +
				"3. Set up environment variables (if any):\n\n" +
				"```sh\n" +
				"cp .env.example .env\n" +
				"# Edit the .env file with your configuration\n" +
				"```\n\n" +
				"4. Run the application:\n\n" +
				"```sh\n" +
				"go run main.go\n" +
				"```\n\n" +
				"## Usage\n\n" +
				"### Running the Project\n\n" +
				"To start the project, use:\n\n" +
				"```sh\n" +
				"gost r\n" +
				"```\n\n" +
				"### Project Structure\n\n" +
				"By default gost creates the following structure for you:\n\n" +
				"```\n" +
				".\n" +
				"├── cmd             # Main applications of the project\n" +
				"├── app             # Private application and library code\n" +
				"├── pkg             # Public library code\n" +
				"├── web             # Web server-related files\n" +
				"│   ├── static      # Static files\n" +
				"│   └── templates   # HTML templates\n" +
				"├── go.mod          # Go module file\n" +
				"├── main.go         # Main entry point of the application\n" +
				"└── README.md       # This file\n" +
				"```\n\n" +
				"### Running Tests\n\n" +
				"To run tests, use:\n\n" +
				"```sh\n" +
				"go test ./...\n" +
				"```\n\n" +
				"## Configuration\n\n" +
				"List any configuration settings for your project:\n\n" +
				"- `DATABASE_URL`: The URL of your database.\n" +
				"- `PORT`: The port on which the server will run.\n\n" +
				"## Contributing\n\n" +
				"We welcome contributions! Please follow these steps to contribute:\n\n" +
				"1. Fork the repository.\n" +
				"2. Create a new branch with your feature or bug fix.\n" +
				"3. Commit your changes.\n" +
				"4. Push the branch to your fork.\n" +
				"5. Create a pull request.\n\n" +
				"## License\n\n" +
				"This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.\n\n" +
				"## Acknowledgements\n\n" +
				"Thanks to the contributors and the open-source community for their valuable input and support.\n"
		}, "MAKEFILE": func() string {
			return `# Makefile for {{.AppName}}

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME = {{.AppName}}
BINARY_UNIX = $(BINARY_NAME)_unix

# Frontend parameters
FRONTEND_DIR = web
NPMCMD = npm
NPMINSTALL = $(NPMCMD) install
NPMRUNBUILD = $(NPMCMD) run build

# All target
all: test build

# Test target
test:
	$(GOTEST) -v ./...

# Build target
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Release target
release: clean
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
	zip $(BINARY_UNIX).zip $(BINARY_UNIX)

# Frontend target
frontend:
	cd $(FRONTEND_DIR) && $(NPMINSTALL) && $(NPMRUNBUILD)

# Clean target
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_UNIX).zip

.PHONY: all test build release frontend clean
`
		},
		".gitignore": func() string {
			return ``
		},
		".air.toml": func() string {
			return `[build]
cmd = "go build -o ./tmp/main ."
bin = "tmp/main"
watch = ["."]
exclude_dir = ["tmp", "vendor"]
exclude_file = ["go.sum", "go.mod", ".gitignore", ".DS_Store", ".idea"]
delay = 200
`
		},
		".env": func() string {
			return `
# Application environment
# PROD or DEV
GOST_ENV=DEV

# HTTP listen port of the application
PORT=:{{.Port}}

# Database Config
DB_DRIVER={{.DbDriver}}
DB_USER=
DB_HOST=
DB_PASSWORD=
DB_NAME=db.db

MIGRATIONS_DIR=app/db/migrations

# Application secret used to secure your sessions.
# The secret will be auto generated on install.
# If you still want to change it make sure its at
# least 32 bytes long.
# NOTE: You might want to change this secret when using
# your app in production.
GOST_SECRET={{.Fingerprint}}

# Authentication Plugin
GOST_AUTH_REDIRECT_AFTER_LOGIN=/profile
GOST_AUTH_SESSION_EXPIRY_IN_HOURS=72
# Skip user email verification
GOST_AUTH_SKIP_VERIFY=true
GOST_BACKEND={{ .BackendPkg }}
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
		return nil, fmt.Errorf("> error getting current working directory: %v", err)
	}

	// Construct the full path to the asset
	assetPath := filepath.Join(cwd, "app", "assets", fileName)

	// Read the file content
	content, err := os.ReadFile(assetPath)
	if err != nil {
		return nil, fmt.Errorf("> error reading file %s: %v", assetPath, err)
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
	}
}

func GenerateFiles(data config.ProjectData) error {
	for path, tmplFunc := range Map() {
		content, err := parser.ParseTemplateString(path, tmplFunc(), data)
		if err != nil {
			fmt.Printf(">> Data Provided:\n%+v\n", data)
			fmt.Printf(">> Template String:\n%+v\n", tmplFunc())
			return fmt.Errorf(">> failed to parse template %s: %w", path, err)
		}

		appNameDowncased := strings.ToLower(data.AppName)
		filePath := filepath.Join(appNameDowncased, path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("✗ failed to create directory: %w", err)
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		} else {
			if strings.HasSuffix(filePath, ".go") {
				if err = cleaner.SortImports(filePath); err != nil {
					return fmt.Errorf("✗ The file was saved but failed to sort its imports %w", err)
				}
			}
		}
	}
	return nil
}
