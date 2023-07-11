/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var itemStatusFlag bool

// itemStatusCmd represents the itemStatus command
var itemStatusCmd = &cobra.Command{
	Use:   "itemStatus",
	Short: "Change the status of an item in the database",
	Long: `Change the status of an item in the database. 
		Needs the item ID and the new status.
	
		Example: todo query -q 1 ItemStatus -s true
	`,
	Run: func(cmd *cobra.Command, args []string) {
		createTodoDb()
		fmt.Println("Running CHANGE_ITEM_STATUS...")
		err := ToDo.ChangeItemDoneStatus(queryFlag, itemStatusFlag)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Ok")
	},
}

func init() {
	queryCmd.AddCommand(itemStatusCmd)

	itemStatusCmd.Flags().BoolVarP(&itemStatusFlag, "status", "s", true, "Change item 'done' status to true or false")
}
