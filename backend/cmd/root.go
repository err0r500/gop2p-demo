package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

const (
	serverModeKey = "server"
	apiPortKey    = "api_port"
	p2pPortKey    = "p2p_port"
)

var rootCmd = &cobra.Command{
	Use:   "gop2p",
	Short: "a p2p messaging system",
	Long:  `gop2p is a peer-to-peer messaging system`,
	Run: func(cmd *cobra.Command, args []string) {

		// when every flag / env var is parsed, we start the app
		// in server or client mode according to the "server" flag
		if viper.GetBool(serverModeKey) {
			startInServerMode(viper.GetInt(apiPortKey))
		} else {
			startInClientMode(
				viper.GetInt(apiPortKey),
				viper.GetInt(p2pPortKey),
			)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// we set here the available flags
func init() {
	viper.AutomaticEnv()

	// we select if the app runs in client or server mode, defaults to client mode
	rootCmd.Flags().BoolP(serverModeKey, "s", false, "Run gop2p in server mode, default: false")
	viper.BindPFlag(serverModeKey, rootCmd.Flags().Lookup(serverModeKey))

	// we select the port on which the API server will listen on, defaults to 3000
	rootCmd.Flags().IntP(apiPortKey, "p", 3000, "The API port to listen on")
	viper.BindPFlag(apiPortKey, rootCmd.Flags().Lookup(apiPortKey))

	// we select the port on which the P2P server will listen on, defaults to 4000
	rootCmd.Flags().Int(p2pPortKey, 4000, "The port used for p2p communication between clients")
	viper.BindPFlag(p2pPortKey, rootCmd.Flags().Lookup(p2pPortKey))
}
