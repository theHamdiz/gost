package db

import (
	"fmt"
	"time"

	"github.com/theHamdiz/gost/config"
	"github.com/theHamdiz/gost/gen/general"
	"github.com/theHamdiz/gost/seeder"
)

type Generator struct {
	Files map[string]func() string
}

func (g *Generator) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenerator() *Generator {
	now := time.Now().UTC().UnixNano()
	return &Generator{
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
	}
}
