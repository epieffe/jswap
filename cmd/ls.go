package cmd

import (
	"errors"
	"strconv"

	"github.com/epiefe/jswap/internal/jdk/adoptium"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls [<major>]",
	Short:   "List installed releases",
	Long:    "List installed releases, matching a given major if provided.",
	Args:    cobra.MaximumNArgs(1),
	Example: "  jswap ls\n" + "  jswap ls 21",
	RunE: func(cmd *cobra.Command, args []string) error {
		major := 0
		if len(args) > 0 {
			var err error
			major, err = strconv.Atoi(args[0])
			if err != nil {
				return errors.New("arg must be an integer")
			}
		}
		if err := adoptium.PrintLocalReleases(major); err != nil {
			return err
		}
		return nil
	},
}
