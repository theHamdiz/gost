package config

type Configurable interface {
	GetAppName() string
}

type ProjectData struct {
	AppName                string
	BackendImport          string
	BackendInit            string
	BackendPkg             string
	ComponentsFramework    string
	ConfigFile             string
	CurrentYear            int
	DbDriver               string
	DbOrm                  string
	Fingerprint            string
	FrontEndFramework      string
	IncludeAuth            bool
	MigrationsDir          string
	ProjectDir             string
	Port                   int
	UiFramework            string
	VersionedBackendImport string
}

type ResourcePluginConfig struct {
	AppName      string
	ResourceName string
	ResourceType string
}

func (p *ProjectData) GetAppName() string {
	return p.AppName
}

func (r *ResourcePluginConfig) GetAppName() string {
	return r.AppName
}
