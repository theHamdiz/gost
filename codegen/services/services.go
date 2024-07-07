package services

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenServicesPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenServicesPlugin) Init() error {
	// Initialize Files
	g.Files = map[string]func() string{
		// 		"app/services/logger.go": func() string {
		// 			return `package services

		// import (
		//     "log"
		//     "os"
		// )

		// var (
		//     InfoLogger  *log.Logger
		//     ErrorLogger *log.Logger
		// )

		// func InitLogger() {
		//     InfoLogger = log.New(os.Stdout, "ℹ INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		//     ErrorLogger = log.New(os.Stderr, "✗ ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		// }

		// func Info(message string) {
		//     InfoLogger.Println(message)
		// }

		// func Error(message string) {
		//     ErrorLogger.Println(message)
		// }
		// `
		// 		},
		// 		"app/services/rateLimiter.go": func() string {
		// 			return `package services

		// import (
		//     "golang.org/x/time/rate"
		// )

		// var limiter *rate.Limiter

		// func InitRateLimiter(limit rate.Limit, burst int) {
		//     limiter = rate.NewLimiter(limit, burst)
		// }

		// func Allow() bool {
		//     return limiter.Allow()
		// }
		// `
		// 		},
	}

	return nil
}

func (g *GenServicesPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenServicesPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenServicesPlugin) Name() string {
	return "GenServicesPlugin"
}

func (g *GenServicesPlugin) Version() string {
	return "1.0.0"
}

func (g *GenServicesPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenServicesPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenServicesPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenServicesPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenServicesPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/services"
}

func (g *GenServicesPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenServicesPlugin(data config.ProjectData) *GenServicesPlugin {
	return &GenServicesPlugin{
		Data: data,
	}
}
