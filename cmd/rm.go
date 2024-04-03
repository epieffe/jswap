package cmd

import (
	"github.com/epiefe/jswap/internal/jdk"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm <release>...",
	Short:   "Remove installed releases",
	Long:    "Remove installed releases.",
	Args:    cobra.MinimumNArgs(1),
	Example: "  jswap rm jdk-21.0.2+13\n" + "  jswap rm jdk-17.0.10+7 jdk-21.0.2+13",
	RunE: func(cmd *cobra.Command, args []string) error {
		return jdk.RemoveReleases(args...)
	},
}
