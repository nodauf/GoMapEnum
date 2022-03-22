package smb

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/smb"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var level logger.Level
var verbose bool
var debug bool
var validUsers string
var output string
var proxyString string

var smbOptions smb.Options

// SmbCmd represents the smb command
var SmbCmd = &cobra.Command{
	Use:   "smb",
	Short: "Commands for smb module",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(validUsers), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	SmbCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	SmbCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	SmbCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")

	SmbCmd.AddCommand(bruteCmd)
}

func initLogger() {
	if debug {
		level = logger.DebugLevel
	} else if verbose {
		level = logger.VerboseLevel
	} else {
		level = logger.InfoLevel
	}

}

func initProxy() {
	if proxyString != "" {

	}
}
