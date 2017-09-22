package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd respresents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List u2f devices attached",
	Long:  `List of all u2f devices attached`,
	Run: func(cmd *cobra.Command, args []string) {

		devices := getOrderedDevices()

		for i, d := range devices {
			fmt.Printf("[%d]: manufacturer = %q, product = %q, vid = 0x%04x, pid = 0x%04x\n", i+1, d.Manufacturer, d.Product, d.ProductID, d.VendorID)
		}
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
