package config

type ProjectData struct {
	AppName                string
	BackendImport          string
	BackendInit            string
	BackendPkg             string
	ComponentsFramework    string
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
