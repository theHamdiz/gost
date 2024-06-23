package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/theHamdiz/gost/cfg"
	_ "github.com/theHamdiz/gost/cli"
	"github.com/theHamdiz/gost/clr"
	"github.com/theHamdiz/gost/dwn"
	genCfg "github.com/theHamdiz/gost/gen/config"
	"github.com/theHamdiz/gost/gen/dirs"
	"github.com/theHamdiz/gost/gen/files"
	"github.com/theHamdiz/gost/gen/fingerprint"
	"github.com/theHamdiz/gost/router"
	"github.com/theHamdiz/gost/runner"
	"github.com/theHamdiz/gost/seeder"
)

var Config *cfg.GostConfig

func BuildConfig(config *cfg.GostConfig) {
	scanner := bufio.NewScanner(os.Stdin)
	existingConfig := isFirstRun()
	var somethingChanged bool
	appName := config.AppName

	if existingConfig != nil {
		*config = *existingConfig
		config.AppName = appName
		fmt.Println(clr.Colorize("[✔] Welcome back to GoSt!", "green"))
		if config.GlobalSettings == "No set it & forget it." {
			fmt.Println(clr.Colorize("[✔] Configuration loaded.", "green"))
			return
		}
	} else {
		somethingChanged = true
		fmt.Println(clr.Colorize("[✔] Welcome to GoSt! Your favorite go starter tool.", "green"))
	}
	if config.PreferredIDE == "" {
		config.PreferredIDE = askChoice(scanner, "[-] Your IDE of choice:", []string{"VSCode", "Goland", "IDEA", "Cursor", "Zed", "Sublime", "Vim", "Nvim", "Nano", "Notepad++", "Zeus", "LiteIDE", "Emacs", "Eclipse"})
		somethingChanged = true
	}
	if config.PreferredBackendFramework == "" {
		config.PreferredBackendFramework = askChoice(scanner, "[-] Choose your backend framework:", []string{"Gin", "Chi", "Echo", "StdLib"})
		somethingChanged = true
	}
	if config.PreferredDbDriver == "" {
		config.PreferredDbDriver = askChoice(scanner, "[-] Choose your preferred db driver:", []string{"Sqlite", "Postgresql", "MySql", "MongoDb"})
		somethingChanged = true
	}
	if config.PreferredDbOrm == "" {
		config.PreferredDbOrm = askChoice(scanner, "[-] Choose your preferred ORM:", []string{"Ent", "Gorm", "Bun", "Sqlc", "Bob", "None"})
		somethingChanged = true
	}
	if config.PreferredUiFramework == "" {
		config.PreferredUiFramework = askChoice(scanner, "[-] Choose your preferred UI framework:", []string{"Tailwindcss", "Bootstrap"})
		somethingChanged = true
	}
	if config.PreferredUiFramework == "Tailwindcss" && config.PreferredComponentsFramework == "" {
		config.PreferredComponentsFramework = askChoice(scanner, "[-] Choose your preferred components framework (tailwind only):", []string{"None", "DaisyUI", "Flowbite"})
		somethingChanged = true
	} else if config.PreferredUiFramework != "Tailwindcss" {
		config.PreferredComponentsFramework = "None"
		somethingChanged = true
	}
	if config.PreferredPort == 0 {
		config.PreferredPort = askIntChoice(scanner, "[-] Preferred Port:", []int{9630, 42069, 6666})
		somethingChanged = true
	}
	if config.GlobalSettings == "" {
		config.GlobalSettings = askChoice(scanner, "[-] Should we ask you every time you initiate a project or do you want to set your preferences globally?", []string{"Yes ask me", "No set it & forget it.", "Keep IDE settings only global.", "Keep IDE & port settings only global."})
		somethingChanged = true
	}
	if config.PreferredConfigFormat == "" {
		config.PreferredConfigFormat = askChoice(scanner, "[-] Preferred cfg file format:", []string{"env", "json", "toml"})
		somethingChanged = true
	}

	if somethingChanged {
		saveConfig(*config)
	}
}

