package cmd

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sort"

	"github.com/flynn/hid"
	"github.com/flynn/u2f/u2fhid"
	"github.com/spf13/cobra"
)

var keyHandleFlag string
var appIDFlag string
var challengeFlag string
var publicKeyFlag string
var signatureFlag string
var deviceNum int

func sum256(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

// Ordered Devices by Path
type Devices []*hid.DeviceInfo

func (s Devices) Len() int {
	return len(s)
}

func (s Devices) Less(i, j int) bool {
	return s[i].Path < s[j].Path
}

func (s Devices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Devices come back in a random order, which makes it diffucult to select
// the device on a command line for each run
func getOrderedDevices() []*hid.DeviceInfo {
	devices, err := u2fhid.Devices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening hid usb: %+v", err)
		os.Exit(1)
	}

	if len(devices) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No devices found")
		os.Exit(1)
	}

	sort.Sort(Devices(devices))
	return devices
}

func getDevice() *hid.DeviceInfo {
	devices := getOrderedDevices()

	if deviceNum > len(devices) {
		fmt.Fprintf(os.Stderr, "Error: Device [%d] not found", deviceNum)
		os.Exit(1)
	}

	return devices[deviceNum-1]
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
	RootCmd.PersistentFlags().IntVar(&deviceNum, "d", 1, "Device number if multiple devices available")
}
