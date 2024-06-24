package plugins

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
		},
	}
}