func askChoice(scanner *bufio.Scanner, question string, choices []string) string {
	fmt.Println(clr.Colorize(question, "black"))
	for i, choice := range choices {
		fmt.Printf("%d - %s\n", i+1, clr.Colorize(choice, "black"))
	}
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(choices) {
			return choices[choice-1]
		}
		fmt.Println(clr.Colorize("Invalid choice, please try again.", "black"))
	}
}

func askIntChoice(scanner *bufio.Scanner, question string, choices []int) int {
	fmt.Println(clr.Colorize(question, "black"))
	for i, choice := range choices {
		fmt.Printf("%d - %d\n", i+1, choice)
	}
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(choices) {
			return choices[choice-1]
		}
		fmt.Println(clr.Colorize("Invalid choice, please try again.", "black"))
	}
}

func saveConfig(config cfg.GostConfig) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(clr.Colorize("Error getting user home directory", "red"))
		return
	}

	var filePath string
	switch config.PreferredConfigFormat {
	case "env":
		filePath = filepath.Join(usr.HomeDir, ".gost.env")
	case "json":
		filePath = filepath.Join(usr.HomeDir, ".gost.json")
	case "toml":
		filePath = filepath.Join(usr.HomeDir, ".gost.toml")
	default:
		fmt.Println(clr.Colorize("Invalid cfg format", "red"))
		return
	}

	var saveErr error
	switch config.PreferredConfigFormat {
	case "env":
		saveErr = config.SaveAsEnv(filePath)
	case "json":
		saveErr = config.SaveAsJSON(filePath)
	case "toml":
		saveErr = config.SaveAsTOML(filePath)
	}

	if saveErr != nil {
		fmt.Println(clr.Colorize("Error saving cfg: "+saveErr.Error(), "red"))
	} else {
		fmt.Println(clr.Colorize("Configuration saved to 👉 "+filePath, "green"))
	}
}

func isFirstRun() *cfg.GostConfig {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(clr.Colorize("Error getting user home directory", "red"))
		return nil
	}

	envFilePath := filepath.Join(usr.HomeDir, ".gost.env")
	jsonFilePath := filepath.Join(usr.HomeDir, ".gost.json")
	tomlFilePath := filepath.Join(usr.HomeDir, ".gost.toml")

	if _, err := os.Stat(envFilePath); err == nil {
		config, err := cfg.LoadFromEnv(envFilePath)
		if err == nil {
			return config
		}
		fmt.Println(clr.Colorize("Error loading cfg from .gost.env file", "red"))
	}
	if _, err := os.Stat(jsonFilePath); err == nil {
		config, err := cfg.LoadFromJSON(jsonFilePath)
		if err == nil {
			return config
		}
		fmt.Println(clr.Colorize("Error loading cfg from .gost.json file", "red"))
	}
	if _, err := os.Stat(tomlFilePath); err == nil {
		config, err := cfg.LoadFromTOML(tomlFilePath)
		if err == nil {
			return config
		}
		fmt.Println(clr.Colorize("Error loading cfg from .gost.toml file", "red"))
	}

	return nil
}

func installFrameworks(projectDir string) error {
	err := installTailwind(projectDir)
	if err != nil {
		return err
	}

	err = installBootstrap(projectDir)
	if err != nil {
		return err
	}

	err = installHtmx(projectDir)
	if err != nil {
		return err
	}

	return nil
}

// installTailwind -> Try to install tailwindcss from node first.
func installTailwind(projectDir string) error {
	return runner.RunCommandWithDir(projectDir, "npm", "install", "tailwind@latest", "--force")
}

// installBootstrap -> Try to install bootstrap from node first.
func installBootstrap(projectDir string) error {
	return runner.RunCommandWithDir(projectDir, "npm", "install", "bootstrap@latest", "--force")
}

// installHtmx -> Try to install htmx from node first.
func installHtmx(projectDir string) error {
	return runner.RunCommandWithDir(projectDir, "npm", "install", "htmx.org@latest", "--save", "--force")
}

// installAir -> install the air watcher framework.
func installAir(projectDir string) error {
	if _, err := os.Stat(filepath.Join(projectDir, ".air.toml")); os.IsNotExist(err) {
		return runner.RunCommandWithDir(projectDir, "go", "install", "github.com/air-verse/air@latest")
	}
	return nil
}

