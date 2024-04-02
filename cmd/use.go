package cmd

import (
	"strconv"

	"github.com/epiefe/jswap/internal/jdk"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:     "use {<major> | <release>}",
	Short:   "Modify PATH and JAVA_HOME to use a JDK release",
	Long:    "Modify PATH and JAVA_HOME to use a JDK release. If arg is a major number use the latest installed release.",
	Args:    cobra.MaximumNArgs(1),
	Example: "  jswap use 21\n" + "  jswap use jdk-21.0.2+13",
	RunE: func(cmd *cobra.Command, args []string) error {
		major, err := strconv.Atoi(args[0])
		if err != nil {
			// arg is a release name
			if err := jdk.UseRelease(args[0]); err != nil {
				return err
			}
		} else {
			// arg is a major integer
			if err := jdk.UseMajor(major); err != nil {
				return err
			}
		}
		return nil
	},
}
