package auth

import (
	"fmt"

	"github.com/theHamdiz/gost/plugins"
	"github.com/theHamdiz/gost/plugins/config"
	"github.com/theHamdiz/gost/services"
)

type PluginConfig struct {
	PluginDir string
}

type AuthPlugin struct {
	authService *AuthService
	config      *PluginConfig
}

func (ap *AuthPlugin) Init() error {
	var sc = services.NewServiceContainer()
	ap.authService = sc.GetService("auth").(*AuthService)
	ap.config.PluginDir = "plugins/auth"
	fmt.Println("Auth Plugin Initialized with cfg:", ap.config)
	return nil
}

func (ap *AuthPlugin) Execute() error {
	fmt.Println("Auth Plugin Executing")
	return nil
}

func (ap *AuthPlugin) Shutdown() error {
	fmt.Println("Auth Plugin Shutdown")
	return nil
}

func (g *AuthPlugin) Name() string {
	return "GenTypesPlugin"
}

func (g *AuthPlugin) Version() string {
	return "1.0.0"
}

func (g *AuthPlugin) Dependencies() []string {
	return []string{}
}

func (g *AuthPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *AuthPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *AuthPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *AuthPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/types"
}

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func RegisterAuthPlugin() {
	sc := services.NewServiceContainer()
	authService := NewAuthService()
	sc.RegisterService("auth", authService)
	pmConfig := &config.PluginManagerConfig{
		PluginsDir: "plugins",
	}
	pm := plugins.NewPluginManager(pmConfig)
	authPlugin := &AuthPlugin{}
	pm.RegisterPlugin(authPlugin)

	if err := pm.InitPlugins(); err != nil {
		fmt.Println("Error initializing plugins:", err)
		return
	}

	if err := pm.ExecutePlugins(); err != nil {
		fmt.Println("Error executing plugins:", err)
		return
	}

	if err := pm.ShutdownPlugins(); err != nil {
		fmt.Println("Error shutting down plugins:", err)
		return
	}
}
