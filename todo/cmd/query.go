/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var queryFlag int

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query an item in the database",
	Long: `Query an item in the database.
	Needs the item ID.
		
	Example: 
		- todo db query -i 1
		- todo db query --id 1
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("query called")
		createTodoDb()
		item, err := ToDo.GetItem(queryFlag)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		ToDo.PrintItem(item)
		fmt.Println("Ok")
	},
}

func init() {
	dbCmd.AddCommand(queryCmd)
	queryCmd.PersistentFlags().IntVarP(&queryFlag, "id", "i", 0, "Query an item in the database")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
