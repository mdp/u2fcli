package cmd

import (
	"encoding/base64"
	"encoding/json"
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
		device := getDevice()

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

		fmt.Fprintf(os.Stderr, "Registering, press the button on your U2F device #%d [%s]", deviceNum, device.Product)

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

		// Parse the response
		pubKey := res[1:66]
		res = res[66:]
		khLen := int(res[0])
		res = res[1:]
		keyHandle := res[:khLen]

		// output for easier consumption by another program
		output := map[string]interface{}{
			"RegisteredData": base64.RawURLEncoding.EncodeToString(res),
			"PublicKey":      base64.RawURLEncoding.EncodeToString(pubKey),
			"KeyHandle":      base64.RawURLEncoding.EncodeToString(keyHandle),
		}
		jsonOut, _ := json.MarshalIndent(output, "", "  ")

		fmt.Println(string(jsonOut))
	},
}

func init() {
	RootCmd.AddCommand(regCmd)
}
