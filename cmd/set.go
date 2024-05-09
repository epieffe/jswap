package cmd

import (
	"strconv"

	"github.com/epiefe/jswap/internal/jdk"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:     "set {<major> | <release>}",
	Short:   "Modify PATH and JAVA_HOME to use a JDK release",
	Long:    "Modify PATH and JAVA_HOME to use a JDK release. If arg is a major number use the latest installed release.",
	Args:    cobra.MaximumNArgs(1),
	Example: "  jswap set 21\n" + "  jswap set jdk-21.0.2+13",
	RunE: func(cmd *cobra.Command, args []string) error {
		major, err := strconv.Atoi(args[0])
		if err != nil {
			// arg is a release name
			return jdk.SetRelease(args[0])
		} else {
			// arg is a major integer
			return jdk.SetMajor(major)
		}
	},
}
