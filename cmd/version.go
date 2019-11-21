package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	Version   = ""
	Commit    = ""
	BuildDate = ""
	GoArch    = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
	GoVersion = runtime.Version()
)


var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of jxcore",
	Long:  `All software has versions. This jxcore-backend`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version   : ", Version)
		fmt.Println("Commit    : ", Commit)
		fmt.Println("BuildDate : ", BuildDate)
		fmt.Println("GoArch    : ", GoArch)
		fmt.Println("GoVersion : ", GoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
