package azure

import (
	"GoMapEnum/src/utils"
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// BASE_XML is the data of the request which will be sent to authenticate
var BASE_XML = `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
  <s:Header>
    <a:Action s:mustUnderstand="1">http://schemas.xmlsoap.org/ws/2005/02/trust/RST/Issue</a:Action>
    <a:MessageID>MessageIDPlaceholder</a:MessageID>
    <a:ReplyTo>
      <a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
    </a:ReplyTo>
    <a:To s:mustUnderstand="1">https://autologon.microsoftazuread-sso.com/dewi.onmicrosoft.com/winauth/trust/2005/usernamemixed?client-request-id=30cad7ca-797c-4dba-81f6-8b01f6371013</a:To>
    <o:Security xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" s:mustUnderstand="1">
      <u:Timestamp u:Id="_0">
        <u:Created>%s</u:Created>
        <u:Expires>%s</u:Expires>
      </u:Timestamp>
      <o:UsernameToken u:Id="UsernameTokenPlaceholder">
        <o:Username>%s</o:Username>
        <o:Password>%s</o:Password>
      </o:UsernameToken>
    </o:Security>
  </s:Header>
  <s:Body>
    <trust:RequestSecurityToken xmlns:trust="http://schemas.xmlsoap.org/ws/2005/02/trust">
      <wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
        <a:EndpointReference>
          <a:Address>urn:federation:MicrosoftOnline</a:Address>
        </a:EndpointReference>
      </wsp:AppliesTo>
      <trust:KeyType>http://schemas.xmlsoap.org/ws/2005/05/identity/NoProofKey</trust:KeyType>
      <trust:RequestType>http://schemas.xmlsoap.org/ws/2005/02/trust/Issue</trust:RequestType>
    </trust:RequestSecurityToken>
  </s:Body>
</s:Envelope>
`

// AZURE_URL is the url to authenticate on
var AZURE_URL = "https://autologon.microsoftazuread-sso.com/%s/winauth/trust/2005/usernamemixed?client-request-id=%s"

func UserEnum(optionsInterface *interface{}, username string) bool {
	options := (*optionsInterface).(*Options)
	valid := false
	if len(strings.Split(username, "@")) == 1 {
		options.Log.Error("Only email format is supported, skipping " + username)
		return false
	}
	client := &http.Client{

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           options.ProxyHTTP,
		},
	}
	domain := strings.Split(username, "@")[1]
	// Random password for authentication
	password := utils.RandomString(10)
	// Generate time for the POST data
	now := time.Now()
	created := now.Format(time.RFC3339)
	expired := now.Add(10 * time.Minute).Format(time.RFC3339)
	// Random UUID for the POST data
	uuid, _ := utils.NewUUID()
	dataToSend := fmt.Sprintf(BASE_XML, created, expired, username, password)
	url := fmt.Sprintf(AZURE_URL, domain, uuid)
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(dataToSend))
	resp, err := client.Do(req)
	if err != nil {
		options.Log.Error("Error on response.\n[ERRO] - " + err.Error())
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var structResponseAzure azureResponse
	xml.Unmarshal(body, &structResponseAzure)
	// Parse the response
	code := strings.Split(structResponseAzure.Body.Fault.Detail.Error.Internalerror.Text, ":")[0]
	switch code { // https://docs.microsoft.com/en-us/azure/active-directory/develop/reference-aadsts-error-codes
	case "AADSTS81016":
		options.Log.Success("The user " + username + " seems to exist")
		valid = true
	case "AADSTS50053":
		options.Log.Info(username + " is locked")
	case "AADSTS50126":
		options.Log.Success(username + " exists but the password is wrong")
		valid = true
	case "AADSTS50055":
		options.Log.Success(username + " exists but the password is expired")
		valid = true
	case "AADSTS50056":
		options.Log.Success(username + " exists but there is no password")
		valid = true
	case "AADSTS50014":
		options.Log.Success(username + " exists but max passthru auth time exceeded")
		valid = true
	case "AADSTS50076": // Due to a configuration change made by your administrator, or because you moved to a new location, you must use multi-factor authentication to access
		options.Log.Info(username + " MFA needed")
	case "AADSTS700016":
		options.Log.Error(username + " The application wasn't found in the directory/tenant")
	case "AADSTS50034": // UserAccountNotFound - To sign into this application, the account must be added to the directory.
		options.Log.Fail(username + " does not exist")
	case "AADSTS90002":
		options.Log.Error("The Tenant '" + domain + "' does not exist")
	default:
		options.Log.Error("Unknow error: " + structResponseAzure.Body.Fault.Detail.Error.Internalerror.Text)

	}
	return valid
}
