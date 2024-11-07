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

	"github.com/spf13/cobra"
)

// upgradeVersionCmd represents the upgradeVersion command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the go version of the specified application.",
	Long: `Upgrade a single application go version or a specified package version or update all applications go version or package. For example:

gam upgrade --name [your_app_name] --version [version_wanted]

gam upgrade --version [version_wanted] --all

gam upgrade --version [version_wanted] -a
`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		newVersion, _ := cmd.Flags().GetString("version")
		allApp, _ := cmd.Flags().GetBool("all")

		if name != "" {
			if newVersion != "" {
				appPath := getAppPath(name)
				bumpOneGoVersion(appPath, newVersion)
			} else {
				fmt.Println("Please refer a new version for your application")
			}
			return
		}
		if newVersion != "" {
			if allApp {
				appPath := getAllPath()
				bumpAllGoVersion(appPath, newVersion)
			} else {
				fmt.Println("Use the all flag or specified an application name")
			}
		} else {
			fmt.Println("Please refer a new version or use a different flag")
		}
	},
}

func bumpOneGoVersion(appPath string, newVersion string) {
	path := strings.Split(appPath, "app path: ")
	os.Chdir(path[1])
	f, err := os.ReadFile("go.mod")
	if err != nil {
		log.Println(err)
	}

	lines := strings.Split(string(f), "\n")

	lines[2] = "go " + newVersion

	output := strings.Join(lines, "\n")
	err = os.WriteFile("go.mod", []byte(output), 0644)

	updateDependencies := exec.Command("go", "mod", "tidy")
	if err := updateDependencies.Run(); err != nil {
		log.Fatalf("failed to run go mod tidy: %v", err)
	} else {
		log.Println("Dependencies updated")
	}
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Go version of the application upgraded successfully.")
	}
}

func bumpAllGoVersion(appPath []string, newVersion string) {
	for i := range appPath {
		bumpOneGoVersion(appPath[i], newVersion)
	}
	log.Println("All go version upgraded successfully")
}

func getAllPath() []string {
	var outpout []string
	os.Chdir(installedPath + "/storage")
	f, err := os.ReadFile("app.txt")
	if err != nil {
		log.Println(err)
	}
	fileSplit := strings.Split(string(f), "\n")
	for _, value := range fileSplit {
		if strings.HasPrefix(value, "app path: ") {
			outpout = append(outpout, value)
		}
	}
	return outpout
}

func getAppPath(appName string) string {
	os.Chdir(installedPath + "/storage")
	f, err := os.ReadFile("app.txt")
	if err != nil {
		log.Println(err)
	}
	outpout := strings.Split(string(f), "\n")
	for i := range outpout {
		if outpout[i] == "name: "+appName {
			return outpout[i+1]
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeVersionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	upgradeCmd.PersistentFlags().String("version", "", "new version of the app")
	upgradeCmd.PersistentFlags().String("name", "", "name of the application")
	upgradeCmd.Flags().BoolP("all", "a", false, "Select all application")
}
