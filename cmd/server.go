package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"context"
	"github.com/jc-m/rhub/server"
)

func init() {
	params := server.NewParams()

	runCommand := &cobra.Command{
		Use:   "server",
		Short: "Start the rhub server",
		Long:  "Start the rhub server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Server")
			ctx := context.Background()

			server := server.New(params)
			server.Run(ctx)
		},
	}

	runCommand.Flags().StringVarP(&params.ConfigPath, "config", "c", "", "config file name")
	runCommand.MarkFlagRequired("config")


	RootCmd.AddCommand(runCommand)
}