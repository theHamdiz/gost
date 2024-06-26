package plugins

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenPluginsPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenPluginsPlugin) Init() error {
	// Initialize Files
	g.Files = map[string]func() string{
		"plugins/core/config.go": func() string {
			return `package core

import (
	"encoding/json"
	"fmt"
	"os"
)

type PluginConfig struct {
	Settings map[string]interface{}
}

func LoadConfig(filePath string) (*PluginConfig, error) {
	config := &PluginConfig{
		Settings: make(map[string]interface{}),
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config.Settings); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *PluginConfig) GetPluginConfig(pluginName string) map[string]interface{} {
	if cfg, ok := c.Settings[pluginName].(map[string]interface{}); ok {
		return cfg
	}
	return nil
}

type PluginManagerConfig struct {
	PluginsDir string
}
`
		},
		"plugins/db/dialects/dialects.go": func() string {
			return `package dialects

// Dialect interface for different SQL dialects
type Dialect interface {
	Select(columns ...string) string
	From(table string) string
	Where(condition string) string
	InsertInto(table string, columns ...string) string
	Values(values ...interface{}) string
	Update(table string) string
	Set(assignments ...string) string
	DeleteFrom(table string) string
	Limit(limit int) string
	Offset(offset int) string
	Join(table, condition string) string
	LeftJoin(table, condition string) string
	RightJoin(table, condition string) string
	OrderBy(columns ...string) string
	GroupBy(columns ...string) string
	Having(condition string) string
	Returning(columns ...string) string
	Placeholder() string
}
`
		},
		"plugins/db/dialects/sqlite.go": func() string {
			return `package dialects
			type SQLiteDialect struct{}

func (d *SQLiteDialect) Select(columns ...string) string {
	return fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
}

func (d *SQLiteDialect) From(table string) string {
	return fmt.Sprintf("FROM %s ", table)
}

func (d *SQLiteDialect) Where(condition string) string {
	return fmt.Sprintf("WHERE %s ", condition)
}

func (d *SQLiteDialect) InsertInto(table string, columns ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) ", table, strings.Join(columns, ", "))
}

func (d *SQLiteDialect) Values(values ...interface{}) string {
	valStr := make([]string, len(values))
	for i, v := range values {
		valStr[i] = fmt.Sprintf("'%v'", v)
	}
	return fmt.Sprintf("VALUES (%s) ", strings.Join(valStr, ", "))
}

func (d *SQLiteDialect) Update(table string) string {
	return fmt.Sprintf("UPDATE %s ", table)
}

func (d *SQLiteDialect) Set(assignments ...string) string {
	return fmt.Sprintf("SET %s ", strings.Join(assignments, ", "))
}

func (d *SQLiteDialect) DeleteFrom(table string) string {
	return fmt.Sprintf("DELETE FROM %s ", table)
}

func (d *SQLiteDialect) Limit(limit int) string {
	return fmt.Sprintf("LIMIT %d ", limit)
}

func (d *SQLiteDialect) Offset(offset int) string {
	return fmt.Sprintf("OFFSET %d ", offset)
}

func (d *SQLiteDialect) Join(table, condition string) string {
	return fmt.Sprintf("JOIN %s ON %s ", table, condition)
}

func (d *SQLiteDialect) LeftJoin(table, condition string) string {
	return fmt.Sprintf("LEFT JOIN %s ON %s ", table, condition)
}

func (d *SQLiteDialect) RightJoin(table, condition string) string {
	return fmt.Sprintf("RIGHT JOIN %s ON %s ", table, condition)
}

func (d *SQLiteDialect) OrderBy(columns ...string) string {
	return fmt.Sprintf("ORDER BY %s ", strings.Join(columns, ", "))
}

func (d *SQLiteDialect) GroupBy(columns ...string) string {
	return fmt.Sprintf("GROUP BY %s ", strings.Join(columns, ", "))
}

func (d *SQLiteDialect) Having(condition string) string {
	return fmt.Sprintf("HAVING %s ", condition)
}

func (d *SQLiteDialect) Returning(columns ...string) string {
	return "" // SQLite does not support RETURNING directly
}

func (d *SQLiteDialect) Placeholder() string {
	return "?"
}

			`
		},
		"plugins/db/dialects/postgresql.go": func() string {
			return `package dialects

import (
	"fmt"
	"strings"
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

func (d *PostgreSQLDialect) Limit(limit int) string {
	return fmt.Sprintf("LIMIT %d ", limit)
}

func (d *PostgreSQLDialect) Offset(offset int) string {
	return fmt.Sprintf("OFFSET %d ", offset)
}

func (d *PostgreSQLDialect) Join(table, condition string) string {
	return fmt.Sprintf("JOIN %s ON %s ", table, condition)
}

func (d *PostgreSQLDialect) LeftJoin(table, condition string) string {
	return fmt.Sprintf("LEFT JOIN %s ON %s ", table, condition)
}

func (d *PostgreSQLDialect) RightJoin(table, condition string) string {
	return fmt.Sprintf("RIGHT JOIN %s ON %s ", table, condition)
}

func (d *PostgreSQLDialect) OrderBy(columns ...string) string {
	return fmt.Sprintf("ORDER BY %s ", strings.Join(columns, ", "))
}

func (d *PostgreSQLDialect) GroupBy(columns ...string) string {
	return fmt.Sprintf("GROUP BY %s ", strings.Join(columns, ", "))
}

func (d *PostgreSQLDialect) Having(condition string) string {
	return fmt.Sprintf("HAVING %s ", condition)
}

func (d *PostgreSQLDialect) Returning(columns ...string) string {
	return fmt.Sprintf("RETURNING %s ", strings.Join(columns, ", "))
}

func (d *PostgreSQLDialect) Placeholder() string {
	return "$"
}

`
		},
		"plugins/db/dialects/mysql.go": func() string {
			return `package dialects

import (
	"fmt"
	"strings"
)

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

func (d *MySQLDialect) Limit(limit int) string {
	return fmt.Sprintf("LIMIT %d ", limit)
}

func (d *MySQLDialect) Offset(offset int) string {
	return fmt.Sprintf("OFFSET %d ", offset)
}

func (d *MySQLDialect) Join(table, condition string) string {
	return fmt.Sprintf("JOIN %s ON %s ", table, condition)
}

func (d *MySQLDialect) LeftJoin(table, condition string) string {
	return fmt.Sprintf("LEFT JOIN %s ON %s ", table, condition)
}

func (d *MySQLDialect) RightJoin(table, condition string) string {
	return fmt.Sprintf("RIGHT JOIN %s ON %s ", table, condition)
}

func (d *MySQLDialect) OrderBy(columns ...string) string {
	return fmt.Sprintf("ORDER BY %s ", strings.Join(columns, ", "))
}

func (d *MySQLDialect) GroupBy(columns ...string) string {
	return fmt.Sprintf("GROUP BY %s ", strings.Join(columns, ", "))
}

func (d *MySQLDialect) Having(condition string) string {
	return fmt.Sprintf("HAVING %s ", condition)
}

func (d *MySQLDialect) Returning(columns ...string) string {
	return "" // MySQL does not support RETURNING directly
}

func (d *MySQLDialect) Placeholder() string {
	return "?"
}

`
		},
		"plugins/db/dialects/oracle.go": func() string {
			return `package oracle
// OracleDialect is an implementation of Dialect for Oracle
type OracleDialect struct{}

func (d *OracleDialect) Select(columns ...string) string {
	return fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
}

func (d *OracleDialect) From(table string) string {
	return fmt.Sprintf("FROM %s ", table)
}

func (d *OracleDialect) Where(condition string) string {
	return fmt.Sprintf("WHERE %s ", condition)
}

func (d *OracleDialect) InsertInto(table string, columns ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) ", table, strings.Join(columns, ", "))
}

func (d *OracleDialect) Values(values ...interface{}) string {
	valStr := make([]string, len(values))
	for i, v := range values {
		valStr[i] = fmt.Sprintf("'%v'", v)
	}
	return fmt.Sprintf("VALUES (%s) ", strings.Join(valStr, ", "))
}

func (d *OracleDialect) Update(table string) string {
	return fmt.Sprintf("UPDATE %s ", table)
}

func (d *OracleDialect) Set(assignments ...string) string {
	return fmt.Sprintf("SET %s ", strings.Join(assignments, ", "))
}

func (d *OracleDialect) DeleteFrom(table string) string {
	return fmt.Sprintf("DELETE FROM %s ", table)
}

func (d *OracleDialect) Limit(limit int) string {
	return fmt.Sprintf("FETCH FIRST %d ROWS ONLY ", limit)
}

func (d *OracleDialect) Offset(offset int) string {
	return fmt.Sprintf("OFFSET %d ROWS ", offset)
}

func (d *OracleDialect) Join(table, condition string) string {
	return fmt.Sprintf("JOIN %s ON %s ", table, condition)
}

func (d *OracleDialect) LeftJoin(table, condition string) string {
	return fmt.Sprintf("LEFT JOIN %s ON %s ", table, condition)
}

func (d *OracleDialect) RightJoin(table, condition string) string {
	return fmt.Sprintf("RIGHT JOIN %s ON %s ", table, condition)
}

func (d *OracleDialect) OrderBy(columns ...string) string {
	return fmt.Sprintf("ORDER BY %s ", strings.Join(columns, ", "))
}

func (d *OracleDialect) GroupBy(columns ...string) string {
	return fmt.Sprintf("GROUP BY %s ", strings.Join(columns, ", "))
}

func (d *OracleDialect) Having(condition string) string {
	return fmt.Sprintf("HAVING %s ", condition)
}

func (d *OracleDialect) Returning(columns ...string) string {
	return fmt.Sprintf("RETURNING %s INTO ", strings.Join(columns, ", "))
}

func (d *OracleDialect) Placeholder() string {
	return ":param"
}

`
		},
		"plugins/db/dialects/sqlserver.go": func() string {
			return `package dialects
// SQLServerDialect is an implementation of Dialect for SQL Server
type SQLServerDialect struct{}

func (d *SQLServerDialect) Select(columns ...string) string {
	return fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
}

func (d *SQLServerDialect) From(table string) string {
	return fmt.Sprintf("FROM %s ", table)
}

func (d *SQLServerDialect) Where(condition string) string {
	return fmt.Sprintf("WHERE %s ", condition)
}

func (d *SQLServerDialect) InsertInto(table string, columns ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) ", table, strings.Join(columns, ", "))
}

func (d *SQLServerDialect) Values(values ...interface{}) string {
	valStr := make([]string, len(values))
	for i, v := range values {
		valStr[i] = fmt.Sprintf("'%v'", v)
	}
	return fmt.Sprintf("VALUES (%s) ", strings.Join(valStr, ", "))
}

func (d *SQLServerDialect) Update(table string) string {
	return fmt.Sprintf("UPDATE %s ", table)
}

func (d *SQLServerDialect) Set(assignments ...string) string {
	return fmt.Sprintf("SET %s ", strings.Join(assignments, ", "))
}

func (d *SQLServerDialect) DeleteFrom(table string) string {
	return fmt.Sprintf("DELETE FROM %s ", table)
}

func (d *SQLServerDialect) Limit(limit int) string {
	return fmt.Sprintf("TOP %d ", limit)
}

func (d *SQLServerDialect) Offset(offset int) string {
	return fmt.Sprintf("OFFSET %d ROWS ", offset)
}

func (d *SQLServerDialect) Join(table, condition string) string {
	return fmt.Sprintf("JOIN %s ON %s ", table, condition)
}

func (d *SQLServerDialect) LeftJoin(table, condition string) string {
	return fmt.Sprintf("LEFT JOIN %s ON %s ", table, condition)
}

func (d *SQLServerDialect) RightJoin(table, condition string) string {
	return fmt.Sprintf("RIGHT JOIN %s ON %s ", table, condition)
}

func (d *SQLServerDialect) OrderBy(columns ...string) string {
	return fmt.Sprintf("ORDER BY %s ", strings.Join(columns, ", "))
}

func (d *SQLServerDialect) GroupBy(columns ...string) string {
	return fmt.Sprintf("GROUP BY %s ", strings.Join(columns, ", "))
}

func (d *SQLServerDialect) Having(condition string) string {
	return fmt.Sprintf("HAVING %s ", condition)
}

func (d *SQLServerDialect) Returning(columns ...string) string {
	return "" // SQL Server does not support RETURNING directly
}

func (d *SQLServerDialect) Placeholder() string {
	return "@param"
}

`
		},
		"plugins/db/dialects/mariadb.go": func() string {
			return `package dialects
// MariaDBDialect is an implementation of Dialect for MariaDB
type MariaDBDialect struct{}

func (d *MariaDBDialect) Select(columns ...string) string {
	return fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
}

func (d *MariaDBDialect) From(table string) string {
	return fmt.Sprintf("FROM %s ", table)
}

func (d *MariaDBDialect) Where(condition string) string {
	return fmt.Sprintf("WHERE %s ", condition)
}

func (d *MariaDBDialect) InsertInto(table string, columns ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) ", table, strings.Join(columns, ", "))
}

func (d *MariaDBDialect) Values(values ...interface{}) string {
	valStr := make([]string, len(values))
	for i, v := range values {
		valStr[i] = fmt.Sprintf("'%v'", v)
	}
	return fmt.Sprintf("VALUES (%s) ", strings.Join(valStr, ", "))
}

func (d *MariaDBDialect) Update(table string) string {
	return fmt.Sprintf("UPDATE %s ", table)
}

func (d *MariaDBDialect) Set(assignments ...string) string {
	return fmt.Sprintf("SET %s ", strings.Join(assignments, ", "))
}

func (d *MariaDBDialect) DeleteFrom(table string) string {
	return fmt.Sprintf("DELETE FROM %s ", table)
}

func (d *MariaDBDialect) Limit(limit int) string {
	return fmt.Sprintf("LIMIT %d ", limit)
}

func (d *MariaDBDialect) Offset(offset int) string {
	return fmt.Sprintf("OFFSET %d ", offset)
}

func (d *MariaDBDialect) Join(table, condition string) string {
	return fmt.Sprintf("JOIN %s ON %s ", table, condition)
}

func (d *MariaDBDialect) LeftJoin(table, condition string) string {
	return fmt.Sprintf("LEFT JOIN %s ON %s ", table, condition)
}

func (d *MariaDBDialect) RightJoin(table, condition string) string {
	return fmt.Sprintf("RIGHT JOIN %s ON %s ", table, condition)
}

func (d *MariaDBDialect) OrderBy(columns ...string) string {
	return fmt.Sprintf("ORDER BY %s ", strings.Join(columns, ", "))
}

func (d *MariaDBDialect) GroupBy(columns ...string) string {
	return fmt.Sprintf("GROUP BY %s ", strings.Join(columns, ", "))
}

func (d *MariaDBDialect) Having(condition string) string {
	return fmt.Sprintf("HAVING %s ", condition)
}

func (d *MariaDBDialect) Returning(columns ...string) string {
	return "" // MariaDB does not support RETURNING directly
}

func (d *MariaDBDialect) Placeholder() string {
	return "?"
}

`
		},
		"plugins/db/dialects/firebird.go": func() string {
			return `package dialects
// FirebirdDialect is an implementation of Dialect for Firebird SQL
type FirebirdDialect struct{}

func (d *FirebirdDialect) Select(columns ...string) string {
	return fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
}

func (d *FirebirdDialect) From(table string) string {
	return fmt.Sprintf("FROM %s ", table)
}

func (d *FirebirdDialect) Where(condition string) string {
	return fmt.Sprintf("WHERE %s ", condition)
}

func (d *FirebirdDialect) InsertInto(table string, columns ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) ", table, strings.Join(columns, ", "))
}

func (d *FirebirdDialect) Values(values ...interface{}) string {
	valStr := make([]string, len(values))
	for i, v := range values {
		valStr[i] = fmt.Sprintf("'%v'", v)
	}
	return fmt.Sprintf("VALUES (%s) ", strings.Join(valStr, ", "))
}

func (d *FirebirdDialect) Update(table string) string {
	return fmt.Sprintf("UPDATE %s ", table)
}

func (d *FirebirdDialect) Set(assignments ...string) string {
	return fmt.Sprintf("SET %s ", strings.Join(assignments, ", "))
}

func (d *FirebirdDialect) DeleteFrom(table string) string {
	return fmt.Sprintf("DELETE FROM %s ", table)
}

func (d *FirebirdDialect) Limit(limit int) string {
	return fmt.Sprintf("FIRST %d ", limit)
}

func (d *FirebirdDialect) Offset(offset int) string {
	return fmt.Sprintf("SKIP %d ", offset)
}

func (d *FirebirdDialect) Join(table, condition string) string {
	return fmt.Sprintf("JOIN %s ON %s ", table, condition)
}

func (d *FirebirdDialect) LeftJoin(table, condition string) string {
	return fmt.Sprintf("LEFT JOIN %s ON %s ", table, condition)
}

func (d *FirebirdDialect) RightJoin(table, condition string) string {
	return fmt.Sprintf("RIGHT JOIN %s ON %s ", table, condition)
}

func (d *FirebirdDialect) OrderBy(columns ...string) string {
	return fmt.Sprintf("ORDER BY %s ", strings.Join(columns, ", "))
}

func (d *FirebirdDialect) GroupBy(columns ...string) string {
	return fmt.Sprintf("GROUP BY %s ", strings.Join(columns, ", "))
}

func (d *FirebirdDialect) Having(condition string) string {
	return fmt.Sprintf("HAVING %s ", condition)
}

func (d *FirebirdDialect) Returning(columns ...string) string {
	return fmt.Sprintf("RETURNING %s ", strings.Join(columns, ", "))
}

func (d *FirebirdDialect) Placeholder() string {
	return "?"
}

`
		},
		"plugins/db/dialects/db2.go": func() string {
			return `package dialects
// DB2Dialect is an implementation of Dialect for IBM Db2
type DB2Dialect struct{}

func (d *DB2Dialect) Select(columns ...string) string {
	return fmt.Sprintf("SELECT %s ", strings.Join(columns, ", "))
}

func (d *DB2Dialect) From(table string) string {
	return fmt.Sprintf("FROM %s ", table)
}

func (d *DB2Dialect) Where(condition string) string {
	return fmt.Sprintf("WHERE %s ", condition)
}

func (d *DB2Dialect) InsertInto(table string, columns ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) ", table, strings.Join(columns, ", "))
}

func (d *DB2Dialect) Values(values ...interface{}) string {
	valStr := make([]string, len(values))
	for i, v := range values {
		valStr[i] = fmt.Sprintf("'%v'", v)
	}
	return fmt.Sprintf("VALUES (%s) ", strings.Join(valStr, ", "))
}

func (d *DB2Dialect) Update(table string) string {
	return fmt.Sprintf("UPDATE %s ", table)
}

func (d *DB2Dialect) Set(assignments ...string) string {
	return fmt.Sprintf("SET %s ", strings.Join(assignments, ", "))
}

func (d *DB2Dialect) DeleteFrom(table string) string {
	return fmt.Sprintf("DELETE FROM %s ", table)
}

func (d *DB2Dialect) Limit(limit int) string {
	return fmt.Sprintf("FETCH FIRST %d ROWS ONLY ", limit)
}

func (d *DB2Dialect) Offset(offset int) string {
	return fmt.Sprintf("OFFSET %d ROWS ", offset)
}

func (d *DB2Dialect) Join(table, condition string) string {
	return fmt.Sprintf("JOIN %s ON %s ", table, condition)
}

func (d *DB2Dialect) LeftJoin(table, condition string) string {
	return fmt.Sprintf("LEFT JOIN %s ON %s ", table, condition)
}

func (d *DB2Dialect) RightJoin(table, condition string) string {
	return fmt.Sprintf("RIGHT JOIN %s ON %s ", table, condition)
}

func (d *DB2Dialect) OrderBy(columns ...string) string {
	return fmt.Sprintf("ORDER BY %s ", strings.Join(columns, ", "))
}

func (d *DB2Dialect) GroupBy(columns ...string) string {
	return fmt.Sprintf("GROUP BY %s ", strings.Join(columns, ", "))
}

func (d *DB2Dialect) Having(condition string) string {
	return fmt.Sprintf("HAVING %s ", condition)
}

func (d *DB2Dialect) Returning(columns ...string) string {
	return fmt.Sprintf("RETURNING %s ", strings.Join(columns, ", "))
}

func (d *DB2Dialect) Placeholder() string {
	return "@param"
}

			`
		},
		"plugins/core/core.go": func() string {
			return `package core

import (
	"errors"
	"fmt"
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
	config          *PluginManagerConfig
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

func NewPluginManager(config *PluginManagerConfig, pluginDirectory string) *PluginManager {
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
		"plugins/db/db.go": func() string {
			return `package plugins

import (
	"database/sql"
	"fmt"
	"goat/plugins/db/dialects"
	"reflect"
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

// Model represents a generic model that can provide its table name and column mappings.
type Model interface {
	TableName() string
	ColumnMappings() map[string]string // field name to column name
}

// DbPlugin interface represents the DB plugin
type DbPlugin interface {
	Plugin
	NewDbBuilder(dialect dialects.Dialect) *DbBuilder
}

// DB manages the database connection and provides methods to start new queries.
type DB struct {
	dialect dialects.Dialect
	db      *sql.DB
}

// Init initializes the DB plugin
func (o *DB) Init() error {
	// Initialization logic if needed
	return nil
}

// Execute executes the DB plugin (example placeholder implementation)
func (o *DB) Execute() error {
	// Execution logic if needed
	return nil
}

// Shutdown cleans up the DB plugin
func (o *DB) Shutdown() error {
	// Shutdown logic if needed
	return nil
}

// Name returns the name of the DB plugin
func (o *DB) Name() string {
	return "Natural Orm For Gost"
}

// Version returns the version of the DB plugin
func (o *DB) Version() string {
	return "1.0.0"
}

// Dependencies returns the dependencies of the DB plugin
func (o *DB) Dependencies() []string {
	return []string{}
}

// AuthorName returns the author's name
func (o *DB) AuthorName() string {
	return "Ahmad Hamdi"
}

// AuthorEmail returns the author's email
func (o *DB) AuthorEmail() string {
	return "contact@hamdiz.me"
}

// Website returns the website of the plugin
func (o *DB) Website() string {
	return "https://hamdiz.me"
}

// GitHub returns the GitHub URL of the plugin
func (o *DB) GitHub() string {
	return "https://github.com/theHamdiz/gost/plugins/orm"
}

// NewDbBuilder creates a new DbBuilder for a given dialect
func (o *DB) NewDbBuilder(dialect dialects.Dialect) *DbBuilder {
	return NewDbBuilder(dialect)
}

func (o *DB) Select(columns ...string) *DbBuilder {
	builder := &DbBuilder{
		dialect: o.dialect,
		db:      o.db,
	}
	return builder.Select(columns...)
}

// NewDB creates a new DB instance
func NewDB(dialect dialects.Dialect, db *sql.DB) *DB {
	return &DB{
		dialect: dialect,
		db:      db,
	}
}

// ----------------------- //
// Next:DbBuilder Type
// ---------------------- //

// DbBuilder provides a fluent API for building queries through the builder pattern.
type DbBuilder struct {
	dialect dialects.Dialect
	query   strings.Builder
	args    []interface{}
	state   State
	db      *sql.DB
}

/*
Limit adds a LIMIT clause to the query.
This method allows specifying the maximum number of rows to return from the query.

Parameters:

	limit (int): The maximum number of rows to return.

Returns:

	*DbBuilder: The current DbBuilder instance with the LIMIT clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("id", "name").
	              From("users").
	              Where("active = ?", true).
	              OrderBy("name ASC").
	              Limit(5)
*/
func (b *DbBuilder) Limit(limit int) *DbBuilder {
	b.query.WriteString(b.dialect.Limit(limit))
	return b
}

/*
Offset adds an OFFSET clause to the query.
This method allows specifying the number of rows to skip before starting to return rows from the query.

Parameters:

	offset (int): The number of rows to skip.

Returns:

	*DbBuilder: The current DbBuilder instance with the OFFSET clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("id", "name").
	              From("users").
	              Where("active = ?", true).
	              OrderBy("name ASC").
	              Offset(10).
	              Limit(5)
*/
func (b *DbBuilder) Offset(offset int) *DbBuilder {
	b.query.WriteString(b.dialect.Offset(offset))
	return b
}

/*
Join adds a JOIN clause to the query.
This method allows specifying a table and condition for a JOIN operation.

Parameters:

	table (string): The name of the table to join.
	condition (string): The condition for the JOIN.

Returns:

	*DbBuilder: The current DbBuilder instance with the JOIN clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("a.id", "b.name").
	              From("table_a a").
	              Join("table_b b", "a.id = b.a_id").
	              Where("a.active = ?", true)
*/
func (b *DbBuilder) Join(table, condition string) *DbBuilder {
	b.query.WriteString(b.dialect.Join(table, condition))
	return b
}

/*
LeftJoin adds a LEFT JOIN clause to the query.
This method allows specifying a table and condition for a LEFT JOIN operation.

Parameters:

	table (string): The name of the table to join.
	condition (string): The condition for the LEFT JOIN.

Returns:

	*DbBuilder: The current DbBuilder instance with the LEFT JOIN clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("a.id", "b.name").
	              From("table_a a").
	              LeftJoin("table_b b", "a.id = b.a_id").
	              Where("a.active = ?", true)
*/
func (b *DbBuilder) LeftJoin(table, condition string) *DbBuilder {
	b.query.WriteString(b.dialect.LeftJoin(table, condition))
	return b
}

/*
RightJoin adds a RIGHT JOIN clause to the query.
This method allows specifying a table and condition for a RIGHT JOIN operation.

Parameters:

	table (string): The name of the table to join.
	condition (string): The condition for the RIGHT JOIN.

Returns:

	*DbBuilder: The current DbBuilder instance with the RIGHT JOIN clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("a.id", "b.name").
	              From("table_a a").
	              RightJoin("table_b b", "a.id = b.a_id").
	              Where("a.active = ?", true)
*/
func (b *DbBuilder) RightJoin(table, condition string) *DbBuilder {
	b.query.WriteString(b.dialect.RightJoin(table, condition))
	return b
}

/*
OrderBy adds an ORDER BY clause to the query.
This method allows specifying columns by which the result set should be ordered.

Parameters:

	columns (...string): The columns by which to order the result set.

Returns:

	*DbBuilder: The current DbBuilder instance with the ORDER BY clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("id", "name").
	              From("users").
	              Where("active = ?", true).
	              OrderBy("name ASC", "id DESC")
*/
func (b *DbBuilder) OrderBy(columns ...string) *DbBuilder {
	b.query.WriteString(b.dialect.OrderBy(columns...))
	return b
}

/*
GroupBy adds a GROUP BY clause to the query.
This method allows specifying columns by which the result set should be grouped.

Parameters:

	columns (...string): The columns by which to group the result set.

Returns:

	*DbBuilder: The current DbBuilder instance with the GROUP BY clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("department", "COUNT(*) as num_employees").
	              From("employees").
	              GroupBy("department")
*/
func (b *DbBuilder) GroupBy(columns ...string) *DbBuilder {
	b.query.WriteString(b.dialect.GroupBy(columns...))
	return b
}

/*
Having adds a HAVING clause to the query.
This method allows specifying a condition for groups, typically used in conjunction with a GROUP BY clause.

Parameters:

	condition (string): The condition for the HAVING clause.

Returns:

	*DbBuilder: The current DbBuilder instance with the HAVING clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("department", "COUNT(*) as num_employees").
	              From("employees").
	              GroupBy("department").
	              Having("num_employees > 10")
*/
func (b *DbBuilder) Having(condition string) *DbBuilder {
	b.query.WriteString(b.dialect.Having(condition))
	return b
}

/*
Returning adds a RETURNING clause to the query.
This method allows specifying which columns should be returned after an INSERT, UPDATE, or DELETE operation.

Parameters:

	columns (...string): The columns to be returned by the query.

Returns:

	*DbBuilder: The current DbBuilder instance with the RETURNING clause added.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              InsertInto("users", "name", "active").
	              Values("John Doe", true).
	              Returning("id", "name")
*/
func (b *DbBuilder) Returning(columns ...string) *DbBuilder {
	b.query.WriteString(b.dialect.Returning(columns...))
	return b
}

/*
Scan method to map query results to generic structs using reflection.
This method executes the built query and maps the results to the provided destination slice of structs.

Parameters:

	dest (interface{}): A pointer to a slice of structs where the query results will be stored.

Returns:

	error: An error object if the query execution or result scanning fails, otherwise nil.

Example usage:

	var users []User
	err := db.NewDbBuilder(dialect).
	           Select("id", "name").
	           From("users").
	           Where("active = ?", true).
	           Scan(&users)
	if err != nil {
	    log.Fatalf("Failed to scan query results: %v", err)
	}

	for _, user := range users {
	    fmt.Printf("User: %+v\n", user)
	}
*/
func (b *DbBuilder) Scan(dest interface{}) error {
	rows, err := b.db.Query(b.query.String(), b.args...)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}(rows)

	return scanRows(rows, dest)
}

/*
Exec method for executing queries without returning rows.
This method executes the built query against the database without expecting any rows in return.

Returns:

	error: An error object if the execution fails, otherwise nil.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              InsertInto("users", "name", "active").
	              Values("John Doe", true)
	err := builder.Exec()
	if err != nil {
	    log.Fatalf("Failed to execute query: %v", err)
	}
*/
func (b *DbBuilder) Exec() error {
	_, err := b.db.Exec(b.query.String(), b.args...)
	return err
}

/*
Select adds a SELECT clause to the query.
This method sets the state to Selecting and ensures that the SELECT clause is
only added at the beginning of the query construction.

Parameters:

	columns (...string): The columns to be selected in the query.

Returns:

	*DbBuilder: The current DbBuilder instance with the SELECT clause added.

Panics:

	Will panic if called after the initial state, as SELECT can only be called at the beginning of the query.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("id", "name").
	              From("users").
	              Where("active = ?", true)
*/
func (b *DbBuilder) Select(columns ...string) *DbBuilder {
	if b.state != Initial {
		panic("Select can only be called at the beginning of the query")
	}
	b.query.WriteString(b.dialect.Select(columns...))
	b.state = Selecting
	return b
}

/*
From adds a FROM clause to the query.
This method sets the state to Froming and ensures that the FROM clause is
only added after a SELECT clause.

Parameters:

	table (string): The name of the table from which to select data.

Returns:

	*DbBuilder: The current DbBuilder instance with the FROM clause added.

Panics:

	Will panic if called before a SELECT clause, as FROM must be called after SELECT.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("id", "name").
	              From("users").
	              Where("active = ?", true)
*/
func (b *DbBuilder) From(table string) *DbBuilder {
	if b.state != Selecting {
		panic("From must be called after Select")
	}
	b.query.WriteString(b.dialect.From(table))
	b.state = Froming
	return b
}

/*
Where adds a WHERE clause to the query.
This method sets the state to Whereing and ensures that the WHERE clause is
only added after a FROM or UPDATE clause.

Parameters:

	condition (string): The condition for the WHERE clause.

Returns:

	*DbBuilder: The current DbBuilder instance with the WHERE clause added.

Panics:

	Will panic if called before a FROM or UPDATE clause, as WHERE must be called after FROM or UPDATE.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("id", "name").
	              From("users").
	              Where("active = ?", true)
*/
func (b *DbBuilder) Where(condition string) *DbBuilder {
	if b.state != Froming && b.state != Updating {
		panic("Where must be called after From or Update")
	}
	b.query.WriteString(b.dialect.Where(condition))
	b.state = Whereing
	return b
}

/*
Insert inserts a new record into the table represented by the given model.
This method uses reflection to dynamically determine the columns and values from the model's fields.

Parameters:

	model (Model): The model representing the table and data to be inserted.

Returns:

	*OrmBuilder: The current OrmBuilder instance with the INSERT INTO clause and VALUES clause added.

Panics:

	Will panic if called after the initial state, as Insert can only be called at the beginning of the query.

Example usage:

	user := User{
	    Name:   "John Doe",
	    Active: true,
	}
	err := db.NewOrmBuilder(dialect).
	              Insert(user).
	              Exec()
	if err != nil {
	    log.Fatalf("Failed to insert user: %v", err)
	}
*/
func (b *DbBuilder) Insert(model Model) *DbBuilder {
	if b.state != Initial {
		panic("Insert can only be called at the beginning of the query")
	}

	// Get the table name and column mappings from the model
	table := model.TableName()
	mappings := model.ColumnMappings()

	// Extract column names and values from the model using reflection
	columns := make([]string, 0, len(mappings))
	values := make([]interface{}, 0, len(mappings))
	modelValue := reflect.ValueOf(model).Elem()

	for field, column := range mappings {
		columns = append(columns, column)
		values = append(values, modelValue.FieldByName(field).Interface())
	}

	// Build the INSERT INTO clause
	b.query.WriteString(b.dialect.InsertInto(table, columns...))
	b.state = Inserting

	// Build the VALUES clause
	b.query.WriteString(b.dialect.Values(values...))
	b.args = append(b.args, values...)
	b.state = Valuing

	return b
}

/*
InsertInto adds an INSERT INTO clause to the query.
This method sets the state to Inserting and ensures that the INSERT INTO clause is
only added at the beginning of the query construction.

Parameters:

	table (string): The name of the table into which rows should be inserted.
	columns (...string): The columns into which values should be inserted.

Returns:

	*DbBuilder: The current DbBuilder instance with the INSERT INTO clause added.

Panics:

	Will panic if called after the initial state, as INSERT INTO can only be called at the beginning of the query.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              InsertInto("users", "name", "active").
	              Values("John Doe", true)
*/
func (b *DbBuilder) InsertInto(table string, columns ...string) *DbBuilder {
	if b.state != Initial {
		panic("InsertInto can only be called at the beginning of the query")
	}
	b.query.WriteString(b.dialect.InsertInto(table, columns...))
	b.state = Inserting
	return b
}

/*
Values adds a VALUES clause to the query.
This method sets the state to Valuing and ensures that the VALUES clause is
only added after an INSERT INTO clause.

Parameters:

	values (...interface{}): The values to be inserted into the table.

Returns:

	*DbBuilder: The current DbBuilder instance with the VALUES clause added.

Panics:

	Will panic if called before an INSERT INTO clause, as VALUES must be called after INSERT INTO.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              InsertInto("users", "name", "active").
	              Values("John Doe", true)
*/
func (b *DbBuilder) Values(values ...interface{}) *DbBuilder {
	if b.state != Inserting {
		panic("Values must be called after InsertInto")
	}
	b.query.WriteString(b.dialect.Values(values...))
	b.state = Valuing
	return b
}

/*
Update adds an UPDATE clause to the query.
This method sets the state to Updating and ensures that the UPDATE clause is
only added at the beginning of the query construction.

Parameters:

	table (string): The name of the table to be updated.

Returns:

	*DbBuilder: The current DbBuilder instance with the UPDATE clause added.

Panics:

	Will panic if called after the initial state, as UPDATE can only be called at the beginning of the query.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Update("users").
	              Set("name = ?", "active = ?").
	              Where("id = ?", 1)
*/
func (b *DbBuilder) Update(table string) *DbBuilder {
	if b.state != Initial {
		panic("Update can only be called at the beginning of the query")
	}
	b.query.WriteString(b.dialect.Update(table))
	b.state = Updating
	return b
}

/*
Set adds a SET clause to the query.
This method sets the state to Setting and ensures that the SET clause is
only added after an UPDATE clause.

Parameters:

	assignments (...string): The column assignments for the SET clause in the form of "column = value".

Returns:

	*DbBuilder: The current DbBuilder instance with the SET clause added.

Panics:

	Will panic if called before an UPDATE clause, as SET must be called after UPDATE.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Update("users").
	              Set("name = ?", "active = ?").
	              Where("id = ?", 1)
*/
func (b *DbBuilder) Set(assignments ...string) *DbBuilder {
	if b.state != Updating {
		panic("Set must be called after Update")
	}
	b.query.WriteString(b.dialect.Set(assignments...))
	b.state = Setting
	return b
}

/*
DeleteFrom adds a DELETE FROM clause to the query.
This method sets the state to Deleting and ensures that the DELETE FROM clause is
only added at the beginning of the query construction.

Parameters:

	table (string): The name of the table from which rows should be deleted.

Returns:

	*DbBuilder: The current DbBuilder instance with the DELETE FROM clause added.

Panics:

	Will panic if called after the initial state, as DELETE FROM can only be called at the beginning of the query.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              DeleteFrom("users").
	              Where("active = ?", false)
*/
func (b *DbBuilder) DeleteFrom(table string) *DbBuilder {
	if b.state != Initial {
		panic("DeleteFrom can only be called at the beginning of the query")
	}
	b.query.WriteString(b.dialect.DeleteFrom(table))
	b.state = Deleting
	return b
}

/*
Build returns the final SQL query as a string.
It concatenates all the parts of the query that have been built up so far.
This method can be called after the various query-building methods (e.g., Select, From, Where)
to get the complete SQL query as a string.

Example usage:

	builder := db.NewDbBuilder(dialect).
	              Select("id", "name").
	              From("users").
	              Where("active = ?", true).
	              OrderBy("name")
	query := builder.Build()
*/
func (b *DbBuilder) Build() string {
	return b.query.String()
}

/*
	scanRows -> maps the results of the SQL query to a destination slice of structs.

rows: The SQL rows returned from the query.
dest: A pointer to a slice of structs to which the results will be mapped.
Returns an error if there is any issue with the rows or reflection operations.
*/
func scanRows(rows *sql.Rows, dest interface{}) error {
	// Ensure dest is a pointer to a slice
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to a slice")
	}

	// Get the element type of the slice
	destValue = destValue.Elem()
	destType := destValue.Type().Elem()

	// Get the columns from the SQL rows
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	// Iterate over the rows
	for rows.Next() {
		// Create a new instance of the destination struct type
		model := reflect.New(destType).Elem()

		// Create a slice of field pointers to scan the row values into
		fieldPtrs := make([]interface{}, len(columns))
		for i, col := range columns {
			// Find the struct field that matches the column name
			field := model.FieldByNameFunc(func(name string) bool {
				field, _ := destType.FieldByName(name)
				return strings.EqualFold(field.Tag.Get("db"), col)
			})
			if field.IsValid() {
				// Use the address of the struct field
				fieldPtrs[i] = field.Addr().Interface()
			} else {
				// Use a placeholder if the struct does not have a matching field
				var placeholder interface{}
				fieldPtrs[i] = &placeholder
			}
		}

		// Scan the row values into the field pointers
		if err := rows.Scan(fieldPtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		// Append the populated struct to the destination slice
		destValue.Set(reflect.Append(destValue, model))
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error encountered during rows iteration: %w", err)
	}

	return nil
}

// NewDbBuilder creates a new DbBuilder for a given dialect!
func NewDbBuilder(dialect dialects.Dialect) *DbBuilder {
	return &DbBuilder{
		dialect: dialect,
		state:   Initial,
	}
}

`
		},
	}
	return nil
}

func (g *GenPluginsPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenPluginsPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenPluginsPlugin) Name() string {
	return "GenPluginsPlugin"
}

func (g *GenPluginsPlugin) Version() string {
	return "1.0.0"
}

func (g *GenPluginsPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenPluginsPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenPluginsPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenPluginsPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenPluginsPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/plugins"
}

func (g *GenPluginsPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenPluginsPlugin(data config.ProjectData) *GenPluginsPlugin {
	return &GenPluginsPlugin{
		Data: data,
	}
}