func getLatestGoPackageVersion(packageName string) string {
	cmd := exec.Command("go", "list", "-m", "-versions", packageName)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error getting latest version of %s: %v", packageName, err)
	}

	versions := strings.Fields(string(output))
	return versions[len(versions)-1]
}

func setupFrameworks(projectDir string) error {
	err := dwn.DownloadFile("https://unpkg.com/htmx.org@latest", filepath.Join(projectDir, "app/assets/static/js/htmx.min.js"))
	if err != nil {
		log.Printf("Error downloading htmx: %v", err)
		return err
	}

	err = dwn.DownloadFile("https://cdn.tailwindcss.com", filepath.Join(projectDir, "app/assets/static/js/tailwind.min.js"))
	if err != nil {
		log.Printf("Error downloading htmx: %v", err)
		return err
	}

	return nil
}

func AddCommands(config cfg.GostConfig) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "gost",
		Short: "GoSt CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				config.AppName = args[0]
				BuildConfig(&config)
			} else {
				fmt.Println(clr.Colorize("Please specify a command or an app name.", "red"))
			}
		},
	}

	// Project Initialization Commands
	rootCmd.AddCommand(initCmd(config))
	rootCmd.AddCommand(createCmd(config))
	rootCmd.AddCommand(newCmd(config))

	// Add other command groups
	addDbCommands(rootCmd)
	addServerCommands(rootCmd)
	addConfigCommands(rootCmd)
	addGenerateCommands(rootCmd)
	addPluginCommands(rootCmd)
	addTestCommands(rootCmd)

	// Add your existing commands
	addRunCommand(rootCmd)

	return rootCmd
}

func GenerateProjectDir(config *cfg.GostConfig) {
	data := CreateProjectDataFromConfig(config)
	projectDir, err := router.GetProjectPath(data.AppName)
	if err != nil {
		log.Fatal(err)
	}

	if data.AppName == "" {
		panic(clr.Colorize("Please specify a project name!", "red"))
	}

	err = dirs.Generate(projectDir)
	if err != nil {
		panic(err)
	}

	fmt.Println(clr.Colorize(fmt.Sprintf("Project Dir created successfully 👉 %s", projectDir), "green"))
}
func initCmd(config cfg.GostConfig) *cobra.Command {
	return &cobra.Command{
		Use:     "init [app name]",
		Short:   "Initialize a new project",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				config.AppName = args[0]
			}
			BuildConfig(&config)
			GenerateProjectDir(&config)
			generateProject(&config)
		},
	}
}

func createCmd(config cfg.GostConfig) *cobra.Command {
	return &cobra.Command{
		Use:     "create project [app name]",
		Short:   "Create a new project",
		Long:    "Create a new go web project",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				config.AppName = args[0]
				fmt.Println("Create>> App Name:", config.AppName)
			} else {
				panic("Please specify a project name")
			}
			fmt.Println("Create>>After>> App Name:", config.AppName)
			BuildConfig(&config)
			GenerateProjectDir(&config)
			generateProject(&config)
		},
	}
}

func newCmd(config cfg.GostConfig) *cobra.Command {
	return &cobra.Command{
		Use:     "new project [app name]",
		Short:   "Create a new project",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"n"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				config.AppName = args[0]
			}
			BuildConfig(&config)
			GenerateProjectDir(&config)
			generateProject(&config)
		},
	}
}

func addDbCommands(rootCmd *cobra.Command) {
	var dbCmd = &cobra.Command{
		Use:     "db",
		Short:   "Database commands",
		Aliases: []string{"d"},
	}

	var migrateCmd = &cobra.Command{
		Use:     "migrate",
		Short:   "Run database migrations",
		Aliases: []string{"m"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call migrator.MigrateDB
			fmt.Println("Running database migrations")
			fmt.Println("Migrations Complete!")
		},
	}

	var seedCmd = &cobra.Command{
		Use:     "seed",
		Short:   "Seed the database",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call seeder.SeedDBData
			fmt.Println("Seeding database data")
			fmt.Println("Seeding Completed Successfully!")
		},
	}

	var rollbackCmd = &cobra.Command{
		Use:     "rollback",
		Short:   "Rollback database migrations",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call db.RollbackDB
			fmt.Println("Rolling back database migrations")
			fmt.Println("Rollback Completed Successfully!")
		},
	}

	dbCmd.AddCommand(migrateCmd, seedCmd, rollbackCmd)
	rootCmd.AddCommand(dbCmd)
}

