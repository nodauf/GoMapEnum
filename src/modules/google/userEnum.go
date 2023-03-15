package google

import (
	"GoMapEnum/src/utils"
	"crypto/tls"
	"fmt"
	"net/http"
)

var URL_GOOGLE_USER_ENUM = "https://mail.google.com/mail/gxlu?email=%s"

func UserEnum(optionsInterface *interface{}, username string) bool {

	options := (*optionsInterface).(*Options)
	var valid = false

	url := fmt.Sprintf(URL_GOOGLE_USER_ENUM, username)
	header := make(map[string]string)
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; rv:61.0) Gecko/20100101 Firefox/61.0"
	header["Accept-Language"] = "en-US,en;q=0.5"
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           options.ProxyHTTP,
		},
	}
	req.Header.Set("User-Agent", utils.GetUserAgent())
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	resp, err := client.Do(req)
	if err != nil {
		options.Log.Error("Error on response.\n[ERRO] - " + err.Error())
		return false
	}

	if len(resp.Cookies()) > 0 {
		options.Log.Success(username)
		valid = true
	} else {
		options.Log.Fail(username)
	}

	return valid
}
