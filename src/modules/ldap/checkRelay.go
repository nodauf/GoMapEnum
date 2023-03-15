package ldap

import (
	"GoMapEnum/src/logger"
	"reflect"
)

func (options *Options) CheckRelay() {
	optionsInterface := reflect.ValueOf(options).Interface()
	options.Log.Debug("Initializing LDAP")
	if !RetrieveTargetInfo(&optionsInterface) {
		options.Log.Error("Cannot initialize LDAP")
	}

	// Dirty trick. log.fail display only if verbose, here we want to always print for fail
	if options.Log.Level <= logger.InfoLevel {
		options.Log.Level = logger.VerboseLevel
	}
	/*var credentialProvided = true
	if options.Users == "" {
		credentialProvided = false
		options.Users = "guest"
	}
	_, err := options.authenticate(options.Users, options.Passwords)

	// Check for issues with LDAPs
	if err != nil && strings.Contains(err.Error(), "connection reset by peer") {
		options.Log.Error("Something is wrong with LDAPS (maybe an issue with the certificate on the ldap server, missing certificate or whatever)")
		return
	}
	// Check for ldap signing
	if credentialProvided {
		if !ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			if options.TLS == "NoTLS" {
				if ldap.IsErrorWithCode(err, ldap.LDAPResultStrongAuthRequired) {
					options.Log.Fail("LDAP signing is required !")
				} else {
					options.Log.Success("LDAP signing is not required")
				}
			} else {
				options.Log.Debug("Using TLS, not checking LDAP signing")
			}
		} else {
			options.Log.Error("Credentials are not valid. Cannot check for LDAP signing")
		}
	} else {
		options.Log.Error("Credentials are needed to check ldap signing")
	}*/

	if options.Users != "" {
		signingRequired, err := isLDAPSigningRequired(options.Target, options.Users, options.Passwords, options.Log, options.Timeout)
		if err != nil {
			options.Log.Error(err.Error())
		} else {
			if signingRequired {
				options.Log.Fail("LDAP signing is required !")
			} else {
				options.Log.Success("LDAP signing is not required")
			}
		}
	} else {
		options.Log.Error("Credentials are needed to check ldap signing")
	}

	// Check for ldap binding
	LDAPBinding, err := isLDAPBindingEnforced(options.Target, options.Log, options.Timeout)
	if err != nil {
		options.Log.Error(err.Error())
	} else {
		if LDAPBinding {
			options.Log.Fail("Chanel binding is enforced")
		} else {
			options.Log.Success("Channel binding is not enforced")
		}
	}

}
