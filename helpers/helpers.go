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
	"sync"
	"time"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/theHamdiz/gost/cfg"
	_ "github.com/theHamdiz/gost/cli"
	"github.com/theHamdiz/gost/clr"
	"github.com/theHamdiz/gost/codegen"
	"github.com/theHamdiz/gost/codegen/dirs"
	"github.com/theHamdiz/gost/codegen/fingerprint"
	genCfg "github.com/theHamdiz/gost/config"
	"github.com/theHamdiz/gost/dwn"
	"github.com/theHamdiz/gost/git"
	"github.com/theHamdiz/gost/npm"
	"github.com/theHamdiz/gost/router"
	"github.com/theHamdiz/gost/runner"
	"github.com/theHamdiz/gost/seeder"
	"github.com/theHamdiz/gost/tailwind"
)

// TODO: implement an env checker that can install anything for the user.
type EnvChecker interface {
	CheckDotEnvExists() bool
	CheckGoInstalled() bool
	CheckGitInstalled() bool
	CheckNodeInstalled() bool
	CheckAirInstalled() bool
	CheckNpmInstalled() bool
}

var Config *cfg.GostConfig
var ProjectData *genCfg.ProjectData

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
		config.PreferredDbOrm = askChoice(scanner, "[-] Choose your preferred ORM:", []string{"Built In", "Ent", "Gorm", "Bun", "Sqlc", "Bob"})
		somethingChanged = true
	}
	if config.PreferredFrontEndFramework == "" {
		config.PreferredUiFramework = askChoice(scanner, "[-] Choose your preferred frontend framework:", []string{"Htmx", "React", "Svelte", "Vue"})
		somethingChanged = true
	}
	if config.PreferredUiFramework == "" {
		config.PreferredUiFramework = askChoice(scanner, "[-] Choose your preferred UI framework:", []string{"Tailwindcss", "Bootstrap"})
		somethingChanged = true
	}
	if config.PreferredUiFramework == "Tailwindcss" && config.PreferredComponentsFramework == "" {
		config.PreferredComponentsFramework = askChoice(scanner, "[-] Choose your preferred components framework (tailwind only):", []string{"None", "DaisyUI", "Flowbite", "PrelineUI", "TW-Elements"})
		somethingChanged = true
	} else if config.PreferredUiFramework != "Tailwindcss" {
		config.PreferredComponentsFramework = "None"
		somethingChanged = true
	}
	if config.PreferredPort == 0 {
		config.PreferredPort = askIntChoice(scanner, "[-] Preferred Port:", []int{9630, 42069, 6666, 8080})
		somethingChanged = true
	}
	if config.GlobalSettings == "" {
		config.GlobalSettings = askChoice(scanner, "[-] Should we ask you every time you initiate a project or do you want to set your preferences globally?", []string{"Yes ask me", "No set it & forget it.", "Keep IDE settings only global.", "Keep IDE & port settings only global."})
		somethingChanged = true
	}
	if config.PreferredConfigFormat == "" {
		config.PreferredConfigFormat = askChoice(scanner, "[-] Preferred cfg file format:", []string{"env", "json", "toml", "yaml"})
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
	case "yaml":
		filePath = filepath.Join(usr.HomeDir, ".gost.yaml")
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
	case "yaml":
		saveErr = config.SaveAsYAML(filePath)
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
	switch strings.ToLower(ProjectData.UiFramework) {
	case "tailwind", "tailwindcss":
		err := installTailwind(projectDir)
		if err != nil {
			fmt.Printf(clr.Colorize("Error installing tailwind: %s\n", "red"), err)
		}
	case "bootstrap", "bootstrapjs", "bootstrapcss":
		err := installBootstrap(projectDir)
		if err != nil {
			fmt.Printf(clr.Colorize("Error installing bootstrap: %s\n", "red"), err)
		}
	}

	err := installHtmx(projectDir)
	if err != nil {
		fmt.Printf(clr.Colorize("Error installing htmx: %s\n", "red"), err)
	}

	err = installAir(projectDir)
	if err != nil {
		fmt.Printf(clr.Colorize("Error installing air: %s\n", "red"), err)
	}
	return nil
}

// installTailwind -> Try to install tailwindcss from node first.
func installTailwind(projectDir string) error {
	return runner.RunCommandWithDir(projectDir, "npm", "i", "tailwind@latest", "--force")
}

// installBootstrap -> Try to install bootstrap from node first.
func installBootstrap(projectDir string) error {
	return runner.RunCommandWithDir(projectDir, "npm", "i", "bootstrap@latest", "--force")
}

