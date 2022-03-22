package ldap

import (
	"strings"
)

func Authenticate(optionsInterface *interface{}, username, password string) bool {
	options := (*optionsInterface).(*Options)
	options.Log.Verbose("Using domain " + options.Domain + " for authentication. Hostname: " + options.Hostname)

	err := options.authenticateNTLM(username, password, options.IsHash)
	defer options.ldapConn.Close()
	if err != nil {
		if !strings.Contains(err.Error(), "Invalid Credentials") {
			options.Log.Error("fail to authenticate: %s", err.Error())
		}
		return false
	}
	return true
}
