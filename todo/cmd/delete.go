/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteFlag int

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an item from the database",
	Long: `Delete an item from the database.
		Needs the item ID.

		Example: todo delete -d 1
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running DELETE_DB_ITEM...")
		err := ToDo.DeleteItem(deleteFlag)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Ok")
	},
}

func init() {
	dbCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().IntVar(&deleteFlag, "d", 0, "Delete an item from the database")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
