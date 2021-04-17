package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mal-as/imss/storage"
	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "storage command",
	Run: func(cmd *cobra.Command, args []string) {
		err := storage.Init(os.Getenv("STORAGE_PORT"))
		if err != nil {
			log.Fatal(err)
		}

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		<-interrupt

	},
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
