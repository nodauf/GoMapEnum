package smtp

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/smtp"
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

var smtpOptions smtp.Options

// SMTPCmd represents the owa command
var SMTPCmd = &cobra.Command{
	Use:   "smtp",
	Short: "Commands for owa module",
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
	SMTPCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	SMTPCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	SMTPCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	SMTPCmd.PersistentFlags().IntVar(&smtpOptions.Timeout, "timeout", 2, "Timeout for the SMTP connection in seconds")
	SMTPCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Socks5 proxy to use (ex: localhost:8080)")

	SMTPCmd.AddCommand(enumCmd)
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
		defaultDailer := &net.Dialer{Timeout: time.Duration(smtpOptions.Timeout * int(time.Second))}
		smtpOptions.ProxyTCP, err = proxy.SOCKS5("tcp", proxyString, nil, defaultDailer)
		if err != nil {
			fmt.Println("fail to use the proxy " + proxyString + ": " + err.Error())
			os.Exit(1)
		}
	}
}
