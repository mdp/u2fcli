package cmd

import (
	"fmt"
	"os"

	"github.com/flynn/u2f/u2fhid"
	"github.com/spf13/cobra"
)

// lsCmd respresents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List u2f devices attached",
	Long:  `List of all u2f devices attached`,
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := u2fhid.Devices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening hid usb%+v", err)
			os.Exit(1)
		}

		if len(devices) == 0 {
			fmt.Fprintln(os.Stderr, "Error: No devices found")
			os.Exit(1)
		}

		for i, d := range devices {
			fmt.Printf("[%d]: manufacturer = %q, product = %q, vid = 0x%04x, pid = 0x%04x\n", i, d.Manufacturer, d.Product, d.ProductID, d.VendorID)
		}
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
