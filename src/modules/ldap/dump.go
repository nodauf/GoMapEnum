package ldap

import (
	"GoMapEnum/src/utils"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func (options *Options) InitLDAP() bool {
	optionsInterface := reflect.ValueOf(options).Interface()
	options.Log.Debug("Initializing LDAP")
	if !RetrieveTargetInfo(&optionsInterface) {
		options.Log.Error("Cannot initialize LDAP")
		return false
	}
	valid, err := options.authenticate(options.Users, options.Passwords)
	if !valid || err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			options.Log.Error("fail to authenticate: Invalid credential")
		} else {
			options.Log.Error("fail connect to the LDAP. Unkown error: %v", err)
		}
		return false
	}
	err = options.GetDefaultNamingContext()
	if err != nil {
		options.Log.Error("fail to retrieve the default naming context, error: %v", err)
		return false
	}
	options.Log.Debug("Naming context: " + options.BaseDN)
	return true
}

func (options *Options) Dump() string {
	if !options.InitLDAP() {
		options.Log.Error("Cannot initialize LDAP")
		return ""
	}
	defer options.ldapConn.Close()

	var ldapData []*ldap.Entry
	if utils.StringInSlice(options.DumpObjects, "all") {
		for object := range options.queries {
			ldapData = append(ldapData, options.DumpObject(object)...)
		}
	} else {
		//objectsToDump := strings.Split(strings.ToLower(options.DumpObjects), ",")
		for _, object := range options.DumpObjects {
			options.Log.Debug("Dumping %s", object)
			object = strings.ToLower(object)
			ldapData = append(ldapData, options.DumpObject(object)...)
			if ldapData == nil {
				continue
			}
			if options.JSON {
				jsonOutput, err := json.MarshalIndent(&ldapData, "", "\t")
				// If cannot convert to json, return the raw data
				if err != nil {
					options.Log.Error("fail to convert the data to json, error: %v", err)
				} else {
					ioutil.WriteFile(object+".json", []byte(jsonOutput), fs.ModePerm)

				}
				if err != nil {
					options.Log.Error("fail to convert data to json: %s", err.Error())
					// return the data to avoid loosing them

				}
				return string(jsonOutput)
			}

			if options.HTML {
				columns := strings.Split(options.queries[object]["attributs"], ",")
				rows := ParseLDAPData(ldapData, columns)
				template := utils.DataToHTML(rows, columns, object)
				ioutil.WriteFile(object+".html", template.Bytes(), fs.ModePerm)

			}
		}
	}
	var output string
	for _, entry := range ldapData {
		output += fmt.Sprintf("DN: %s\n", entry.DN)
		for _, attr := range entry.Attributes {
			output += fmt.Sprintf("%s: %s\n", attr.Name, attr.Values)
		}
	}
	return output

}

func (options *Options) DumpObject(object string) []*ldap.Entry {
	if options.queries == nil {
		options.initDumpMap()
	}
	if _, ok := options.queries[object]; !ok {
		options.Log.Error("Not able to dump %s. The query is not implemented.", object)
		return nil
	}
	ldapResult := executeLdapQuery(options.ldapConn, options.BaseDN, options.queries[object])

	// Print the results
	for _, entry := range ldapResult.Entries {
		options.Log.Success(object + ": " + entry.DN)
	}
	return ldapResult.Entries
}
