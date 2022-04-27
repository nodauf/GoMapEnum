package smb

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/smb"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
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
	SmbCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Socks5 proxy to use (ex: localhost:8080)")

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
		var err error
		defaultDailer := &net.Dialer{Timeout: time.Duration(smbOptions.Timeout * int(time.Second))}
		smbOptions.ProxyTCP, err = proxy.SOCKS5("tcp", proxyString, nil, defaultDailer)
		if err != nil {
			fmt.Println("fail to use the proxy " + proxyString + ": " + err.Error())
			os.Exit(1)
		}
	}
}
