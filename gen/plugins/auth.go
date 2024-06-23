package plugins

import (
	"fmt"
	"myframework/framework"
)

type AuthPlugin struct {
	authService *AuthService
}

func (ap *AuthPlugin) Init(config map[string]interface{}, sc *framework.ServiceContainer) error {
	ap.authService = sc.GetService("auth").(*AuthService)
	fmt.Println("Auth Plugin Initialized with cfg:", config)
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

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func main() {
	config, err := framework.LoadConfig("cfg.json")
	if err != nil {
		fmt.Println("Error loading cfg:", err)
		return
	}

	sc := framework.NewServiceContainer()
	authService := NewAuthService()
	sc.RegisterService("auth", authService)

	pm := framework.NewPluginManager(*config, sc)
	authPlugin := &AuthPlugin{}
	pm.RegisterPlugin("auth", authPlugin)

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