func addServerCommands(rootCmd *cobra.Command) {
	var serverCmd = &cobra.Command{
		Use:     "server",
		Short:   "Server commands",
		Aliases: []string{"s"},
	}

	var startCmd = &cobra.Command{
		Use:     "start",
		Short:   "Start the server",
		Aliases: []string{"st"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call server.Start
			fmt.Println("Starting the server...")
		},
	}

	var stopCmd = &cobra.Command{
		Use:     "stop",
		Short:   "Stop the server",
		Aliases: []string{"sp"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call server.Stop
			fmt.Println("Stopping the server...")
		},
	}

	serverCmd.AddCommand(startCmd, stopCmd)
	rootCmd.AddCommand(serverCmd)
}

func addConfigCommands(rootCmd *cobra.Command) {
	var configCmd = &cobra.Command{
		Use:     "config",
		Short:   "Configuration commands",
		Aliases: []string{"cfg", "cnfg"},
	}

	var setCmd = &cobra.Command{
		Use:     "set <key> <value>",
		Short:   "Set a configuration value",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Call config.Set
			fmt.Printf("Setting %s to %s\n", args[0], args[1])
		},
	}

	var getCmd = &cobra.Command{
		Use:     "get <key>",
		Short:   "Get a configuration value",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Call config.Get
			fmt.Printf("Value of %s\n", args[0])
		},
	}

	var removeCmd = &cobra.Command{
		Use:     "remove <key>",
		Short:   "Remove a configuration value",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Call config.Remove
			fmt.Printf("Removing %s\n", args[0])
		},
	}

	var showCmd = &cobra.Command{
		Use:     "show",
		Short:   "Show configuration",
		Aliases: []string{"sh"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call config.Show
			cfg.ShowConfig("")
		},
	}

	var resetCmd = &cobra.Command{
		Use:     "reset",
		Short:   "Reset configuration",
		Aliases: []string{"r", "rs"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call config.Reset
			fmt.Println("Resetting configuration...")
			cfg.ResetConfig("")
			fmt.Println("Reset Completed Successfully!")
		},
	}

	configCmd.AddCommand(setCmd, getCmd, removeCmd, showCmd, resetCmd)
	rootCmd.AddCommand(configCmd)
}

