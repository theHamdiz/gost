package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"

	"github.com/theHamdiz/gost/clr"
	"github.com/theHamdiz/gost/dirs"
	"github.com/theHamdiz/gost/dwn"
	"github.com/theHamdiz/gost/router"
	"github.com/theHamdiz/gost/runner"
	"github.com/theHamdiz/gost/seeder"

	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/theHamdiz/gost/gen"
)

func buildConfig(config *gen.GostConfig) {
	scanner := bufio.NewScanner(os.Stdin)

	existingConfig := isFirstRun()
	if existingConfig != nil {
		*config = *existingConfig
		fmt.Println(clr.Colorize("[✔] Welcome back to GoSt!", "green"))
		fmt.Println(clr.Colorize("[✔] Configuration loaded.", "green"))
	} else {
		fmt.Println(clr.Colorize("[✔] Welcome to GoSt! Your favorite go starter tool.", "green"))
	}

	if config.PreferredIDE == "" {
		config.PreferredIDE = askChoice(scanner, "[-] Your IDE of choice:", []string{"VSCode", "Goland", "IDEA", "Cursor", "Zed", "Sublime", "Vim", "Nvim", "Nano", "Notepad++", "Zeus", "LiteIDE", "Emacs", "Eclipse"})
	}
	if config.PreferredBackendFramework == "" {
		config.PreferredBackendFramework = askChoice(scanner, "[-] Choose your backend framework:", []string{"Gin", "Chi", "Echo", "StdLib"})
	}
	if config.PreferredDbDriver == "" {
		config.PreferredDbDriver = askChoice(scanner, "[-] Choose your preferred db driver:", []string{"Sqlite", "Postgresql", "MySql", "MongoDb"})
	}
	if config.PreferredDbOrm == "" {
		config.PreferredDbOrm = askChoice(scanner, "[-] Choose your preferred ORM:", []string{"Ent", "Gorm", "Bun", "Sqlc", "Bob", "None"})
	}
	if config.PreferredUiFramework == "" {
		config.PreferredUiFramework = askChoice(scanner, "[-] Choose your preferred UI framework:", []string{"Tailwindcss", "Bootstrap"})
	}
	if config.PreferredUiFramework == "Tailwindcss" && config.PreferredComponentsFramework == "" {
		config.PreferredComponentsFramework = askChoice(scanner, "[-] Choose your preferred components framework (tailwind only):", []string{"None", "DaisyUI", "Flowbite"})
	} else if config.PreferredUiFramework != "Tailwindcss" {
		config.PreferredComponentsFramework = "None"
	}
	if config.PreferredPort == 0 {
		config.PreferredPort = askIntChoice(scanner, "[-] Preferred Port:", []int{9630, 42069, 6666})
	}
	if config.GlobalSettings == "" {
		config.GlobalSettings = askChoice(scanner, "[-] Should we ask you every time you initiate a project or do you want to set your preferences globally?", []string{"Yes ask me", "No set it & forget it.", "Keep IDE settings only global.", "Keep IDE & port settings only global."})
	}
	if config.PreferredConfigFormat == "" {
		config.PreferredConfigFormat = askChoice(scanner, "[-] Preferred config file format:", []string{"env", "json", "toml"})
	}

	saveConfig(*config)
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
func saveConfig(config gen.GostConfig) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(clr.Colorize("Error getting user home directory", "red"))
		return
	}

	var filePath string
	switch config.PreferredConfigFormat {
	case "env":
		filePath = filepath.Join(usr.HomeDir, ".gost")
	case "json":
		filePath = filepath.Join(usr.HomeDir, ".gost.json")
	case "toml":
		filePath = filepath.Join(usr.HomeDir, ".gost.toml")
	default:
		fmt.Println(clr.Colorize("Invalid config format", "red"))
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
		fmt.Println(clr.Colorize("Error saving config: "+saveErr.Error(), "red"))
	} else {
		fmt.Println(clr.Colorize("Configuration saved to 👉 "+filePath, "green"))
	}
}

