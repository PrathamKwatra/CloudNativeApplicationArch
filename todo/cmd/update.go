/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateFlag string

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an item in the database",
	Long: `Update an item in the database.
		Needs the JSON string of the item.
		
		Example: todo update -u '{"id":1,"name":"test","done":false}'
	`,
	Run: func(cmd *cobra.Command, args []string) {
		createTodoDb()
		fmt.Println("Running UPDATE_DB_ITEM...")
		item, err := ToDo.JsonToItem(updateFlag)
		if err != nil {
			fmt.Println("Update option requires a valid JSON todo item string")
			fmt.Println("Error: ", err)
			return
		}
		if err := ToDo.UpdateItem(item); err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Ok")
	},
}

func init() {
	dbCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&updateFlag, "u", "", "Update an item in the database")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
