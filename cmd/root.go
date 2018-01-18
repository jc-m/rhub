package cmd

import (
	"path"
	"os"
	"github.com/spf13/cobra"
	"fmt"
)

var RootCmd = &cobra.Command{
	Use:   path.Base(os.Args[0]),
	Short: "Radio Hub",
	Long:  `Radio Hub`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rhub")

	},
}

