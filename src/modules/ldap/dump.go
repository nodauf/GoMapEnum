package ldap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func (options *Options) initDumpMap() {
	options.queries = make(map[string]map[string]string)

	computers := make(map[string]string)
	computers["filter"] = "(objectClass=Computer)"
	computers["attributs"] = "cn,dNSHostName,operatingSystem,operatingSystemVersion,operatingSystemServicePack,whenCreated,lastLogon,objectSid,objectClass"
	options.queries["computers"] = computers

	users := make(map[string]string)
	users["filter"] = "(objectClass=user)"
	users["attributs"] = "cn,sAMAccountName,userPrincipalName,objectClass"
	options.queries["users"] = users
}
func (options *Options) Dump() string {
	optionsInterface := reflect.ValueOf(options).Interface()
	options.initDumpMap()
	RetrieveTargetInfo(&optionsInterface)
	options.Log.Verbose("Using domain " + options.Domain + " for authentication. Hostname: " + options.Hostname)
	err := options.authenticateNTLM(options.Users, options.Passwords, options.IsHash)
	defer options.ldapConn.Close()
	if err != nil {
		if !strings.Contains(err.Error(), "Invalid Credentials") {
			options.Log.Error("fail to authenticate: %s", err.Error())
		} else {
			options.Log.Error("cannot connect to the LDAP")
		}
		return ""
	}
	err = options.getDefaultNamingContext()
	if err != nil {
		options.Log.Error("fail to retrieve the default naming context")
		return ""
	}
	options.Log.Debug("Naming context: " + options.BaseDN)
	var ldapData []*ldap.Entry
	if strings.ToLower(options.DumpObjects) == "all" {
		for object := range options.queries {
			ldapData = append(ldapData, options.dumpObject(object)...)
		}
	} else {
		objectsToDump := strings.Split(strings.ToLower(options.DumpObjects), ",")
		for _, object := range objectsToDump {
			ldapData = append(ldapData, options.dumpObject(object)...)
		}
	}

	jsonOutput, err := json.MarshalIndent(&ldapData, "", "\t")
	if err != nil {
		options.Log.Error("fail to convert data to json: %s", err.Error())
		// return the data to avoid loosing them
		var output string
		for _, entry := range ldapData {
			output += fmt.Sprintf("DN: %s\n", entry.DN)
			for _, attr := range entry.Attributes {
				output += fmt.Sprintf("%s: %s\n", attr.Name, attr.Values)
			}
		}
		return output
	}
	return string(jsonOutput)

}
