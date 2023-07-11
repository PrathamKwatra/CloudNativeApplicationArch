/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addFlag string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item to the database",
	Long: `Add an item to the database.
		Needs the JSON string of the item.

		Example: todo add -a '{"id": 1, "name":"test", "done":false}'
	`,
	Run: func(cmd *cobra.Command, args []string) {
		createTodoDb()
		fmt.Println("Running ADD_DB_ITEM...")
		item, err := ToDo.JsonToItem(addFlag)
		if err != nil {
			fmt.Println("Add option requires a valid JSON todo item string")
			fmt.Println("Error: ", err)
			return
		}
		if err := ToDo.AddItem(item); err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Ok")
	},
}

func init() {
	dbCmd.AddCommand(addCmd)
	// Set the value of the addFlag variable when the user passes a valid JSON string
	addCmd.Flags().StringVarP(&addFlag, "add", "a", "", "Add an item to the database")
}
