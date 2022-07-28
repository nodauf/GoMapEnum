package ldap

import (
	"GoMapEnum/src/modules/smb"
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

func RetrieveTargetInfo(optionsInterface *interface{}) bool {
	options := (*optionsInterface).(*Options)
	var err error
	if options.Domain == "" {
		options.Log.Debug("Domain is not specified, trying to retrieve it from the target through smb")
		options.Domain, options.Hostname, err = smb.GetTargetInfo(options.Target, options.Timeout, options.ProxyTCP)
		if err != nil {
			options.Log.Error("Fail to connect to smb to retrieve the domain name: %s. Please provide the domain with -d flag.", err.Error())
			return false
		}
	} else if options.Target == "" {
		options.Log.Debug("Target is not specified, trying to retrieve it from the target through DNS")
		ldapServers, _ := findLDAPServers(options.Domain)
		options.Target = ldapServers[0]
	}
	return true
}

func (options *Options) GetDefaultNamingContext() error {
	if options.BaseDN != "" {
		return nil
	}
	if options.ldapConn == nil {
		options.ldapConn, _ = establisheConnection(options.Target, options.TLS, options.Timeout, options.ProxyTCP)
	}
	sr := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"defaultNamingContext"},
		nil)
	options.Log.Debug("LDAP request (objectClass=*) with attribut defaultNamingContext")
	res, err := options.ldapConn.Search(sr)
	if err != nil {
		return err
	}
	if len(res.Entries) == 0 {
		return fmt.Errorf("error getting metadata: No LDAP responses from server")
	}
	defaultNamingContext := res.Entries[0].GetAttributeValue("defaultNamingContext")
	if defaultNamingContext == "" {
		return fmt.Errorf("error getting metadata: attribute defaultNamingContext missing")
	}
	options.BaseDN = defaultNamingContext
	return nil

}
