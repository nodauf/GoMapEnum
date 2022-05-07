package ldap

import "github.com/go-ldap/ldap/v3"

func Authenticate(optionsInterface *interface{}, username, password string) bool {
	options := (*optionsInterface).(*Options)
	valid, err := options.authenticate(username, password)
	if err != nil && !ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
		options.Log.Error("LDAP error: %s", err)
	}
	options.ldapConn.Close()
	return valid
}
