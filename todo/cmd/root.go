/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"drexel.edu/todo/db"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "todo application",
	Long: `A todo application written in Go for CLI.
		The application uses a local store as a database.
		You can add, list, delete, update and complete tasks.
		A task looks like:
		{
			"id": 1,
			"title": "task title",
			"completed": false,
		}
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var (
	dbFileNameFlag string
	listFlag       bool
	itemStatusFlag bool
	queryFlag      int
	addFlag        string
	updateFlag     string
	deleteFlag     int
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func CreateNewDb() *db.ToDo {
	todo, err := db.New(dbFileNameFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return todo
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
