package cmd

import (
	"errors"
	"strconv"

	"github.com/epiefe/jswap/internal/adoptium"
	"github.com/spf13/cobra"
)

var lsRemoteCmd = &cobra.Command{
	Use:   "versions [<release>]",
	Short: "List remote versions available for install",
	Long:  "List remote versions available for install.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		release := 0
		if len(args) > 0 {
			var err error
			release, err = strconv.Atoi(args[0])
			if err != nil {
				return errors.New("release arg must be a number")
			}
		}
		if err := adoptium.PrintRemoteVersions(release); err != nil {
			return err
		}
		return nil
	},
}
