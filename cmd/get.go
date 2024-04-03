package cmd

import (
	"strconv"

	"github.com/epiefe/jswap/internal/jdk"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get {<major> | <release>}",
	Short:   "Download and install a JDK release",
	Long:    "Download and install a JDK release. If arg is a major number search for the latest available release.",
	Args:    cobra.ExactArgs(1),
	Example: "  jswap get 21\n" + "  jswap get jdk-21.0.2+13",
	RunE: func(cmd *cobra.Command, args []string) error {
		major, err := strconv.Atoi(args[0])
		if err != nil {
			// arg is a release name
			return jdk.GetRelease(args[0])
		} else {
			// arg is a major integer
			return jdk.GetLatest(major)
		}
	},
}
