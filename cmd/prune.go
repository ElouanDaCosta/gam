/*
Copyright © 2024 Elouan DA COSTA PEIXOTO elouandacostapeixoto@gmail.com
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove an application from the storage",
	Long: `Remove an application from the storage,
clear completely the storage where all the application is saved 
or remove an application with the remove-app flag. For example:

gam prune --name new_app
gam prune --all
gam prune --app --name new_app
`,
	Run: func(cmd *cobra.Command, args []string) {
		appName, _ := cmd.Flags().GetString("name")
		clearAllApp, _ := cmd.Flags().GetBool("all")
		removeApp, _ := cmd.Flags().GetBool("app")

		if appName != "" {
			if removeApp {
				deleteApp(appName)
			} else {
				removeFromStorage(appName)
			}
			return
		}
		if removeApp {
			log.Println("Please refer an application name to remove.")
		} else if !clearAllApp {
			log.Println("Please refer an application name.")
		} else {
			cleanAllStorage()
		}
	},
}

func removeFromStorage(appName string) {
	os.Chdir(installedPath + "/storage")

	filePath := "app.txt"
	targetLine := appName

	// read the app storage file
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	lines := strings.Split(string(content), "\n")

	// this code snippet is iterating over each line in the `lines` slice, which contains the content of a
	// file read earlier. For each line, it checks if the line matches the pattern `"name: "+targetLine`.
	// If a match is found, it assigns the index of that line to the `index` variable and breaks out of
	// the loop.
	// start the index at -1 to cover all the file.
	index := -1
	for i, line := range lines {
		if line == "name: "+targetLine {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Line not found")
		return
	}

	// this block of code is responsible for creating a new slice of strings called `newLines` that will
	// contain the updated content of the file after removing a specific line
	var newLines []string
	if index < len(lines)-1 {
		newLines = append(lines[:index], lines[index+2:]...)
	} else {
		newLines = lines[:index]
	}

	// this block of code is responsible for updating the content of the file after removing a specific
	// line.
	output := strings.Join(newLines, "\n")
	err = os.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		log.Fatalf("unable to write file: %v", err)
	}

	fmt.Println("Lines deleted successfully")
}

func cleanAllStorage() {
	os.Chdir(installedPath + "/storage")

	f, err := os.OpenFile("app.txt", os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	var bs []byte
	buf := bytes.NewBuffer(bs)

	for scanner.Scan() {
		_, err := buf.WriteString("")
		if err != nil {
			panic("Couldn't replace line")
		}
	}
	log.Println("Storage clean successfully")
	f.Truncate(0)
	f.Seek(0, 0)
	buf.WriteTo(f)
}

func deleteApp(appName string) {
	os.Chdir(".")
	fmt.Println("Are you sure to delete the app?(y/n)")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	if char == 'y' {
		removeContentDirectory(appName)
		removeFromStorage(appName)
		os.Remove(appName)
		if err != nil {
			fmt.Println(err)
		} else {
			log.Println("Application deleted successfully")
		}
	} else if char == 'n' {
		fmt.Println("Action canceled")
	} else {
		deleteApp(appName)
	}
}

func removeContentDirectory(appName string) {
	os.Chdir(".")
	exist, err := directoryExist(appName)
	if err != nil {
		fmt.Println(err)
	}
	if !exist {
		fmt.Println(appName + " no such file or directory")
		return
	}
	remove := os.RemoveAll(appName)
	if remove != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Content in", appName, "removed successfully")
	}
}

func directoryExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	removeCmd.PersistentFlags().String("name", "", "Clear the given application from the saved application storage")
	removeCmd.Flags().BoolP("all", "a", false, "Clear all the storage from the saved application.")
	removeCmd.Flags().BoolP("app", "r", false, "Delete the working directory of the given application. Not reversible !")
}
