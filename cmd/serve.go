package cmd

import (
	"os"

	"github.com/meghashyamc/auth/api"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	Long:  `Start HTTP server and serve endpoints related to user authentication`,
	Run: func(cmd *cobra.Command, args []string) {
		listener, err := api.NewHTTPListener()
		if err != nil {
			os.Exit(1)
		}
		listener.Listen()
	},
}
