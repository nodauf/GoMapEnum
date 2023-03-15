package ldap

import (
	"github.com/go-ldap/ldap/v3"
)

func Authenticate(optionsInterface *interface{}, username, password string) bool {
	options := (*optionsInterface).(*Options)
	valid, err := options.authenticate(username, password)
	if err != nil && !ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
		// Handle errors
		if ldap.IsErrorWithCode(err, ldap.LDAPResultStrongAuthRequired) {
			options.Log.Error("LDAP signing is required. Please use --tls TLS")
		}
		options.Log.Error("%s - LDAP error: %v", username, err)

	}
	options.ldapConn.Close()
	return valid
}
