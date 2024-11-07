/*
Copyright © 2024 Elouan DA COSTA PEIXOTO elouandacostapeixoto@gmail.com
*/
package cmd

import (
	"os"

	"github.com/ElouanDaCosta/gam/utils"
	"github.com/spf13/cobra"
)

var installedPath = utils.GetInstalledPath()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "1.2.0",
	Use:     "gam",
	Short:   "Generate a basic template of an application in Golang",
	Long:    `Generate template for applications in Golang with the package that you want, like gin, gRPC etc.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.microservice-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