func isFirstRun() *gen.GostConfig {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(clr.Colorize("Error getting user home directory", "red"))
		return nil
	}

	envFilePath := filepath.Join(usr.HomeDir, ".gost")
	jsonFilePath := filepath.Join(usr.HomeDir, ".gost.json")
	tomlFilePath := filepath.Join(usr.HomeDir, ".gost.toml")

	if _, err := os.Stat(envFilePath); err == nil {
		config, err := gen.LoadFromEnv(envFilePath)
		if err == nil {
			return config
		}
		fmt.Println(clr.Colorize("Error loading config from .gost file", "red"))
	}
	if _, err := os.Stat(jsonFilePath); err == nil {
		config, err := gen.LoadFromJSON(jsonFilePath)
		if err == nil {
			return config
		}
		fmt.Println(clr.Colorize("Error loading config from .gost.json file", "red"))
	}
	if _, err := os.Stat(tomlFilePath); err == nil {
		config, err := gen.LoadFromTOML(tomlFilePath)
		if err == nil {
			return config
		}
		fmt.Println(clr.Colorize("Error loading config from .gost.toml file", "red"))
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

// installTailwind -> Try to install tailwindcss from nodefirst.
func installTailwind(projectDir string) error {
	if _, err := os.Stat(filepath.Join(projectDir, "node_modules")); os.IsNotExist(err) {
		return runner.RunCommand("npm", "install", "tailwind@latest", "--force")
	}
	return nil
}

// installBootstrap -> Try to install bootstrap from node first.
func installBootstrap(projectDir string) error {
	if _, err := os.Stat(filepath.Join(projectDir, "node_modules")); os.IsNotExist(err) {
		return runner.RunCommand("npm", "install", "bootstrap@latest", "--force")
	}
	return nil
}

// installHtmx -> Try to install htmx from node first.
func installHtmx(projectDir string) error {
	if _, err := os.Stat(filepath.Join(projectDir, "node_modules")); os.IsNotExist(err) {
		return runner.RunCommand("npm", "install", "htmx.org@latest", "--save", "--force")
	}
	return nil
}

// installAir -> install the air watcher framework.
func installAir(projectDir string) error {
	if _, err := os.Stat(filepath.Join(projectDir, ".air.toml")); os.IsNotExist(err) {
		return runner.RunCommand("go", "install", "github.com/air-verse/air@latest")
	}
	return nil
}

func getLatestVersion(packageName string) string {
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
func addCommands(config gen.GostConfig) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "gost",
		Short: "GoSt CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				config.AppName = args[0]
				buildConfig(&config)
			} else {
				fmt.Println("Please specify a command or an app name.")
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&config.AppName, "name", "n", "", "Application name")
	rootCmd.PersistentFlags().StringVarP(&config.PreferredUiFramework, "ui", "u", "", "UI framework")
	rootCmd.PersistentFlags().StringVarP(&config.PreferredComponentsFramework, "component", "c", "none", "Component framework")
	rootCmd.PersistentFlags().StringVarP(&config.PreferredBackendFramework, "backend", "b", "", "Backend framework")

	var createCmd = &cobra.Command{
		Use:   "create [app name] [ui framework] [component framework] [backend framework]",
		Short: "Create a new project",
		Args:  cobra.MaximumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				config.AppName = args[0]
			}
			if len(args) > 1 {
				config.PreferredUiFramework = args[1]
			}
			if len(args) > 2 {
				config.PreferredComponentsFramework = args[2]
			}
			if len(args) > 3 {
				config.PreferredBackendFramework = args[3]
			}
			buildConfig(&config)
			generateProject(config)
		},
	}

	rootCmd.AddCommand(createCmd)
	addDbCommands(rootCmd)
	addRunCommand(rootCmd)

	return rootCmd
}

func addDbCommands(rootCmd *cobra.Command) {
	var dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Database commands",
	}

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			// Call migrator.MigrateDB
			fmt.Println("Running database migrations")
		},
	}

	var seedCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seed the database",
		Run: func(cmd *cobra.Command, args []string) {
			// Call seeder.SeedDBData
			fmt.Println("Seeding database data")
		},
	}

	var fakeCmd = &cobra.Command{
		Use:   "fake",
		Short: "Fake database data",
		Run: func(cmd *cobra.Command, args []string) {
			// Call faker.FakeDBData
			fmt.Println("Faking database data")
		},
	}

	dbCmd.AddCommand(migrateCmd, seedCmd, fakeCmd)
	rootCmd.AddCommand(dbCmd)
}

