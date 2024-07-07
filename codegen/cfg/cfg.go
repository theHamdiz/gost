package cfg

import (
	"bytes"
	"html/template"

	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenConfPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

const cfgTemplate = `package cfg

import (
    "os"
    "strings"
)

var(
	{{- if eq .PreferredConfigFormat ".env"}}
	ConfigFile string = ".gost.env"
	{{- else if eq .PreferredConfigFormat ".json"}}
	ConfigFile string = ".gost.json"
	{{- else if eq .PreferredConfigFormat ".toml"}}
	ConfigFile string = ".gost.toml"
	{{- else if eq .PreferredConfigFormat ".yaml"}}
	ConfigFile string = ".gost.yaml"
	{{- end }}
)

type Configurable interface {
	{{- if eq .PreferredConfigFormat ".env"}}
	SaveAsEnv(filePath string) error
	{{- else if eq .PreferredConfigFormat ".json"}}
	SaveAsJSON(filePath string) error
	{{- else if eq .PreferredConfigFormat ".toml"}}
	SaveAsTOML(filePath string) error
	{{- else if eq .PreferredConfigFormat ".yaml"}}
	SaveAsYAML(filePath string) error
	{{- end }}
}

type Config struct {
	AppName                         string
	DbDriver                        string
	DbHost                          string
	DbName                          string
	DbOrm                           string
	DbPassword                      string
	DbUri                           string
	DbUser                          string
	GostAuthRedirectAfterLogin      string
	GostAuthSessionExpiryInHours    string
	GostAuthSkipVerify              bool
	BackendPkg                      string
	GostEnv                         string
	GostSecret                      string
	MigrationsDir                   string
	Port                            string
	RedisDb                         string
	RedisPassword                   string
	RedisUri                        string
}

func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.GostEnv) == "dev" || strings.ToLower(c.GostEnv) == "development"
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

func LoadConfig() (*Config, error) {
	{{- if eq .ConfigFile ".env"}}
    return &Config{
        GostEnv:                      getEnv("GOST_ENV", "DEV"),
        Port:                         getEnv("PORT", ":9630"),
        DbDriver:                     getEnv("DB_DRIVER", "sqlite3"),
        DbUser:                       getEnv("DB_USER", ""),
        DbHost:                       getEnv("DB_HOST", ""),
        DbPassword:                   getEnv("DB_PASSWORD", ""),
        DbName:                       getEnv("DB_NAME", "data.db"),
        DbOrm:                        getEnv("DB_ORM", "entgo"),
        MigrationsDir:                getEnv("MIGRATIONS_DIR", "app/db/migrations"),
        GostSecret:                   getEnv("GOST_SECRET", "49cf26a7d274d62ad902ead6e69f5d71b4ffe703b4b07d25652c117cab74fcb1"),
        GostAuthRedirectAfterLogin:   getEnv("GOST_AUTH_REDIRECT_AFTER_LOGIN", "/profile"),
        GostAuthSessionExpiryInHours: getEnv("GOST_AUTH_SESSION_EXPIRY_IN_HOURS", "72"),
        GostAuthSkipVerify:           getEnvBool("GOST_AUTH_SKIP_VERIFY", true),
        BackendPkg:                   getEnv("GOST_BACKEND", "gin"),
    }, nil

	{{- else if eq .PreferredConfigFormat ".json"}}
		file, err := os.Open(ConfigFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		config := &Config{}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(config)
		return config, err

	{{- else if eq .PreferredConfigFormat ".toml"}}
		file, err := os.Open(ConfigFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		config := &Config{}
		decoder := toml.NewDecoder(file)
		err = decoder.Decode(config)
		return config, err

	{{- else if eq .PreferredConfigFormat ".yaml"}}
		file, err := os.Open(ConfigFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		config := &Config{}
		decoder := yaml.NewDecoder(file)
		return config, decoder.Decode(config)

	{{- end }}
}
`

func (g *GenConfPlugin) Init() error {
	g.Files = map[string]func() string{
		"app/cfg/cfg.go": func() string {
			tmpl, err := template.New("cfg").Parse(cfgTemplate)
			if err != nil {
				panic(err)
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, g.Data); err != nil {
				panic(err)
			}

			return buf.String()
		},
	}
	return nil
}

func (g *GenConfPlugin) Execute() error {
	return general.GenerateFiles(g.Data, g.Files)
}

func (g *GenConfPlugin) Shutdown() error {
	return nil
}

func (g *GenConfPlugin) Name() string {
	return "Configuration Code GenConfPlugin"
}

func (g *GenConfPlugin) Version() string {
	return "1.0.0"
}

func (g *GenConfPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenConfPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenConfPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenConfPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenConfPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/cfg"
}

func (g *GenConfPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenConfPlugin(data config.ProjectData) *GenConfPlugin {
	return &GenConfPlugin{
		Data: data,
	}
}
