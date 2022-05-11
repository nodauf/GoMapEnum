package kerberos

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/ldap"
	"strings"

	kconfig "github.com/nodauf/gokrb5/v8/config"
)

// InitSession initializes the kerberos client
func InitSession(optionsInterface *interface{}) bool {
	var err error
	options := (*optionsInterface).(*Options)
	// If the domain is not specified we try to get it with a LDAP request for the default naming context
	if options.Domain == "" {
		//options.Domain, _, err = smb.GetTargetInfo(options.Target, options.Timeout, options.ProxyTCP)
		var optionsLDAP ldap.Options
		tmpLog := new(logger.Logger)
		*tmpLog = *options.Log
		optionsLDAP.Log = tmpLog
		optionsLDAP.Log.Level = options.Log.Level
		optionsLDAP.Log.Module = "LDAP\t"
		optionsLDAP.TLS = "NoTLS"
		optionsLDAP.Target = options.Target
		optionsLDAP.ProxyHTTP = options.ProxyHTTP
		err := optionsLDAP.GetDefaultNamingContext()
		if err != nil {
			options.Log.Error("Fail to connect to ldap to retrieve the default naming context to get the domain name in netbios format: %v. Please provide the domain with -d flag.", err.Error())
			return false
		}
		options.Domain = strings.Replace(strings.ReplaceAll(optionsLDAP.BaseDN, "DC=", ""), ",", ".", -1)
	}
	options.Domain = strings.ToUpper(options.Domain)
	configstring := buildKrb5Template(options.Domain, options.Target)
	options.kerberosConfig, err = kconfig.NewFromString(configstring)
	if err != nil {
		panic(err)
	}

	_, options.kdcs, err = options.kerberosConfig.GetKDCs(options.Domain, false)
	if err != nil {
		options.Log.Error("Couldn't find any KDCs for realm %s (%v). Please specify a Domain Controller", options.Domain, err)
		return false
	}
	options.Log.Debug("Found %d KDCs for realm %s: %v", len(options.kdcs), options.Domain, options.kdcs)

	return true
}
