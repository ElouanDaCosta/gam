/*
Copyright © 2024 Elouan DA COSTA PEIXOTO elouandacostapeixoto@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ElouanDaCosta/gam/templates"
	"github.com/ElouanDaCosta/gam/utils"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Folder struct {
	Name       string   `mapstructure:"name"`
	Subfolders []Folder `mapstructure:"subfolders"`
}

type Config struct {
	ServiceName string   `mapstructure:"service_name"`
	Folders     []Folder `mapstructure:"folders"`
}

type promptContent struct {
	errorMsg string
	label    string
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a golang application.",
	Long: `Initialize a golang application using the specified technology between gin, gRPC or just basic http. For example:

gam init 
gam init --name [your_app_name]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		appName, _ := cmd.Flags().GetString("name")
		dockerfile, _ := cmd.Flags().GetBool("dockerfile")

		if appName != "" {
			if dockerfile {
				generateFromStructureFile(appName, true)
			} else {
				generateFromStructureFile(appName, false)
			}
		} else {
			generateFromStructureFile("new_app", false)
		}
	},
}

func generateFromStructureFile(appName string, dockerfile bool) {
	newService := exec.Command("mkdir", appName)
	_, newServiceErr := newService.Output()

	if newServiceErr != nil {
		fmt.Println("Error creating the directory:", newServiceErr.Error())
		return
	}

	currentPath := utils.GetAbsolutePath()

	if err := os.Chdir(appName); err != nil {
		log.Fatalf("unable to change directory to %s, %v", appName, err)
	}

	runGoModInit(appName)

	appType := promptContent{
		"Please select a package.",
		"Which package do you want your app to be based of ?",
	}

	newAppType := askUserForPackage(appType)
	config, configErr := readStructureFile(newAppType, installedPath)
	if configErr != nil {
		fmt.Println(configErr)
	}
	createFolders(currentPath+"/"+appName, config.Folders)
	log.Println("Application structure created")
	addPackageToApp(newAppType, appName)
	if dockerfile {
		createDockerfile(currentPath + "/" + appName)
	}

	writeInSaveAppFile(appName, installedPath, currentPath)

	fmt.Printf("%s created successfully\n", appName)
}

func runGoModInit(serviceName string) {
	cmd := exec.Command("go", "mod", "init", serviceName)
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run go mod init: %v", err)
	} else {
		log.Println("Go application initialized")
	}
}

func createFolders(basePath string, folders []Folder) {
	os.Chdir(basePath)
	for _, folder := range folders {
		// folderPath := fmt.Sprintf("%s/%s", basePath, folder.Name)
		os.Mkdir(folder.Name, 0755)
		createFolders(basePath, folder.Subfolders)
	}
}

// pass the structure file to the flag without the extension
func readStructureFile(appType string, basePath string) (Config, error) {
	os.Chdir(basePath)
	viper.AddConfigPath("configs")
	switch appType {
	case "gin":
		viper.SetConfigName("config-gin")
	case "gRPC":
		viper.SetConfigName("config-grpc")
	case "basic http":
		viper.SetConfigName("config-http")
	}

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		return Config{}, err
	}
	log.Println("Config file found")
	return config, nil
}

func askUserForPackage(pc promptContent) string {
	items := []string{"gin", "gRPC", "basic http"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label: pc.label,
			Items: items,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}

func addPackageToApp(appType string, newAppBasePath string) {
	os.Chdir(newAppBasePath)
	exec.Command("touch", "main.go").Output()
	if appType == "gin" {
		exec.Command("go", "get", "-u", "github.com/gin-gonic/gin@latest").Output()
		writeMainGo(newAppBasePath, "gin")
	}
	if appType == "gRPC" {
		exec.Command("go", "get", "-u", "google.golang.org/grpc").Output()
		writeMainGo(newAppBasePath, "gRPC")
	}
	if appType == "basic http" {
		writeMainGo(newAppBasePath, "basic http")
	}
}

func writeInSaveAppFile(appName string, basePath string, currentPath string) {
	os.Chdir(basePath + "/storage")
	f, err := os.OpenFile("app.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write([]byte("name: " + appName + "\n"))
	if err != nil {
		log.Fatal(err)
	}

	f.Write([]byte("app path: " + currentPath + "/" + appName + "\n\n"))

	f.Close()
}

func createDockerfile(appName string) {
	errWd := os.Chdir(appName)

	if errWd != nil {
		fmt.Println(errWd)
	}

	f, err := os.OpenFile("dockerfile", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	appVersion, err := utils.GetGoVersion(appName)

	appVersionSplit := strings.Split(appVersion, " ")

	if err != nil {
		fmt.Println(err)
	}

	var content = fmt.Sprintf("%v", templates.RenderDockerfileTemplate(appVersionSplit[1]))

	_, err = f.WriteString(content)

	if err != nil {
		fmt.Println(err)
	} else {
		log.Println("dockerfile generated successfully.")
	}
}

func writeMainGo(basePath string, appType string) {
	os.Chdir(basePath)
	f, err := os.OpenFile("main.go", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	content := ""

	switch appType {
	case "gin":
		content = fmt.Sprintf("%v", templates.RenderGinTemplate())
	case "gRPC":
		content = fmt.Sprintf("%v", templates.RenderGrpcTemplate())
	case "basic http":
		content = fmt.Sprintf("%v", templates.RenderHttpTemplate())
	}

	_, err = f.WriteString(content)

	if err != nil {
		fmt.Println(err)
	} else {
		log.Println("main.go generated successfully.")
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().String("name", "", "Generate an app with the given name (default new_app)")
	generateCmd.Flags().BoolP("dockerfile", "d", false, "Generate a dockerfile in the new app directory")
}
