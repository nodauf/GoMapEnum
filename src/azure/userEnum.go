package azure

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
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

var log *logger.Logger

// UserEnum return a valid list of users according the provided options
func (options *Options) UserEnum() []string {
	log = options.Log
	mux := &sync.Mutex{}
	var validusers []string
	options.Users = utils.GetStringOrFile(options.Users)
	emailList := strings.Split(options.Users, "\n")
	var wg sync.WaitGroup
	queue := make(chan string)

	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for email := range queue {
				valid := false
				if len(strings.Split(email, "@")) == 1 {
					log.Error("Only email format is supported, skipping " + email)
					break
				}
				domain := strings.Split(email, "@")[1]
				// Random password for authentication
				password := utils.RandomString(10)
				// Generate time for the POST data
				now := time.Now()
				created := now.Format(time.RFC3339)
				expired := now.Add(10 * time.Minute).Format(time.RFC3339)
				// Random UUID for the POST data
				uuid, _ := utils.NewUUID()
				dataToSend := fmt.Sprintf(BASE_XML, created, expired, email, password)
				url := fmt.Sprintf(AZURE_URL, domain, uuid)
				req, _ := http.NewRequest("POST", url, bytes.NewBufferString(dataToSend))
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Error("Error on response.\n[ERRO] - " + err.Error())
				}

				body, _ := ioutil.ReadAll(resp.Body)
				var structResponseAzure azureResponse
				xml.Unmarshal(body, &structResponseAzure)
				// Parse the response
				code := strings.Split(structResponseAzure.Body.Fault.Detail.Error.Internalerror.Text, ":")[0]
				switch code { // https://docs.microsoft.com/en-us/azure/active-directory/develop/reference-aadsts-error-codes
				case "AADSTS81016":
					log.Success("The user " + email + " seems to exist")
					valid = true
				case "AADSTS50053":
					log.Info(email + " is locked")
				case "AADSTS50126":
					log.Success(email + " exists but the password is wrong")
					valid = true
				case "AADSTS50055":
					log.Success(email + " exists but the password is expired")
					valid = true
				case "AADSTS50056":
					log.Success(email + " exists but there is no password")
					valid = true
				case "AADSTS50014":
					log.Success(email + " exists but max passthru auth time exceeded")
					valid = true
				case "AADSTS50076": // Due to a configuration change made by your administrator, or because you moved to a new location, you must use multi-factor authentication to access
					log.Info(email + " MFA needed")
				case "AADSTS700016":
					log.Fail(email + " The application wasn't found in the directory/tenant")
				case "AADSTS50034": // UserAccountNotFound - To sign into this application, the account must be added to the directory.
					log.Fail(email + " does not exist")
				case "AADSTS90002":
					log.Fail("The Tenant '" + domain + "' does not exist")
				default:
					log.Error("Unknow error: " + structResponseAzure.Body.Fault.Detail.Error.Internalerror.Text)

				}
				if valid {
					mux.Lock()
					validusers = append(validusers, email)
					mux.Unlock()
				}
			}
		}(i)
	}

	// Trim emails and send them to the pool of workers
	for _, email := range emailList {
		email = strings.ToValidUTF8(email, "")
		email = strings.Trim(email, "\r")
		email = strings.Trim(email, "\n")
		if email == "" {
			continue
		}
		queue <- email
	}

	close(queue)
	wg.Wait()
	return validusers
}
