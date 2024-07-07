package files

import (
	"strings"

	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenFilesPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenFilesPlugin) Init() error {
	g.Files = map[string]func() string{
		"cmd/server/main.go": func() string {
			return `package main

import (
	{{- if eq .BackendPkg "echo" }}
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	{{- end }}
    "{{.AppName}}/app/cfg"
    "{{.AppName}}/app/db"
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
    sig := <- sigs
    log.Printf("Received signal: %s", sig)
}

func main() {
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
		"cmd/worker/main.go": func() string {
			return `package worker

import (
	{{- if eq .BackendPkg "echo" }}
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	{{- end }}
    "{{.AppName}}/app/cfg"
    "{{.AppName}}/app/db"
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
    sig := <- sigs
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
		"go.mod": func() string {
			return `module {{.AppName}}

go 1.22.4

{{ if ne .VersionedBackendImport ""}}
require {{.VersionedBackendImport}}
{{ end }}
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
				"- Feature 4\n\n" +
				"## Installation\n\n" +
				"To install and run this project, follow these steps:\n\n" +
				"1. Clone the repository:\n\n" +
				"```sh\n" +
				"git clone https://github.com/yourusername/yourproject.git\n" +
				"cd yourproject\n" +
				"```\n\n" +
				"2. Install dependencies if not already installed:\n\n" +
				"```sh\n" +
				"go mod tidy\n" +
				"```\n\n" +
				"3. Set up environment variables (if any):\n\n" +
				"```sh\n" +
				"cp .gost.env .env\n" +
				"# Edit the .env file with your configuration\n" +
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
				`├── app
│   ├── api
│   │   └── v1
│   ├── assets
│   │   └── static
│   │       ├── css
│   │       ├── img
│   │       └── js
│   ├── cfg
│   │   └── cfg.go
│   ├── db
│   │   ├── migrations
│   │   │   ├── create_db_1719424521947950600.sql
│   │   │   └── create_db_1719424522725851300.sql
│   │   ├── data.db
│   │   └── db.go
│   ├── events
│   │   └── events.go
│   ├── handlers
│   │   ├── api
│   │   │   └── api.go
│   │   ├── backend
│   │   │   ├── about.go
│   │   │   ├── auth.go
│   │   │   ├── landing.go
│   │   │   └── views.go
│   │   └── frontend
│   │       ├── about.go
│   │       ├── auth.go
│   │       ├── landing.go
│   │       └── views.go
│   ├── middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── notifier.go
│   │   ├── rateLimiter.go
│   │   ├── recoverer.go
│   │   └── requestId.go
│   ├── router
│   │   └── router.go
│   ├── services
│   │   ├── db.go
│   │   ├── logger.go
│   │   └── rateLimiter.go
│   ├── types
│   │   ├── core
│   │   │   └── gost.go
│   │   └── models
│   └── web
│       ├── backend
│       │   ├── assets
│       │   │   ├── css
│       │   │   └── js
│       │   ├── components
│       │   │   └── index.js
│       │   ├── pages
│       │   │   └── index.js
│       │   ├── store
│       │   │   └── index.js
│       │   ├── README.md
│       │   ├── index.html
│       │   ├── package.json
│       │   ├── robots.txt
│       │   ├── signin.html
│       │   ├── signup.html
│       │   └── vite.config.js
│       ├── components
│       │   ├── footer
│       │   │   └── footer.templ
│       │   ├── header
│       │   │   └── header.templ
│       │   ├── navigation
│       │   │   └── sidebar.templ
│       │   └── head.templ
│       ├── errors
│       │   ├── 404.templ
│       │   └── 500.templ
│       ├── frontend
│       │   ├── assets
│       │   │   ├── css
│       │   │   └── js
│       │   ├── components
│       │   │   └── index.js
│       │   ├── pages
│       │   │   ├── index.js
│       │   │   ├── signin.templ
│       │   │   └── signup.templ
│       │   ├── store
│       │   │   └── index.js
│       │   ├── README.md
│       │   ├── index.html
│       │   ├── package.json
│       │   ├── robots.txt
│       │   └── vite.config.js
│       ├── layouts
│       │   ├── app.templ
│       │   └── base.templ
│       ├── pages
│       │   ├── about.templ
│       │   └── home.templ
│       ├── public
│       │   └── index.html
│       ├── shared
│       ├── README.md
│       ├── embed.go
│       └── views.go
├── cmd
│   ├── app
│   │   └── main.go
│   └── scripts
├── log
├── plugins
│   ├── auth
│   ├── core
│   │   ├── config.go
│   │   └── core.go
│   └── db
│       ├── dialects
│       │   ├── db2.go
│       │   ├── dialects.go
│       │   ├── firebird.go
│       │   ├── mariadb.go
│       │   ├── mysql.go
│       │   ├── oracle.go
│       │   ├── postgresql.go
│       │   ├── sqlite.go
│       │   └── sqlserver.go
│       └── db.go
├── public
│   └── assets
├── storage
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── package-lock.json
└── package.json
` +
				"```\n\n" +
				"### Running Tests\n\n" +
				"To run tests, use:\n\n" +
				"```sh\n" +
				"gost t\n" +
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
		},
		"Makefile": func() string {
			return `# Makefile for {{.AppName}}

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME = {{.AppName}}
MAIN_FILE = cmd/app/main.go
BINARY_UNIX = $(BINARY_NAME)_unix
BINARY_WIN = $(BINARY_NAME)_windows

# Frontend parameters
FRONTEND_DIR = web/front
NPMCMD = npm
NPMINSTALL = $(NPMCMD) install
NPMRUNBUILD = $(NPMCMD) run build

# All target
all: test build run release frontend clean

# Test target
test:
	$(GOTEST) -v ./...

# Run target
run:
	air

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

.PHONY: all test build run release frontend clean
`
		},
		"Dockerfile": func() string {
			return `
`
		},
		".gitignore": func() string {
			return `.git/*
.gitignore
.idea/*
.vscode/*
.gost.env.dev
log/*`
		},
		".air.toml": func() string {
			return `[build]
cmd = "go build -o ./tmp/main ."
bin = "tmp/main"
watch = ["."]
exclude_dir = ["tmp", "vendor"]
exclude_file = ["go.sum", "go.mod", ".gitignore", ".DS_Store", ".idea", ".git", ".vscode", "node_modules", "storage", "log"]
delay = 200
`
		}}

	// build a .gost config file based on the format used by the user

	if strings.HasSuffix(g.Data.ConfigFile, ".env") {
		g.Files[".env"] = func() string {
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
		}
	} else if strings.HasSuffix(g.Data.ConfigFile, ".json") {
		g.Files["config.json"] = func() string {
			return `{
  ".gost.env": {
    "GOST_ENV": "DEV",
    "PORT": ":{{.Port}}",
    "DB_DRIVER": "{{.DbDriver}}",
    "DB_USER": "",
    "DB_HOST": "",
    "DB_PASSWORD": "",
    "DB_NAME": "db.db",
    "MIGRATIONS_DIR": "app/db/migrations",
    "GOST_SECRET": "{{.Fingerprint}}",
    "GOST_AUTH_REDIRECT_AFTER_LOGIN": "/profile",
    "GOST_AUTH_SESSION_EXPIRY_IN_HOURS": "72",
    "GOST_AUTH_SKIP_VERIFY": "true",
    "GOST_BACKEND": "{{.BackendPkg}}"
  }
}
`
		}
	} else if strings.HasSuffix(g.Data.ConfigFile, ".toml") {
		g.Files["config.toml"] = func() string {
			return `
[gost.env]
GOST_ENV = "DEV"
PORT = ":{{.Port}}"
DB_DRIVER = "{{.DbDriver}}"
DB_USER = ""
DB_HOST = ""
DB_PASSWORD = ""
DB_NAME = "db.db"
MIGRATIONS_DIR = "app/db/migrations"
GOST_SECRET = "{{.Fingerprint}}"
GOST_AUTH_REDIRECT_AFTER_LOGIN = "/profile"
GOST_AUTH_SESSION_EXPIRY_IN_HOURS = 72
GOST_AUTH_SKIP_VERIFY = true
GOST_BACKEND = "{{.BackendPkg}}"
`
		}
	} else {
		g.Files["config.yaml"] = func() string {
			return `
GOST_ENV: DEV
PORT: ":{{.Port}}"
DB_DRIVER: "{{.DbDriver}}"
DB_USER: ""
DB_HOST: ""
DB_PASSWORD: ""
DB_NAME: "db.db"
MIGRATIONS_DIR: "app/db/migrations"
GOST_SECRET: "{{.Fingerprint}}"
GOST_AUTH_REDIRECT_AFTER_LOGIN: "/profile"
GOST_AUTH_SESSION_EXPIRY_IN_HOURS: 72
GOST_AUTH_SKIP_VERIFY: true
GOST_BACKEND: "{{.BackendPkg}}"
`
		}
	}
	return nil
}

func (g *GenFilesPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenFilesPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenFilesPlugin) Name() string {
	return "GenFilesPlugin"
}

func (g *GenFilesPlugin) Version() string {
	return "1.0.0"
}

func (g *GenFilesPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenFilesPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenFilesPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenFilesPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenFilesPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/files"
}

func (g *GenFilesPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenFilesPlugin(data config.ProjectData) *GenFilesPlugin {
	return &GenFilesPlugin{
		Data: data,
	}
}
