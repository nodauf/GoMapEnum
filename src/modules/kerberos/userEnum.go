package kerberos

import (
	"strings"

	kconfig "github.com/nodauf/gokrb5/v8/config"
)

func KerberosSession(optionsInterface *interface{}) bool {
	var err error
	options := (*optionsInterface).(*Options)
	options.Domain = strings.ToUpper(options.Domain)
	configstring := buildKrb5Template(options.Domain, options.DomainController)
	options.kerberosConfig, err = kconfig.NewFromString(configstring)
	if err != nil {
		panic(err)
	}
	_, options.kdcs, err = options.kerberosConfig.GetKDCs(options.Domain, false)
	if err != nil {
		options.Log.Error("Couldn't find any KDCs for realm %s (%v). Please specify a Domain Controller", options.Domain, err)
		return false
	}
	return true
}

func UserEnum(optionsInterface *interface{}, username string) bool {
	valid := false
	options := (*optionsInterface).(*Options)
	valid, err := options.TestUsername(username)
	if valid {
		options.Log.Success(username)
		valid = true
	} else if err != nil {
		// This is to determine if the error is "okay" or if we should abort everything
		ok, errorString := handleKerbError(err)
		if ok {
			options.Log.Debug("%s %s", username, errorString)
		} else {
			options.Log.Fatal("%s %s", username, errorString)
		}
	} else {
		options.Log.Debug("Unknown behavior for username %s", username)
	}

	return valid
}
