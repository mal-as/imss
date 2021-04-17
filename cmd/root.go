package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "in memory sharded storage",
		Short: "Starting application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use `imss help` for info")
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
