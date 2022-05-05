package ldap

func Authenticate(optionsInterface *interface{}, username, password string) bool {
	options := (*optionsInterface).(*Options)
	valid, err := options.authenticate(username, password)
	if err != nil {
		options.Log.Error("LDAP error: %s", err)
	}
	options.ldapConn.Close()
	return valid
}
