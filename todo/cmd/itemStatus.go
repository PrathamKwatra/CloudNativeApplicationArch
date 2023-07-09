/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// itemStatusCmd represents the itemStatus command
var itemStatusCmd = &cobra.Command{
	Use:   "itemStatus",
	Short: "Change the status of an item in the database",
	Long: `Change the status of an item in the database. 
		Needs the item ID and the new status.
	
		Example: todo query -q 1 ItemStatus -s true
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("itemStatus args: ", args)
		fmt.Println("itemStatus called")
	},
}

func init() {
	queryCmd.AddCommand(itemStatusCmd)

	itemStatusCmd.PersistentFlags().BoolVar(&itemStatusFlag, "s", false, "Change item 'done' status to true or false")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// itemStatusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// itemStatusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
