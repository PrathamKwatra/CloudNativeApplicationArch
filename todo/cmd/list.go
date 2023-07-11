/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the items in the database",
	Long: `List all the items in the database.
	Example: todo db list
	`,
	Run: func(cmd *cobra.Command, args []string) {
		createTodoDb()
		fmt.Println("list called")
		todoList, err := ToDo.GetAllItems()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		ToDo.PrintAllItems(todoList)
		fmt.Println("THERE ARE", len(todoList), "ITEMS IN THE DB")
		fmt.Println("Ok")
	},
}

func init() {
	dbCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
