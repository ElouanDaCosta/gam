/*
Copyright © 2024 Elouan DA COSTA PEIXOTO elouandacostapeixoto@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List generated application.",
	Long: `List all generated application or one application by name. For example:
gam list
gam list --name [your_app_name]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		appName, _ := cmd.Flags().GetString("name")

		if appName != "" {
			getOneApp(appName)
		} else {
			listAllApp()
		}
	},
}

func listAllApp() {
	os.Chdir(installedPath + "/storage")
	f, err := os.ReadFile("app.txt")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(f))
}

func getOneApp(appName string) {
	os.Chdir(installedPath + "/storage")
	f, err := os.ReadFile("app.txt")
	if err != nil {
		log.Println(err)
	}
	outpout := strings.Split(string(f), "\n")
	for i := range outpout {
		if outpout[i] == "name: "+appName {
			fmt.Println(outpout[i])
			fmt.Println(outpout[i+1])
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().String("name", "", "Return the app with the given name")
}
