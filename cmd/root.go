package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "auth",
	Short: "Sample user authentication code",
	Long:  `Backend sample code for user authentication`,
}

/*Execute sets up the rootCmd
 */
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	setupMigrate()
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(serveCmd)

}
