package cmd

import (
	"errors"
	"strconv"

	"github.com/epiefe/jswap/internal/jdk/adoptium"
	"github.com/spf13/cobra"
)

var releasesCmd = &cobra.Command{
	Use:     "releases [<major>]",
	Short:   "List remote releases available for install",
	Long:    "List remote releases available for install, matching a given major if provided.",
	Args:    cobra.MaximumNArgs(1),
	Example: "  jswap releases\n" + "  jswap releases 21",
	RunE: func(cmd *cobra.Command, args []string) error {
		release := 0
		if len(args) > 0 {
			var err error
			release, err = strconv.Atoi(args[0])
			if err != nil {
				return errors.New("arg must be an integer")
			}
		}
		if err := adoptium.PrintRemoteReleases(release); err != nil {
			return err
		}
		return nil
	},
}
