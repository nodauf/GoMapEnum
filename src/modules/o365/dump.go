package o365

import (
	"GoMapEnum/src/utils"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
)

var dumpURLS = map[string]string{
	"users":                  "https://graph.windows.net/%s/users?api-version=1.61-internal",
	"tenantDetails":          "https://graph.windows.net/%s/tenantDetails?api-version=1.61-internal",
	"policies":               "https://graph.windows.net/%s/policies?api-version=1.61-internal",
	"servicePrincipals":      "https://graph.windows.net/%s/servicePrincipals?api-version=1.61-internal",
	"groups":                 "https://graph.windows.net/%s/groups?api-version=1.61-internal",
	"applications":           "https://graph.windows.net/%s/applications?api-version=1.61-internal",
	"devices":                "https://graph.windows.net/%s/devices?api-version=1.61-internal",
	"directoryRoles":         "https://graph.windows.net/%s/directoryRoles?api-version=1.61-internal",
	"roleDefinitions":        "https://graph.windows.net/%s/roleDefinitions?api-version=1.61-internal",
	"contacts":               "https://graph.windows.net/%s/contacts?api-version=1.61-internal",
	"oauth2PermissionGrants": "https://graph.windows.net/%s/oauth2PermissionGrants?api-version=1.61-internal",
}

func (options *Options) Dump() string {
	username := options.Users
	password := options.Passwords
	respStruct := options.requestOauth2(username, password)
	if respStruct.AccessToken == "" {
		errorMessage := "Something is wrong. Check the credentials. "
		if respStruct.ErrorDescription != "" {
			code := strings.Split(respStruct.ErrorDescription, ":")[0]

			switch code { // https://docs.microsoft.com/en-us/azure/active-directory/develop/reference-aadsts-error-codes
			case "AADSTS50053":
				errorMessage += username + " is locked"
			case "AADSTS50126":
				errorMessage += " exists but the password is wrong"
			case "AADSTS50055":
				errorMessage += username + " exists but the password is expired"
			case "AADSTS50056":
				errorMessage += username + " exists but there is no password"
			case "AADSTS50014":
				errorMessage += username + " exists but max passthru auth time exceeded"
			case "AADSTS50076": // Due to a configuration change made by your administrator, or because you moved to a new location, you must use multi-factor authentication to access
				errorMessage += username + " MFA needed"
			case "AADSTS50057":
				errorMessage += username + " and " + password + " matched but the account is disabled"
			case "AADSTS700016":
				errorMessage += username + " The application wasn't found in the directory/tenant"
			case "AADSTS50034": // UserAccountNotFound - To sign into this application, the account must be added to the directory.
				errorMessage += username + " does not exist"
			case "AADSTS90002":
				errorMessage += "The Tenant '" + username + "' does not exist"
			default:
				errorMessage += "Unknow error: " + respStruct.ErrorDescription
			}
		}
		// The access token is empty. We exit
		return ""
	}
	tenantID, err := getTenantIDFromAccessToken(respStruct.AccessToken)
	if err != nil {
		options.Log.Error("cannot retrieve the tenant ID in the access token: " + err.Error())
		return ""
	}
	dataToDump := strings.Split(options.DumpObjects, ",")
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + respStruct.AccessToken
	// For each urls to query and dump, we send the request and save the json or hmtl
	for kind, url := range dumpURLS {
		if options.DumpObjects != "all" && !utils.StringInSlice(dataToDump, kind) {
			continue
		}
		urlWithTenant := fmt.Sprintf(url, tenantID)
		options.Log.Debug("Request url %s", urlWithTenant)
		dataInterface, err := dumpO365ObjectPaging(urlWithTenant, options.ProxyHTTP, headers)
		if err != nil {
			options.Log.Error(err.Error())
		}

		if options.JSON {
			jsonData, _ := json.MarshalIndent(dataInterface, "", "  ")
			ioutil.WriteFile(kind+".json", []byte(jsonData), fs.ModePerm)
		}
		if options.HTML {
			switch kind {
			case "users":
				var users Users
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &users)

				var columns = []string{"DisplayName", "Mail", "ProxyAddresses"}
				data := parseO365Data(users.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for users from O365")
				ioutil.WriteFile("users.html", template.Bytes(), fs.ModePerm)

			case "tenantDetails":
				var tenantDetails TenantDetails
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &tenantDetails)

				var columns = []string{"DisplayName"}
				data := parseO365Data(tenantDetails.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for tenantDetails from O365")
				ioutil.WriteFile("tenantDetails.html", template.Bytes(), fs.ModePerm)

			case "policies":
				var policies Policies
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &policies)

				var columns = []string{"DisplayName"}
				data := parseO365Data(policies.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for policies from O365")
				ioutil.WriteFile("policies.html", template.Bytes(), fs.ModePerm)

			case "servicePrincipals":
				var servicePrincipals ServicePrincipals
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &servicePrincipals)

				var columns = []string{"DisplayName"}
				data := parseO365Data(servicePrincipals.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for servicePrincipals from O365")
				ioutil.WriteFile("servicePrincipals.html", template.Bytes(), fs.ModePerm)

			case "groups":
				var groups Groups
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &groups)

				var columns = []string{"DisplayName", "Mail"}
				data := parseO365Data(groups.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for groups from O365")
				ioutil.WriteFile("groups.html", template.Bytes(), fs.ModePerm)

			case "applications":
				var applications Application
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &applications)

				var columns = []string{"DisplayName"}
				data := parseO365Data(applications.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for applications from O365")
				ioutil.WriteFile("applications.html", template.Bytes(), fs.ModePerm)

			case "devices":
				var devices Devices
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &devices)

				var columns = []string{"DisplayName"}
				data := parseO365Data(devices.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for devices from O365")
				ioutil.WriteFile("devices.html", template.Bytes(), fs.ModePerm)

			case "directoryRoles":
				var directoryRoles DirectoryRoles
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &directoryRoles)

				var columns = []string{"DisplayName"}
				data := parseO365Data(directoryRoles.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for directoryRoles from O365")
				ioutil.WriteFile("directoryRoles.html", template.Bytes(), fs.ModePerm)

			case "roleDefinitions":
				var roleDefinitions RoleDefinitions
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &roleDefinitions)

				var columns = []string{"DisplayName"}
				data := parseO365Data(roleDefinitions.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for roleDefinitions from O365")
				ioutil.WriteFile("roleDefinitions.html", template.Bytes(), fs.ModePerm)

			case "contacts":
				/*var contacts Contacts
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &contacts)

				var columns = []string{"DisplayName"}
				data := parseO365Data(contacts.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for contacts from O365")
				ioutil.WriteFile("contacts.html", template.Bytes(), fs.ModePerm)*/

				// Not implemented miss some data

			case "oauth2PermissionGrants":
				var oauth2PermissionGrants Oauth2PermissionGrants
				jsonData, _ := json.Marshal(dataInterface)
				json.Unmarshal([]byte(jsonData), &oauth2PermissionGrants)

				var columns = []string{"Scope"}
				data := parseO365Data(oauth2PermissionGrants.Value, columns)
				template := utils.DataToHTML(data, columns, "Results for oauth2PermissionGrants from O365")
				ioutil.WriteFile("oauth2PermissionGrants.html", template.Bytes(), fs.ModePerm)

			}

		}
	}
	return ""
}

// IsObjectCanBeDumped return true if it's available in the map dumpURLS
func IsObjectCanBeDumped(object string) bool {
	for name := range dumpURLS {
		if name == object {
			return true
		}
	}
	return false
}
