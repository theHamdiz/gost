package services

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
		},
	}
}
