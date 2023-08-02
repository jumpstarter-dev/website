/*
Copyright © 2023 Miguel Angel Ajo Pelayo <majopela@redhat.com
*/
package cmd

import (
	"fmt"

	"github.com/redhat-et/jumpstarter/pkg/harness"
	"github.com/spf13/cobra"
)

// powerCmd represents the listDevices command
var attachStorage = &cobra.Command{
	Use:   "attach-storage",
	Short: "Attaches storage to the device",

	Run: func(cmd *cobra.Command, args []string) {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			handleErrorAsFatal(err)
		}

		driver := cmd.Flag("driver").Value.String()
		device, err := harness.FindDevice(driver, args[0])
		handleErrorAsFatal(err)

		fmt.Printf("💾 Attaching storage for %s ... ", args[0])
		err = device.AttachStorage(true)
		handleErrorAsFatal(err)
		fmt.Println("done")

	},
}

func init() {
	rootCmd.AddCommand(attachStorage)
	attachStorage.Flags().StringP("driver", "d", "", "Only list devices for the specified driver")
}
