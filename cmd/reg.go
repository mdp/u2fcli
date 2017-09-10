package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
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
		fmt.Println("reg called")
		devices, err := u2fhid.Devices()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		device := devices[0]

		if challengeFlag == "" {
			fmt.Println("Please supply the challenge using -challenge option.")
			return
		}
		if appIDFlag == "" {
			fmt.Println("Please supply the appID using -appid option.")
			return
		}

		appIDHash := sum256(appIDFlag)
		challengeHash := sum256(challengeFlag)

		dev, err := u2fhid.Open(device)
		if err != nil {
			log.Fatal(err)
		}
		t := u2ftoken.NewToken(dev)

		var res []byte
		log.Println("registering, provide user presence")
		for {
			res, err = t.Register(u2ftoken.RegisterRequest{Challenge: challengeHash, Application: appIDHash})
			if err == u2ftoken.ErrPresenceRequired {
				time.Sleep(200 * time.Millisecond)
				continue
			} else if err != nil {
				log.Fatal(err)
			}
			break
		}
		dev.Close()

		log.Printf("registered: %s", base64.RawURLEncoding.EncodeToString(res))
		pubKey := res[1:66]
		res = res[66:]
		khLen := int(res[0])
		res = res[1:]
		keyHandle := res[:khLen]
		log.Printf("pubKey: %s", base64.RawURLEncoding.EncodeToString(pubKey))
		log.Printf("keyHandle: %s", base64.RawURLEncoding.EncodeToString(keyHandle))
	},
}

func init() {
	RootCmd.AddCommand(regCmd)
}
