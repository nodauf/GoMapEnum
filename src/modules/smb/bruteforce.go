package smb

import (
	"GoMapEnum/src/utils"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/nodauf/go-smb2"
)

func RetrieveTargetInfo(optionsInterface *interface{}) bool {
	options := (*optionsInterface).(*Options)
	var err error
	var domain string
	domain, options.Hostname, err = GetTargetInfo(options.Target, options.Timeout)
	if err != nil {
		options.Log.Error("Fail to connect to smb to retrieve the domain name and hostname. Please provide the domain with -d flag.", err.Error())
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

	// Does not worth it to reuse connection. According to few tests, only two authentication can happen on the same socket. After that it is closed
	netDialer := net.Dialer{Timeout: time.Duration(time.Duration(options.Timeout) * time.Second)}

	// Open a socket on port 445
	smbConnection, err := netDialer.Dial("tcp", options.Target+":445")
	if err != nil || smbConnection == nil {
		options.Log.Error("Fail to connect to " + options.Target)
		if err != nil {
			options.Log.Error(err.Error())
		}
		return false
	} else {

		options.Log.Debug("Connection on port 445 is established")
	}

	defer (smbConnection).Close()

	var smbDialer = &smb2.Dialer{}
	if options.IsHash {
		hash, err := hex.DecodeString(strings.TrimSpace(password))
		if err != nil {
			options.Log.Error("Cannot decode the hash " + password + " from hex to byte: " + err.Error())
			return false
		}
		smbDialer = &smb2.Dialer{
			Initiator: &smb2.NTLMInitiator{
				User:   username,
				Hash:   hash,
				Domain: options.Domain,
			},
		}
	} else {
		smbDialer = &smb2.Dialer{
			Initiator: &smb2.NTLMInitiator{
				User:     username,
				Password: password,
				Domain:   options.Domain,
			},
		}

	}
	s, err := smbDialer.Dial(smbConnection)
	if err != nil {
		if strings.Contains(err.Error(), "The user account has been automatically locked because too many invalid logon attempts or password change attempts have been requested") {
			options.Log.Error("The account %s is locked", username)
			options.lockoutCounter++
		}
		// Fail safe to avoid locking to many account
		if options.lockoutCounter >= options.LockoutThreshold {
			options.Log.Fatal("Too many lockout: " + strconv.Itoa(options.lockoutCounter) + " >= " + strconv.Itoa(options.LockoutThreshold))
		}

		// If it is another error
		if !strings.Contains(err.Error(), "The attempted logon is invalid. This is either due to a bad username or authentication information.") {
			options.Log.Error(err.Error())
		}
		return false
	}
	defer s.Logoff()

	return true
}

func GetTargetInfo(target string, timeout int) (string, string, error) {
	netDialer := net.Dialer{Timeout: time.Duration(time.Duration(timeout) * time.Second)}
	// Open a socket on port 445
	smbConnection, err := netDialer.Dial("tcp", target+":445")
	if err != nil || smbConnection == nil {
		errStr := "Fail to connect to " + target
		if err != nil {
			errStr += " " + err.Error()
		}
		return "", "", fmt.Errorf(errStr)
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
