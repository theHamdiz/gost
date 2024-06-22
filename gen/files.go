package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theHamdiz/gost/cleaner"
)

func FilesMap() map[string]func(data TemplateData) string {
	return map[string]func(data TemplateData) string{
		"cmd/app/main.go": func(data TemplateData) string {
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
    {{- end }}
    {{- if eq .BackendPkg "gin" }}
    "github.com/gin-gonic/gin"
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
    sig := <-sigs
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
		"app/services/logger.go": func(data TemplateData) string {
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
    InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
    InfoLogger.Println(message)
}

func Error(message string) {
    ErrorLogger.Println(message)
}
`
		},
		"app/services/rateLimiter.go": func(data TemplateData) string {
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
		"app/services/db.go": func(data TemplateData) string {
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
		"app/cfg/config.go": func(data TemplateData) string {
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
		"app/router/router.go": func(data TemplateData) string {
			return `package router

import (
    "{{.AppName}}/app/handlers"
    "{{.AppName}}/app/middleware"
    "net/http"
    "{{.BackendImport}}"
)

func InitRoutes() *{{.BackendPkg}}.Mux {
    r := {{.BackendInit}}
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Get("/", handlers.HomeHandler)
    r.Get("/about", handlers.AboutHandler)
    r.Get("/signin", handlers.SigninHandler)
    r.Get("/signup", handlers.SignupHandler)

    return r
}
`
		},
		"app/handlers/handlers.go": func(data TemplateData) string {
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

func AboutHandler(w http.ResponseWriter, r *http.Request) {
    if err := views.ExecuteTemplate(w, "about.templ", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
`
		},
		"app/db/db.go": func(data TemplateData) string {
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
		config := cfg.LoadFromEnv()
		db, err := db.New(config)
		if err != nil {
			log.Fatal(err)
		}
		Query = bun.NewDB(db, sqlitedialect.New())
		if config.IsDevelopment() {
			Query.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		}
}
`
		},
		"app/middleware/auth.go": func(data TemplateData) string {
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
		"app/middleware/cors.go": func(data TemplateData) string {
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
		"app/middleware/logger.go": func(data TemplateData) string {
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
		"app/middleware/notifier.go": func(data TemplateData) string {
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

        notifyChannel <- "System started"

        for {
            select {
            case <-r.Context().Done():
                return
            case msg := <-messageChan:
                fmt.Fprintf(w, "data: %s\n\n", msg)
                flusher.Flush()
            }
        }
    })
}

// Function to notify all connected SSE clients
func notifyClients(message string) {
    for client := range clients {
        client <- message
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
		"app/middleware/recoverer.go": func(data TemplateData) string {
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
		"app/middleware/rateLimiter.go": func(data TemplateData) string {
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
		"app/middleware/requestId.go": func(data TemplateData) string {
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
		"app/views/components/head.templ": func(data TemplateData) string {
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
		"app/views/layouts/base.templ": func(data TemplateData) string {
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
		"app/views/layouts/app.templ": func(data TemplateData) string {
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
		"app/views/components/header/header.templ": func(data TemplateData) string {
			return `package components

templ Header(){
	<header>
    	<h1>Welcome to {{.AppName}}</h1>
    </header>
}
`
		},
		"app/views/components/footer/footer.templ": func(data TemplateData) string {
			return `package components

templ Footer(){
	<footer>
   		<p>© {{.CurrentYear}} {{.AppName}}</p>
    </footer>
}
`
		},
		"app/views/pages/home.templ": func(data TemplateData) string {
			return `package pages

templ Home(){
	<h2>Home Page</h2>
	<p>This is the home page.</p>
}
`
		},
		"app/views/pages/about.templ": func(data TemplateData) string {
			return `package pages

templ About(){
	<h2>About Page</h2>
	<p>This is the about page.</p>
}
`
		},
		"app/views/components/navigation/sidebar.templ": func(data TemplateData) string {
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
		"public/public.go": func(data TemplateData) string {
			return `package public

import "embed"

//go:embed assets
var AssetsFS embed.FS
`
		},
		"go.mod": func(data TemplateData) string {
			return `module {{.AppName}}

go 1.22.4

require {{.VersionedBackendImport}}
`
		},
		"go.sum": func(data TemplateData) string {
			return ``
		},
		"README.md": func(data TemplateData) string {
			return `# {{ .AppName }}

A brief description of what your project does.

## Features

- Feature 1
- Feature 2
- Feature 3

## Installation

To install and run this project, follow these steps:

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/yourproject.git
   cd yourproject
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Set up environment variables (if any):

   ```sh
   cp .env.example .env
   # Edit the .env file with your configuration
   ```

4. Run the application:

   ```sh
   go run main.go
   ```

## Usage

### Running the Project

To start the project, use:

```sh
gost r
```

### Project Structure

By default gost creates the following structure for you:

```
.
├── cmd             # Main applications of the project
├── app             # Private application and library code
├── pkg             # Public library code
├── web             # Web server-related files
│   ├── static      # Static files
│   └── templates   # HTML templates
├── go.mod          # Go module file
├── main.go         # Main entry point of the application
└── README.md       # This file
```

### Running Tests

To run tests, use:

```sh
go test ./...
```

## Configuration

List any configuration settings for your project:

- `DATABASE_URL`: The URL of your database.
- `PORT`: The port on which the server will run.

## Contributing

We welcome contributions! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch with your feature or bug fix.
3. Commit your changes.
4. Push the branch to your fork.
5. Create a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Acknowledgements

Thanks to the contributors and the open-source community for their valuable input and support.

            `
		},
		"MAKEFILE": func(data TemplateData) string {
			return ``
		},
		".gitignore": func(data TemplateData) string {
			return ``
		},
		".air.toml": func(data TemplateData) string {
			return `[build]
cmd = "go build -o ./tmp/main ."
bin = "tmp/main"
watch = ["."]
exclude_dir = ["tmp", "vendor"]
exclude_file = ["go.sum", "go.mod", ".gitignore", ".DS_Store", ".idea"]
delay = 200
`
		},
		".env": func(data TemplateData) string {
			return `
# Application environment
# PROD or DEV
GOST_ENV=DEV

# HTTP listen port of the application
PORT=:9630

# Database Config
DB_DRIVER=sqlite3
DB_USER=
DB_HOST=
DB_PASSWORD=
DB_NAME=data.db

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
		"app/views/views.go": func(data TemplateData) string {
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
		"app/views/errors/404.templ": func(data TemplateData) string {
			return `package errors

templ _404(){
	<div>404 Page Not Found</div>
}
`
		},
		"app/views/errors/500.templ": func(data TemplateData) string {
			return `package errors

templ _500(){
		<div>500 Internal Server Error</div>
}
`
		},
	}
}

func GenerateFiles(data TemplateData) error {
	for path, tmplFunc := range FilesMap() {
		content := tmplFunc(data)
		filePath := filepath.Join(data.AppName, path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		} else {
			if strings.HasSuffix(filePath, ".go") {
				if err = cleaner.SortImports(filePath); err != nil {
					return fmt.Errorf("The file was saved but failed to sort its imports %w", err)
				}
			}
		}
	}
	return nil
}
