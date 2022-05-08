package smb

import (
	"GoMapEnum/src/utils"
	"fmt"
	"strings"

	"github.com/nodauf/go-smb2"
	"golang.org/x/net/proxy"
)

func RetrieveTargetInfo(optionsInterface *interface{}) bool {
	options := (*optionsInterface).(*Options)
	var err error
	var domain string
	domain, options.Hostname, err = GetTargetInfo(options.Target, options.Timeout, options.ProxyTCP)
	if err != nil {
		options.Log.Error("Fail to connect to smb to retrieve the domain name and hostname. Please provide the domain with -d flag. %s", err.Error())
		return false
	}
	if options.Domain == "" {
		options.Domain = domain
	}

	options.Log.Verbose("Using domain " + options.Domain + " for authentication. Hostname: " + options.Hostname)
	return true
}

func Authenticate(optionsInterface *interface{}, username, password string) bool {
	options := (*optionsInterface).(*Options)
	valid, err := options.authenticate(username, password)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "The user account has been automatically locked because too many invalid logon attempts or password change attempts have been requested"):
			options.Log.Error("%s has been locked out", username)
		case strings.Contains(err.Error(), "The attempted logon is invalid. This is either due to a bad username or authentication information."):
			// do nothing
		case options.StopOnLockout && strings.Contains(err.Error(), "The user account has been automatically locked because too many invalid logon attempts or password change attempts have been requested"):
			options.Log.Fatal("%s has been locked out", username)
		case strings.Contains(err.Error(), "The user account has expired"):
			valid = true
			options.Log.Verbose("The password %s of %s has expired", password, username)
		default:
			options.Log.Error(username + " " + err.Error())

		}
	}
	return valid
}

func GetTargetInfo(target string, timeout int, proxyTCP proxy.Dialer) (string, string, error) {
	smbConnection, err := utils.OpenConnectionWoProxy(target, "445", timeout, proxyTCP)
	if err != nil {
		return "", "", fmt.Errorf("cannot open a connection to %s:%d : %v", target, 445, err)
	}
	defer (smbConnection).Close()
	initiator := &smb2.NTLMInitiator{
		User:     utils.RandomString(6),
		Password: utils.RandomString(6),
	}
	smbDialer := &smb2.Dialer{
		Initiator: initiator,
	}
	s, err := smbDialer.Dial(smbConnection)
	s.Logoff()
	return initiator.TargetInfo().DomainName, initiator.TargetInfo().ServerName, nil

}
