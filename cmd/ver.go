package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/spf13/cobra"
)

// Following code pulled from Timothy Stranex's U2F library
// https://github.com/tstranex/u2f

type ecdsaSig struct {
	R, S *big.Int
}

type authResp struct {
	UserPresenceVerified bool
	Counter              uint32
	sig                  ecdsaSig
	raw                  []byte
}

func parseSignResponse(sd []byte) (*authResp, error) {
	if len(sd) < 5 {
		return nil, errors.New("u2f: data is too short")
	}

	var ar authResp

	userPresence := sd[0]
	if userPresence|1 != 1 {
		return nil, errors.New("u2f: invalid user presence byte")
	}
	ar.UserPresenceVerified = userPresence == 1

	ar.Counter = uint32(sd[1])<<24 | uint32(sd[2])<<16 | uint32(sd[3])<<8 | uint32(sd[4])

	ar.raw = sd[:5]

	rest, err := asn1.Unmarshal(sd[5:], &ar.sig)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("u2f: trailing data")
	}

	return &ar, nil
}

func verifyAuthSignature(ar authResp, pubKey *ecdsa.PublicKey, appID string, clientData []byte) error {
	appParam := sha256.Sum256([]byte(appID))
	challenge := sha256.Sum256(clientData)

	var buf []byte
	buf = append(buf, appParam[:]...)
	buf = append(buf, ar.raw...)
	buf = append(buf, challenge[:]...)
	hash := sha256.Sum256(buf)

	if !ecdsa.Verify(pubKey, hash[:], ar.sig.R, ar.sig.S) {
		return errors.New("u2f: invalid signature")
	}

	return nil
}

func verify(appID, challenge string, signature, publicKey []byte) error {
	ar, err := parseSignResponse(signature)
	if err != nil {
		return fmt.Errorf("Error parsing signature: %s\n", err)
	}

	x, y := elliptic.Unmarshal(elliptic.P256(), publicKey)
	if x == nil {
		return fmt.Errorf("Error unmarshalling public key")
	}

	pubKey := &ecdsa.PublicKey{}
	pubKey.Curve = elliptic.P256()
	pubKey.X = x
	pubKey.Y = y

	if err := verifyAuthSignature(*ar, pubKey, appID, []byte(challenge)); err != nil {
		return fmt.Errorf("Signature did not verify")
	}

	return nil
}

// verCmd respresents the verify command
var verCmd = &cobra.Command{
	Use:   "ver",
	Short: "Verify a signed response with the provided Public Key and Challenge",
	Long:  "Verify a signed response with the provided Public Key and Challenge",
	Run: func(cmd *cobra.Command, args []string) {
		publicKeyBytes, err := base64.RawURLEncoding.DecodeString(publicKeyFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Public Key is not valid url encoded base64")
			os.Exit(1)
		}

		signatureBytes, err := base64.RawURLEncoding.DecodeString(signatureFlag)
		if err != nil {
			fmt.Println("Signature is not valid url encoded base64")
			os.Exit(1)
		}

		if err := verify(appIDFlag, challengeFlag, signatureBytes, publicKeyBytes); err != nil {
			fmt.Fprintf(os.Stderr, "Signature did not verify: %s", err)
			os.Exit(1)
		}

		fmt.Println("Signature verified")
	},
}

func init() {
	RootCmd.AddCommand(verCmd)
}
