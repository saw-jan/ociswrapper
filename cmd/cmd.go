package cmd

import (
	"fmt"

	"ociswrapper/ocis/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ociswrapper",
	Short: "ociswrapper is a wrapper for oCIS server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func serveCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Flag("bin").Value)
			fmt.Println(cmd.Flag("url").Value)
			fmt.Println(cmd.Flag("wrapper-port").Value)

			// set configs
			config.Set("bin", cmd.Flag("bin").Value.String())
			config.Set("url", cmd.Flag("url").Value.String())
		},
	}

	// serve command args
	serveCmd.Flags().SortFlags = false
	serveCmd.Flags().StringP("bin", "", config.Get("bin"), "Full oCIS binary path")
	serveCmd.Flags().StringP("url", "", config.Get("url"), "oCIS server url")
	serveCmd.Flags().StringP("wrapper-port", "p", "5000", "Wrapper API server port")

	return serveCmd
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(serveCmd())
	rootCmd.Execute()
}
