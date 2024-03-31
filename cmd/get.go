package cmd

import (
	"strconv"

	"github.com/epiefe/jswap/internal/adoptium"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <version>",
	Short: "Download and install a JDK",
	Long: "Download and install a JDK. If <version> is a release number (e.g., 21) downloads the latest available version for that release, " +
		"otherwise <version> must be a specific version name. Use 'jswap versions' to see all the available version names.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		release, err := strconv.Atoi(args[0])
		if err != nil {
			// <version> is a specific version name
			if err := adoptium.DownloadVersion(args[0]); err != nil {
				return err
			}
		} else {
			// <version> is a release number
			if err := adoptium.DownloadLatest(release); err != nil {
				return err
			}
		}
		return nil
	},
}
