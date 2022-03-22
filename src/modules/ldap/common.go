package ldap

import "GoMapEnum/src/modules/smb"

func RetrieveTargetInfo(optionsInterface *interface{}) bool {
	options := (*optionsInterface).(*Options)
	var err error
	if options.Domain == "" {
		options.Domain, options.Hostname, err = smb.GetTargetInfo(options.Target, options.Timeout)
		if err != nil {
			options.Log.Error("Fail to connect to smb to retrieve the domain name: %s. Please provide the domain with -d flag.", err.Error())
			return false
		}
	}
	return true
}
