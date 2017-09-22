package cmd

import (
	"encoding/base64"
	"encoding/json"
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
		device := getDevice()

		if challengeFlag == "" {
			fmt.Fprintln(os.Stderr, "Please supply the challenge using -challenge option.")
			return
		}
		if appIDFlag == "" {
			fmt.Fprintln(os.Stderr, "Please supply the appID using -appid option.")
			return
		}
		if keyHandleFlag == "" {
			fmt.Fprintln(os.Stderr, "Please supply a valid Key Handle using -keyhandle option.")
			return
		}

		appIDHash := sum256(appIDFlag)
		challengeHash := sum256(challengeFlag)

		keyHandleBytes, err := base64.RawURLEncoding.DecodeString(keyHandleFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Keyhandle is not valid url encoded base64")
		}

		dev, err := u2fhid.Open(device)
		defer dev.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		t := u2ftoken.NewToken(dev)

		req := u2ftoken.AuthenticateRequest{
			Challenge:   challengeHash,
			Application: appIDHash,
			KeyHandle:   keyHandleBytes,
		}

		if err := t.CheckAuthenticate(req); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stderr, "Authenticating, press the button on your U2F device")

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

		// output for easier consumption by another program
		output := map[string]interface{}{
			"Counter":   resp.Counter,
			"Signature": base64.RawURLEncoding.EncodeToString(resp.RawResponse),
		}
		jsonOut, _ := json.MarshalIndent(output, "", "  ")

		fmt.Println(string(jsonOut))

	},
}

func init() {
	RootCmd.AddCommand(sigCmd)
}
