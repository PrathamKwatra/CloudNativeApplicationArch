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

var (
	itemStatusFlag bool
	queryFlag      int
	addFlag        string
	updateFlag     string
	deleteFlag     int
)

var ToDo *db.ToDo

//

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Change the database",
	Long: `Change the database. 
	Defaults to ./data/todo.json
	Has to be a json file.
	Can be relative or absolute path. 

	Example: todo db -f ./data/todo.json
	Example: todo db -f ./data/todo.json add -a '{"id": 1, "name":"test", "done":false}'
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("db is set to: " + dbFileNameFlag)
		createTodoDb()
	},
}

func createTodoDb() {
	todo, err := db.New(dbFileNameFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ToDo = todo
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.Flags().StringVarP(&dbFileNameFlag, "file", "f", "./data/todo.json", "Name of the database file")
	createTodoDb()
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
