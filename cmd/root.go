package cmd

import (
	"github.com/monishth/snail/internal/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snail",
	Short: "A quick, simple, faster-than-snail forward proxy server",
	Long:  "A quick, simple, faster-than-snail forward proxy server",
	PreRun: func(cmd *cobra.Command, args []string) {
		switch serverOptions.AuthProvider {
		case server.HtpasswdAuth:
			cmd.MarkFlagRequired("filename")
		case server.SimpleAuth:
			cmd.MarkFlagRequired("userpass")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		server.RunServer(serverOptions)
	},
}

var serverOptions server.ServerOptions = server.ServerOptions{
	AuthProvider: server.NoAuth,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().VarP(&serverOptions.AuthProvider, "auth", "a", `auth provider - must be "htpasswd", "simple" or "none"`)
	rootCmd.PersistentFlags().IntVarP(&serverOptions.Port, "port", "p", 8080, `server port`)
	rootCmd.PersistentFlags().StringVarP(&serverOptions.HttpasswdFilename, "filename", "f", "", "filename for htpasswd auth ")
	rootCmd.PersistentFlags().StringVarP(&serverOptions.SimpleCredentials, "userpass", "u", "", "user:pass pair for simple ")
}