// installHtmx -> Try to install htmx from node first.
func installHtmx(projectDir string) error {
	return runner.RunCommandWithDir(projectDir, "npm", "i", "htmx.org@2.0.0", "--save", "--force")
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

func GenerateProjectDir(config *cfg.GostConfig) error {
	ProjectData = NewProjectDataFromConfig(config)
	ProjectData.AppName = strings.ToLower(config.AppName)
	if ProjectData.AppName == "" {
		panic(clr.Colorize("Please specify a project name!", "red"))
	}

	projectDir, err := router.GetProjectPath(ProjectData.AppName)
	if err != nil {
		log.Fatal(err)
	}

	ProjectData.ProjectDir = projectDir

	if _, err := os.Stat(ProjectData.ProjectDir); err == nil {
		return fmt.Errorf(">>Gost>> Project already exists")
	}

	dirsGenerator := dirs.NewDirsGenerator()
	err = dirsGenerator.Generate(projectDir)
	if err != nil {
		panic(err)
	}

	fmt.Println(clr.Colorize(fmt.Sprintf("Project Dir created successfully 👉 %s", projectDir), "green"))
	return nil
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
			err := GenerateProjectDir(&config)
			if err != nil {
				fmt.Println(clr.Colorize(err.Error(), "red"))
				return
			}
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
			} else {
				panic("Please specify a project name")
			}
			BuildConfig(&config)
			err := GenerateProjectDir(&config)
			if err != nil {
				fmt.Println(clr.Colorize(err.Error(), "red"))
				return
			}
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
			err := GenerateProjectDir(&config)
			if err != nil {
				fmt.Println(clr.Colorize(err.Error(), "red"))
				return
			}
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
		Aliases: []string{"g", "codegen"},
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
			fmt.Println(">>Gost>> Running tests...")
			err := runner.RunTests(ProjectData.ProjectDir)
			if err != nil {
				fmt.Printf(clr.Colorize("Error running tests: %v\n", "red"), err)
			}
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

func NewProjectDataFromConfig(config *cfg.GostConfig) *genCfg.ProjectData {
	return &genCfg.ProjectData{
		AppName:               CapitalizeFirstLetter(config.AppName),
		BackendPkg:            config.PreferredBackendFramework,
		ComponentsFramework:   config.PreferredComponentsFramework,
		PreferredConfigFormat: config.PreferredConfigFormat,
		ConfigFile:            config.PreferredConfigFormat,
		CurrentYear:           time.Now().Year(),
		DbDriver:              config.PreferredDbDriver,
		UiFramework:           config.PreferredUiFramework,
		Port:                  config.PreferredPort,
		DbOrm:                 config.PreferredDbOrm,
		IncludeAuth:           true,
		MigrationsDir:         "app/db/migrations",
	}
}

func generateProject(config *cfg.GostConfig) {
	fmt.Println(clr.Colorize("App Name:", "teal"), clr.Colorize(ProjectData.AppName, "green"))
	fmt.Println(clr.Colorize("Project Directory:", "teal"), clr.Colorize(ProjectData.ProjectDir, "green"))
	fmt.Println(clr.Colorize("Project Config File:", "teal"), clr.Colorize(ProjectData.ConfigFile, "green"))
	fmt.Println(clr.Colorize("UI Framework:", "teal"), clr.Colorize(ProjectData.UiFramework, ""))
	fmt.Println(clr.Colorize("Component Framework:", "teal"), clr.Colorize(ProjectData.ComponentsFramework, ""))
	fmt.Println(clr.Colorize("Frontend Framework:", "teal"), clr.Colorize(ProjectData.FrontEndFramework, ""))
	fmt.Println(clr.Colorize("Backend Framework:", "teal"), clr.Colorize(ProjectData.BackendPkg, ""))
	fmt.Println(clr.Colorize("DB Driver:", "teal"), clr.Colorize(ProjectData.DbDriver, ""))
	fmt.Println(clr.Colorize("DB Orm:", "teal"), clr.Colorize(ProjectData.DbOrm, ""))

	if err := codegen.ExecuteGeneration(*ProjectData); err != nil {
		fmt.Printf("Error generating files: %v\n", err)
	}

	switch strings.ToLower(ProjectData.BackendPkg) {
	case "gin":
		ProjectData.BackendImport = "github.com/gin-gonic/gin"
		ProjectData.BackendInit = "gin.Default()"
		ProjectData.VersionedBackendImport = fmt.Sprintf("%s@%s", ProjectData.BackendImport, getLatestGoPackageVersion("github.com/gin-gonic/gin"))
		ProjectData.VersionedBackendImport = strings.Replace(ProjectData.VersionedBackendImport, "@", " ", 1)
	case "chi":
		ProjectData.BackendImport = "github.com/go-chi/chi/v5"
		ProjectData.BackendPkg = "chi"
		ProjectData.BackendInit = "chi.NewRouter()"
		ProjectData.VersionedBackendImport = fmt.Sprintf("%s@%s", ProjectData.BackendImport, getLatestGoPackageVersion("github.com/go-chi/chi/v5"))
		ProjectData.VersionedBackendImport = strings.Replace(ProjectData.VersionedBackendImport, "@", " ", 1)
	case "echo":
		ProjectData.BackendImport = "github.com/labstack/echo/v5"
		ProjectData.BackendPkg = "echo"
		ProjectData.BackendInit = "echo.New()"
		ProjectData.VersionedBackendImport = fmt.Sprintf("%s@%s", ProjectData.BackendImport, "v5.0.0-20230722203903-ec5b858dab61")
		ProjectData.VersionedBackendImport = strings.Replace(ProjectData.VersionedBackendImport, "@", " ", 1)
	default:
		log.Fatalf("Unsupported backend framework: %s", config.PreferredBackendFramework)
	}

	fingerPrint, _ := fingerprint.Fingerprint(config.AppName)
	ProjectData.Fingerprint = fingerPrint

	err := codegen.ExecuteGeneration(*ProjectData)
	if err != nil {
		log.Fatal(err)
	}

	err = seeder.DbInit(config.AppName)
	if err != nil {
		fmt.Printf(clr.Colorize("Error seeding database: %v\n", "red"), err)
	}

	err = npm.CheckNPMInstalled()
	if err != nil {
		fmt.Println(err)
	}

	frontEndDir := filepath.Join(ProjectData.ProjectDir, "app/web/frontend/")
	backendDir := filepath.Join(ProjectData.ProjectDir, "app/web/backend/")

	err = installFrameworks(frontEndDir)
	if err != nil {
		err = setupFrameworks(frontEndDir)
		if err != nil {
			fmt.Println(clr.Colorize("Could not install UI or Components frameworks using either npm or direct download, maybe check your internet connection!", "red"))
		}
	}

	// Here install components framework & its config.
	if ProjectData.ComponentsFramework != "" {
		err = runner.RunCommand(tailwind.GetInstallationCommandFor(ProjectData.ComponentsFramework))
		if err == nil {
			contentStr := tailwind.GetContentConfigForComponentFramework(ProjectData.ComponentsFramework)
			if contentStr != "" {
				_ = tailwind.AppendToTailwindConfig(filepath.Join(frontEndDir, "tailwind.config.js"), "content", contentStr)
			}

			pluginStr := tailwind.GetContentConfigForComponentFramework(ProjectData.ComponentsFramework)
			if pluginStr != "" {
				_ = tailwind.AppendToTailwindConfig(filepath.Join(frontEndDir, "tailwind.config.js"), "plugins", pluginStr)
			}
		}
	}

	err = git.CheckGitInstalled()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(clr.Colorize("Project created successfully!", "teal"))
	binary, ideErr := config.GetIDEBinaryName()
	var ideCommands [][]string
	if ideErr == nil {
		// launch preferred ide
		ideCommands = [][]string{
			{binary, ProjectData.ProjectDir},
		}
	}

	// Group commands by program
	npmCommands := [][]string{
		{"npm", "i"},
		{"npm", "audit", "fix"},
	}

	// init a git repo & create a main branch instead of master
	gitCommands := [][]string{
		{"git", "init", ProjectData.ProjectDir},
		{"git", "branch", "-m", "main"},
	}

	// get all dependencies & tidy go mod
	goCommands := [][]string{
		{"go", "get", "-u", "./..."},
		{"go", "mod", "tidy"},
	}

	var wg sync.WaitGroup
	results := make(chan error, len(npmCommands)*2+len(gitCommands)+len(goCommands))

	// Function to run a command and send result to the channel
	runCommand := func(cmd []string) {
		defer wg.Done()
		var cmdDirs []string
		if cmd[0] == "npm" {
			cmdDirs = append(cmdDirs, frontEndDir)
			cmdDirs = append(cmdDirs, backendDir)
		} else {
			cmdDirs = append(cmdDirs, ProjectData.ProjectDir)
		}

		for _, dir := range cmdDirs {
			if err := runner.RunCommandWithDir(dir, cmd[0], cmd[1:]...); err != nil {
				results <- err
			} else {
				results <- nil
			}
		}
	}

	// Function to run commands sequentially within a group
	runCommandGroup := func(commands [][]string) {
		for _, cmd := range commands {
			wg.Add(1)
			runCommand(cmd)
			wg.Wait() // Wait for each command to finish before starting the next
		}
	}

	// Run npm commands concurrently
	go func() {
		runCommandGroup(npmCommands)
	}()

	// Run git commands concurrently
	go func() {
		runCommandGroup(gitCommands)
	}()

	// Run go commands concurrently
	go func() {
		runCommandGroup(goCommands)
	}()

	// only call ide commands if there's no error.
	if ideErr == nil {
		go func() {
			runCommandGroup(ideCommands)
		}()
	}
	// Wait for all commands to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for err := range results {
		if err != nil {
			fmt.Println(clr.Colorize("Error:", "red"), err)
		}
	}
}
