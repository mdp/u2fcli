package cmd

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var keyHandleFlag string
var appIDFlag string
var challengeFlag string
var publicKeyFlag string
var signatureFlag string

func sum256(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

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
	RootCmd.PersistentFlags().StringVar(&challengeFlag, "challenge", "", "Challenge(string)")
	RootCmd.PersistentFlags().StringVar(&appIDFlag, "appid", "", "Applicaiton ID(string)")
	RootCmd.PersistentFlags().StringVar(&keyHandleFlag, "keyhandle", "", "Key Handle ID(base64 string)")
	RootCmd.PersistentFlags().StringVar(&publicKeyFlag, "publickey", "", "Public Key of the signer")
	RootCmd.PersistentFlags().StringVar(&signatureFlag, "signature", "", "Raw Signature from U2F 'sig'")
}
