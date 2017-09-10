package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/flynn/u2f/u2fhid"
	"github.com/flynn/u2f/u2ftoken"
	"github.com/spf13/cobra"
)

// sigCmd respresents the sig command
var sigCmd = &cobra.Command{
	Use:   "sig",
	Short: "Sign a challenge with the provided Key Handle",
	Long:  "Sign a challenge with the provided Key Handle",
	Run: func(cmd *cobra.Command, args []string) {
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
		if keyHandleFlag == "" {
			fmt.Println("Please supply a valid Key Handle using -keyhandle option.")
			return
		}

		appIDHash := sum256(appIDFlag)
		challengeHash := sum256(challengeFlag)

		keyHandleBytes, err := base64.RawURLEncoding.DecodeString(keyHandleFlag)
		if err != nil {
			fmt.Println("Keyhandle is not valid url encoded base64")
		}

		dev, err := u2fhid.Open(device)
		defer dev.Close()
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

		t := u2ftoken.NewToken(dev)

		req := u2ftoken.AuthenticateRequest{
			Challenge:   challengeHash,
			Application: appIDHash,
			KeyHandle:   keyHandleBytes,
		}

		if err := t.CheckAuthenticate(req); err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

		fmt.Println("Authenticating, press the button on your U2F device")
		var resp *u2ftoken.AuthenticateResponse
		for {
			resp, err = t.Authenticate(req)
			if err == u2ftoken.ErrPresenceRequired {
				time.Sleep(200 * time.Millisecond)
				continue
			} else if err != nil {
				log.Fatal(err)
			}
			break
		}
		fmt.Printf("\nCounter: %d\nSignature: %s\nRaw Response: %s\n",
			resp.Counter,
			base64.RawURLEncoding.EncodeToString(resp.Signature),
			base64.RawURLEncoding.EncodeToString(resp.RawResponse))
	},
}

func init() {
	RootCmd.AddCommand(sigCmd)
}
