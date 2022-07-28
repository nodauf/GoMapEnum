package ldap

// Credits: https://github.com/ropnop/go-windapsearch

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/ldap"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	ldaplib "github.com/go-ldap/ldap/v3"
	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
)

var level logger.Level
var verbose bool
var debug bool
var outputCmd string
var output string
var proxyString string

var ldapOptions ldap.Options

// LdapCmd represents the owa command
var LdapCmd = &cobra.Command{
	Use:   "ldap",
	Short: "Commands for ldap module",
	Long: `LDAP servers are an important part of the internal network.
The authentication can be bruteforce and with valid credentials a lot of data could be retrieved and parse`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//if ldapOptions.Passwords == "" && ldapOptions.Hash == "" {
		//	return errors.New("The field password or hash is required")
		if ldapOptions.Passwords != "" && ldapOptions.Hash != "" {
			return errors.New("Only password or hash field should be specified")
		} else if ldapOptions.Hash != "" && !ldapOptions.UseNTLM {
			return errors.New("Hash can only be used with NTLM")
		}

		if ldapOptions.Domain == "" && ldapOptions.Target == "" {
			return errors.New("The field domain or target are required")
		}
		if !strings.EqualFold(ldapOptions.TLS, "NoTLS") && !strings.EqualFold(ldapOptions.TLS, "StartTLS") && !strings.EqualFold(ldapOptions.TLS, "TLS") {
			return errors.New("The field TLS should be NoTLS, TLS or StartTLS")
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(outputCmd), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	LdapCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	LdapCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	LdapCmd.PersistentFlags().BoolVar(&ldapOptions.UseNTLM, "NTLM", true, "Authenticate with NTLM")
	LdapCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The output file")
	LdapCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Socks5 proxy to use (ex: localhost:8080)")
	LdapCmd.PersistentFlags().StringVar(&ldapOptions.TLS, "tls", "NoTLS", "TLS mode (NoTLS, TLS, StartTLS))")
	LdapCmd.PersistentFlags().StringVarP(&ldapOptions.Users, "user", "u", "", "User or file containing the emails")
	LdapCmd.PersistentFlags().StringVarP(&ldapOptions.Passwords, "password", "p", "", "Password or file containing the passwords")
	LdapCmd.PersistentFlags().StringVarP(&ldapOptions.Hash, "hash", "H", "", "Hash or file containing the hashes")
	LdapCmd.PersistentFlags().StringVarP(&ldapOptions.Target, "target", "t", "", "Host pointing to the LDAP server")
	LdapCmd.PersistentFlags().StringVarP(&ldapOptions.Domain, "domain", "d", "", "Domain for the authentication (by default the domain name will be guessed with a smb connection)")
	LdapCmd.PersistentFlags().IntVar(&ldapOptions.Timeout, "timeout", int(ldaplib.DefaultTimeout.Seconds()), "Timeout for the LDAP connection in seconds")

	LdapCmd.MarkFlagRequired("user")

	LdapCmd.AddCommand(bruteCmd)
	LdapCmd.AddCommand(dumpCmd)
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
		defaultDailer := &net.Dialer{Timeout: time.Duration(ldapOptions.Timeout * int(time.Second))}
		ldapOptions.ProxyTCP, err = proxy.SOCKS5("tcp", proxyString, nil, defaultDailer)
		if err != nil {
			fmt.Println("fail to use the proxy " + proxyString + ": " + err.Error())
			os.Exit(1)
		}
	}
}
