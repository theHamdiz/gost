package plugins

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
			"plugins/core.go": func() string {
				return `package plugins

import (
	"errors"
	"fmt"
	"github.com/theHamdiz/gost/plugins/config"
)

type Plugin interface {
	Init() error
	Execute() error
	Shutdown() error
	Name() string
	Version() string
	Dependencies() []string
	AuthorName() string
	AuthorEmail() string
	Website() string
	GitHub() string
}

type GenericPluginManager interface {
	InitPlugins() error
	ExecutePlugins() error
	ShutdownPlugins() error
	RegisterPlugins(plugins []Plugin) error
	RegisterPlugin(plugin Plugin) error
}

type PluginMetadata struct {
	Name         string   ` + "`yaml:\"name\" json:\"name\"`" + `
	Version      string   ` + "`yaml:\"version\"  json:\"version\"`" + `
	Dependencies []string ` + "`yaml:\"dependencies\"  json:\"dependencies\"`" + `
	AuthorName   string   ` + "`yaml:\"author_name\"  json:\"author_name\"`" + `
	AuthorEmail  string   ` + "`yaml:\"author_email\"  json:\"author_email\"`" + `
	Website      string   ` + "`yaml:\"website\"  json:\"website\"`" + `
	GitHub       string   ` + "`yaml:\"github\"  json:\"github\"`" + `
}

type PluginManager struct {
	plugins         map[string]Plugin
	pluginOrder     []string
	config          *config.PluginManagerConfig
	pluginDirectory string
}

func (pm *PluginManager) RegisterPlugin(plugin Plugin) error {
	if plugin.Name() == "" {
		return errors.New("plugin name cannot be empty")
	}
	pm.plugins[plugin.Name()] = plugin
	pm.pluginOrder = append(pm.pluginOrder, plugin.Name())
	return nil
}

func (pm *PluginManager) RegisterPlugins(plugins []Plugin) error {
	for _, plugin := range plugins {
		err := pm.RegisterPlugin(plugin)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewPluginManager(config *config.PluginManagerConfig, pluginDirectory string) *PluginManager {
	return &PluginManager{
		plugins:         make(map[string]Plugin),
		config:          config,
		pluginDirectory: pluginDirectory,
	}
}

func (pm *PluginManager) resolveDependencies() ([]Plugin, error) {
	resolved := make(map[string]bool)
	var order []Plugin

	var resolve func(name string) error
	resolve = func(name string) error {
		if resolved[name] {
			return nil
		}

		plugin, exists := pm.plugins[name]
		if !exists {
			return fmt.Errorf("plugin %s not found", name)
		}

		for _, dep := range plugin.Dependencies() {
			if err := resolve(dep); err != nil {
				return err
			}
		}

		resolved[name] = true
		order = append(order, plugin)
		return nil
	}

	for name := range pm.plugins {
		if err := resolve(name); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (pm *PluginManager) InitPlugins() error {
	for _, name := range pm.pluginOrder {
		plugin := pm.plugins[name]
		if err := plugin.Init(); err != nil {
			return err
		}
	}
	return nil
}

func (pm *PluginManager) ExecutePlugins() error {
	for _, name := range pm.pluginOrder {
		plugin := pm.plugins[name]
		if err := plugin.Execute(); err != nil {
			return err
		}
	}
	return nil
}

func (pm *PluginManager) ShutdownPlugins() error {
	// Reverse the order for shutdown
	for i := len(pm.pluginOrder) - 1; i >= 0; i-- {
		name := pm.pluginOrder[i]
		plugin := pm.plugins[name]
		if err := plugin.Shutdown(); err != nil {
			return err
		}
	}
	return nil
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
