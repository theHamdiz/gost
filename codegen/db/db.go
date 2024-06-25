package db

import (
	"fmt"
	"time"

	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
	"github.com/theHamdiz/gost/seeder"
)

type GenDbPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenDbPlugin) Init() error {
	// Any initialization logic for the plugin
	return nil
}

func (g *GenDbPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenDbPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenDbPlugin) Name() string {
	return "GenDbPlugin"
}

func (g *GenDbPlugin) Version() string {
	return "1.0.0"
}

func (g *GenDbPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenDbPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenDbPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenDbPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenDbPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/db"
}

func (g *GenDbPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenDbPlugin(data config.ProjectData) *GenDbPlugin {
	now := time.Now().UTC().UnixNano()
	return &GenDbPlugin{
		Files: map[string]func() string{
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
			fmt.Sprintf("app/db/migrations/create_db_%d.sql", now): func() string { return seeder.GetSeedingScript() },
		},
		Data: data,
	}
}
