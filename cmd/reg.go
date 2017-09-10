package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/flynn/u2f/u2fhid"
	"github.com/flynn/u2f/u2ftoken"
	"github.com/spf13/cobra"
)

// regCmd respresents the reg command
var regCmd = &cobra.Command{
	Use:   "reg",
	Short: "Register with a U2F device",
	Long: `Register with a U2F device
Requires a challege and appID. For example:

u2fcli reg --challenge MyChallenge --appid https://mysite.com`,
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := u2fhid.Devices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}

		if len(devices) == 0 {
			fmt.Fprintln(os.Stderr, "Error: No devices found")
			os.Exit(1)
		}

		device := devices[0]

		if challengeFlag == "" {
			fmt.Println(os.Stderr, "Please supply the challenge using -challenge option.")
			return
		}
		if appIDFlag == "" {
			fmt.Fprintln(os.Stderr, "Please supply the appID using -appid option.")
			return
		}

		appIDHash := sum256(appIDFlag)
		challengeHash := sum256(challengeFlag)

		dev, err := u2fhid.Open(device)
		defer dev.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening device: %s\n", err)
			os.Exit(1)
		}
		t := u2ftoken.NewToken(dev)

		fmt.Println(os.Stderr, "Registering, press the button on your U2F device")

		var res []byte
		for {
			res, err = t.Register(u2ftoken.RegisterRequest{Challenge: challengeHash, Application: appIDHash})
			if err == u2ftoken.ErrPresenceRequired {
				time.Sleep(200 * time.Millisecond)
				continue
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "Error registering with device: %s\n", err)
				os.Exit(1)
			}
			break
		}

		fmt.Printf("Registered Data: %s\n", base64.RawURLEncoding.EncodeToString(res))
		pubKey := res[1:66]
		res = res[66:]
		khLen := int(res[0])
		res = res[1:]
		keyHandle := res[:khLen]
		fmt.Printf("Public Key: %s\n", base64.RawURLEncoding.EncodeToString(pubKey))
		fmt.Printf("Key Handle: %s\n", base64.RawURLEncoding.EncodeToString(keyHandle))
	},
}

func init() {
	RootCmd.AddCommand(regCmd)
}
