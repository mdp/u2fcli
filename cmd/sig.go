package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
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
		fmt.Println("sig called")
		devices, err := u2fhid.Devices()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		device := devices[0]

		if challenge == "" {
			fmt.Println("Please supply the challenge using -challenge option.")
			return
		}
		if appID == "" {
			fmt.Println("Please supply the appID using -appid option.")
			return
		}
		if keyHandle == "" {
			fmt.Println("Please supply a valid Key Handle using -keyhandle option.")
			return
		}
		h := sha256.New()
		h.Write([]byte(appID))
		appIDHash := h.Sum(nil)

		h = sha256.New()
		h.Write([]byte(challenge))
		challengeHash := h.Sum(nil)

		keyHandleBytes, err := base64.RawURLEncoding.DecodeString(keyHandle)
		if err != nil {
			fmt.Println("Keyhandle is not valid url encoded base64")
		}

		dev, err := u2fhid.Open(device)
		defer dev.Close()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("signing, provide user presence")
		t := u2ftoken.NewToken(dev)

		req := u2ftoken.AuthenticateRequest{
			Challenge:   challengeHash,
			Application: appIDHash,
			KeyHandle:   keyHandleBytes,
		}

		if err := t.CheckAuthenticate(req); err != nil {
			log.Fatal(err)
		}

		log.Println("authenticating, provide user presence")
		var resp *u2ftoken.AuthenticateResponse
		for {
			resp, err = t.Authenticate(req)
			if err == u2ftoken.ErrPresenceRequired {
				time.Sleep(200 * time.Millisecond)
				continue
			} else if err != nil {
				log.Fatal(err)
			}
			log.Printf("counter = %d, signature = %x\n", resp.Counter, resp.Signature)
			break
		}
		log.Printf("counter = %d, signature = %x\n", resp.Counter, resp.Signature)
	},
}

func init() {
	RootCmd.AddCommand(sigCmd)
}
