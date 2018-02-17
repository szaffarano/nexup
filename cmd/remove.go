package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Borra un archivo del repositorio",
	Long:  `Borra un archivo del repositorio`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("@TODO")
	},
}
