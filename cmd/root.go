package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var keyHandle string
var appID string
var challenge string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "u2fcli",
	Short: "Interact with U2F tokens from the command line",
	Long: `u2fcli lets you interact with a hardware U2F token.
	Could be used for debugging, development and demonstrations purposes`,
	// Uncomment the following line if your bare application has an action associated with it
	//	Run: func(cmd *cobra.Command, args []string) { },
}

//Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&challenge, "challenge", "", "Challenge(string)")
	RootCmd.PersistentFlags().StringVar(&appID, "appid", "", "Applicaiton ID(string)")
	RootCmd.PersistentFlags().StringVar(&keyHandle, "keyhandle", "", "Key Handle ID(base64 string)")
}
