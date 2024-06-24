package files

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
			},
			"MAKEFILE": func() string {
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
		},
	}
}
