package smb

import (
	"GoMapEnum/src/utils"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/nodauf/go-smb2"
)

func (options *Options) authenticate(username, password string) (bool, error) {
	smbConnection, err := utils.OpenConnectionWoProxy(options.Target, "445", options.Timeout, options.ProxyTCP)
	if err != nil {
		return false, fmt.Errorf("cannot open a connection to %s:%d : %v", options.Target, 445, err)
	}

	options.Log.Debug("Connection on port 445 is established")

	defer smbConnection.Close()

	var smbDialer = &smb2.Dialer{}
	if options.IsHash {
		hash, err := hex.DecodeString(strings.TrimSpace(password))
		if err != nil {
			return false, fmt.Errorf("Cannot decode the hash %s from hex to byte: %v", password, err)
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
		//if strings.Contains(err.Error(), "The user account has been automatically locked because too many invalid logon attempts or password change attempts have been requested") {
		//	return false, err
		//}

		// Fail safe to avoid locking to many account
		//if options.lockoutCounter >= options.LockoutThreshold {
		//	options.Log.Fatal("Too many lockout: " + strconv.Itoa(options.lockoutCounter) + " >= " + strconv.Itoa(options.LockoutThreshold))
		//}

		// If it is another error
		//if !strings.Contains(err.Error(), "The attempted logon is invalid. This is either due to a bad username or authentication information.") {
		//	options.Log.Error(username + " " + err.Error())
		//	return false, err
		//}
		return false, err
	}
	defer s.Logoff()

	return true, nil
}