func addRunCommand(rootCmd *cobra.Command) {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the project",
		Run: func(cmd *cobra.Command, args []string) {
			// Call runner.RunProject
			fmt.Println("Running project")
		},
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "r",
		Short: "Run the project",
		Run: func(cmd *cobra.Command, args []string) {
			// Call runner.RunProject
			fmt.Println("Running project")
		},
	})
}

func generateProject(config gen.GostConfig) {
	data := gen.TemplateData{
		AppName:             config.AppName,
		UiFramework:         config.PreferredUiFramework,
		ComponentsFramework: config.PreferredComponentsFramework,
		BackendPkg:          config.PreferredBackendFramework,
		DbDriver:            "sqlite3",
	}

	fmt.Println("App Name:", data.AppName)
	fmt.Println("UI Framework:", data.UiFramework)
	fmt.Println("Component Framework:", data.ComponentsFramework)
	fmt.Println("Backend Framework:", data.BackendPkg)
	fmt.Println("DB Framework:", data.DbDriver)

	if err := gen.GenerateFiles(data); err != nil {
		fmt.Printf("Error generating files: %v\n", err)
	}

	switch data.BackendPkg {
	case "gin":
		data.BackendImport = "github.com/gin-gonic/gin"
		data.BackendInit = "gin.Default()"
		data.VersionedBackendImport = fmt.Sprintf("%s@%s", data.BackendImport, getLatestVersion("github.com/gin-gonic/gin"))
		runner.RunCommand("go", "get", data.VersionedBackendImport)
		data.VersionedBackendImport = strings.Replace(data.VersionedBackendImport, "@", " ", 1)
	case "chi":
		data.BackendImport = "github.com/go-chi/chi/v5"
		data.BackendPkg = "chi"
		data.BackendInit = "chi.NewRouter()"
		data.VersionedBackendImport = fmt.Sprintf("%s@%s", data.BackendImport, getLatestVersion("github.com/go-chi/chi/v5"))
		runner.RunCommand("go", "get", data.VersionedBackendImport)
		data.VersionedBackendImport = strings.Replace(data.VersionedBackendImport, "@", " ", 1)
	case "echo":
		data.BackendImport = "github.com/labstack/echo/v5"
		data.BackendPkg = "echo"
		data.BackendInit = "echo.New()"
		data.VersionedBackendImport = fmt.Sprintf("%s@%s", data.BackendImport, "v5.0.0-20230722203903-ec5b858dab61")
		runner.RunCommand("go", "get", data.VersionedBackendImport)
		data.VersionedBackendImport = strings.Replace(data.VersionedBackendImport, "@", " ", 1)
	default:
		log.Fatalf("Unsupported backend framework: %s", config.PreferredBackendFramework)
	}

	projectDir, err := router.GetProjectPath(config.AppName)
	if err != nil {
		log.Fatal(err)
	}

	err = dirs.Generate(projectDir, data.AppName)
	if err != nil {
		panic(err)
	}

	fingerPrint, _ := gen.Fingerprint(config.AppName)
	data.Fingerprint = fingerPrint

	err = gen.GenerateFiles(data)
	if err != nil {
		log.Fatal(err)
	}

	err = seeder.DbInit(config.AppName)
	if err != nil {
		return
	}

	err = installFrameworks(projectDir)
	if err != nil {
		err = setupFrameworks(projectDir)
		if err != nil {
			panic("Could not install UI or Components frameworks using either npm or direct download, maybe check your internet connection!")
		}
	}

	runner.RunCommand("go", "mod", "tidy")

	if _, err := os.Stat(filepath.Join(projectDir, ".git")); os.IsNotExist(err) {
		runner.RunCommand("git", "init", projectDir)
	}

	fmt.Println("Project created successfully!")
	binary, err := config.GetIDEBinaryName()
	if err != nil {
		fmt.Println("Could not launch your preferred IDE: ", err)
		return
	}
	runner.RunCommand(binary, projectDir)
}