func addGenerateCommands(rootCmd *cobra.Command) {
	var generateCmd = &cobra.Command{
		Use:     "generate",
		Short:   "Generate a new component",
		Aliases: []string{"g", "gen"},
	}

	var modelCmd = &cobra.Command{
		Use:     "model <name> [fields]",
		Short:   "Generate a new model",
		Aliases: []string{"m", "mod", "mdl", "md"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fields := args[1:]
			// Call generator.GenerateModel
			fmt.Printf("Generating model named %s with fields %v\n", name, fields)
		},
	}

	var viewCmd = &cobra.Command{
		Use:     "view <name>",
		Short:   "Generate a new view",
		Aliases: []string{"v", "vu", "vw", "veiw"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			// Call generator.GenerateView
			fmt.Printf("Generating view named %s\n", name)
		},
	}

	var handlerCmd = &cobra.Command{
		Use:     "handler <name>",
		Short:   "Generate a new handler",
		Aliases: []string{"h", "hnd", "hndl", "hdlr", "handlre"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			// Call generator.GenerateHandler
			fmt.Printf("Generating handler named %s\n", name)
		},
	}

	var eventCmd = &cobra.Command{
		Use:     "event <name>",
		Short:   "Generate a new event",
		Aliases: []string{"e", "ev", "eve", "evn", "evnt"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			// Call generator.GenerateEvent
			fmt.Printf("Generating event named %s\n", name)
		},
	}

	var pluginCmd = &cobra.Command{
		Use:     "plugin <name>",
		Short:   "Generate a new plugin",
		Aliases: []string{"p", "pg", "pgn", "plug", "plug-in"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			// Call generator.GeneratePlugin
			fmt.Printf("Generating plugin named %s\n", name)
		},
	}

	var migrationCmd = &cobra.Command{
		Use:     "migration <name> [fields]",
		Short:   "Generate a new migration",
		Aliases: []string{"m", "mig", "migrate", "migrat"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fields := args[1:]
			// Call generator.GenerateMigration
			fmt.Printf("Generating migration named %s with fields %v\n", name, fields)
		},
	}

	var resourceCmd = &cobra.Command{
		Use:     "resource <name> [fields]",
		Short:   "Generate a new resource",
		Aliases: []string{"r", "rs", "resorce", "resourc", "resurce", "resurc"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fields := args[1:]
			// Call generator.GenerateResource
			fmt.Printf("Generating resource %s with fields %v\n", name, fields)
			// Generate associated model, migration, handler, and view
			fmt.Printf("Generating model %s with fields %v\n", name, fields)
			fmt.Printf("Generating migration for %s with fields %v\n", name, fields)
			fmt.Printf("Generating handler %sHandler\n", name)
			fmt.Printf("Generating view %sView\n", name)
		},
	}

	generateCmd.AddCommand(modelCmd, viewCmd, handlerCmd, eventCmd, pluginCmd, migrationCmd, resourceCmd)
	rootCmd.AddCommand(generateCmd)
}

func addPluginCommands(rootCmd *cobra.Command) {
	var pluginCmd = &cobra.Command{
		Use:     "plugin",
		Short:   "Plugin management commands",
		Aliases: []string{"plg"},
	}

	var installCmd = &cobra.Command{
		Use:     "install <plugin_name>",
		Short:   "Install a plugin",
		Aliases: []string{"i"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Call plugin.Install
			fmt.Printf("Installing plugin %s\n", args[0])
		},
	}

	var uninstallCmd = &cobra.Command{
		Use:     "uninstall <plugin_name>",
		Short:   "Uninstall a plugin",
		Aliases: []string{"u"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Call plugin.Uninstall
			fmt.Printf("Uninstalling plugin %s\n", args[0])
		},
	}

	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List installed plugins",
		Aliases: []string{"l", "ls"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call plugin.List
			fmt.Println("Listing installed plugins")
		},
	}

	pluginCmd.AddCommand(installCmd, uninstallCmd, listCmd)
	rootCmd.AddCommand(pluginCmd)
}

func addTestCommands(rootCmd *cobra.Command) {
	var testCmd = &cobra.Command{
		Use:     "test",
		Short:   "Run tests",
		Aliases: []string{"t", "ts"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call test.RunTests
			fmt.Println("Running tests...")
		},
	}

	rootCmd.AddCommand(testCmd)
}

// Existing user commands
func addRunCommand(rootCmd *cobra.Command) {
	var runCmdFull = &cobra.Command{
		Use:     "run",
		Short:   "Run the project",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			// Call runner.RunProject
			fmt.Println("Running project...")
			err := runner.RunProject(Config.AppName)
			if err != nil {
				fmt.Println("Error running project:", err)
				return
			}
		},
	}

	rootCmd.AddCommand(runCmdFull)
}

// CapitalizeFirstLetter capitalizes the first letter of the input string.
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	first := unicode.ToUpper(runes[0])
	rest := string(runes[1:])

	return string(first) + rest
}

