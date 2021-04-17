package cmd

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mal-as/imss/router"
	"github.com/spf13/cobra"
)

var routerCmd = &cobra.Command{
	Use:   "router",
	Short: "router command",
	Run: func(cmd *cobra.Command, args []string) {
		storages := strings.Split(os.Getenv("STORAGE_ADDRESES"), ",")
		if len(storages) == 0 {
			log.Fatal("empty storages list")
		}

		err := router.Init(storages, os.Getenv("ROUTER_PORT"))
		if err != nil {
			log.Fatal(err)
		}

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		<-interrupt

	},
}

func init() {
	rootCmd.AddCommand(routerCmd)
}