func CreateProjectDataFromConfig(config *cfg.GostConfig) *genCfg.ProjectData {
	return &genCfg.ProjectData{
		AppName:             CapitalizeFirstLetter(config.AppName),
		BackendPkg:          config.PreferredBackendFramework,
		ComponentsFramework: config.PreferredComponentsFramework,
		CurrentYear:         time.Now().Year(),
		DbDriver:            config.PreferredDbDriver,
		UiFramework:         config.PreferredUiFramework,
		Port:                config.PreferredPort,
		DbOrm:               config.PreferredDbOrm,
		IncludeAuth:         true,
		MigrationsDir:       "app/db/migrations",
	}
}
func generateProject(config *cfg.GostConfig) {
	data := CreateProjectDataFromConfig(config)

	fmt.Println("App Name:", data.AppName)
	fmt.Println("UI Framework:", data.UiFramework)
	fmt.Println("Component Framework:", data.ComponentsFramework)
	fmt.Println("Backend Framework:", data.BackendPkg)
	fmt.Println("DB Framework:", data.DbDriver)
	fmt.Println("DB Orm:", data.DbOrm)

	if err := files.GenerateFiles(*data); err != nil {
		fmt.Printf("Error generating files: %v\n", err)
	}

	switch strings.ToLower(data.BackendPkg) {
	case "gin":
		data.BackendImport = "github.com/gin-gonic/gin"
		data.BackendInit = "gin.Default()"
		data.VersionedBackendImport = fmt.Sprintf("%s@%s", data.BackendImport, getLatestGoPackageVersion("github.com/gin-gonic/gin"))
		err := runner.RunCommand("go", "get", data.VersionedBackendImport)
		if err != nil {
			fmt.Printf("Error running go get: %+v\n", err)
			return
		}
		data.VersionedBackendImport = strings.Replace(data.VersionedBackendImport, "@", " ", 1)
	case "chi":
		data.BackendImport = "github.com/go-chi/chi/v5"
		data.BackendPkg = "chi"
		data.BackendInit = "chi.NewRouter()"
		data.VersionedBackendImport = fmt.Sprintf("%s@%s", data.BackendImport, getLatestGoPackageVersion("github.com/go-chi/chi/v5"))
		err := runner.RunCommand("go", "get", data.VersionedBackendImport)
		if err != nil {
			fmt.Println("Error running go get:", err)
			return
		}
		data.VersionedBackendImport = strings.Replace(data.VersionedBackendImport, "@", " ", 1)
	case "echo":
		data.BackendImport = "github.com/labstack/echo/v5"
		data.BackendPkg = "echo"
		data.BackendInit = "echo.New()"
		data.VersionedBackendImport = fmt.Sprintf("%s@%s", data.BackendImport, "v5.0.0-20230722203903-ec5b858dab61")
		err := runner.RunCommand("go", "get", data.VersionedBackendImport)
		if err != nil {
			fmt.Println("Error running go get:", err)
			return
		}
		data.VersionedBackendImport = strings.Replace(data.VersionedBackendImport, "@", " ", 1)
	default:
		log.Fatalf("Unsupported backend framework: %s", config.PreferredBackendFramework)
	}

	fingerPrint, _ := fingerprint.Fingerprint(config.AppName)
	data.Fingerprint = fingerPrint

	err := files.GenerateFiles(*data)
	if err != nil {
		log.Fatal(err)
	}

	err = seeder.DbInit(config.AppName)
	if err != nil {
		return
	}

	err = installFrameworks(data.ProjectDir)
	if err != nil {
		err = setupFrameworks(data.ProjectDir)
		if err != nil {
			panic("Could not install UI or Components frameworks using either npm or direct download, maybe check your internet connection!")
		}
	}

	err = runner.RunCommand("go", "mod", "tidy")
	if err != nil {
		fmt.Println("Error running go mod tidy:", err)
		return
	}

	if _, err := os.Stat(filepath.Join(data.ProjectDir, ".git")); os.IsNotExist(err) {
		err := runner.RunCommand("git", "init", data.ProjectDir)
		if err != nil {
			fmt.Println("Error running git init:", err)
			return
		}
	}

	err = runner.RunCommand("npm", "audit", "fix")
	if err != nil {
		fmt.Println("Error running npm audit fix:", err)
		return
	}

	fmt.Println("Project created successfully!")
	binary, err := config.GetIDEBinaryName()
	if err != nil {
		fmt.Println("Could not launch your preferred IDE: ", err)
		return
	}
	err = runner.RunCommand(binary, data.ProjectDir)
	if err != nil {
		fmt.Println("Error running go get:", err)
		return
	}
}
