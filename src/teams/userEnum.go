package teams

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var URL_PRESENCE_TEAMS = "https://presence.teams.microsoft.com/v1/presence/getpresence/"
var URL_TEAMS = "https://teams.microsoft.com/api/mt/emea/beta/users/%s/externalsearchv3"
var CLIENT_VERSION = "27/1.0.0.2021011237"

// UserEnum: Return a valid list of users according the provided options
func (options *Options) UserEnum(logArg *logger.Logger) []string {
	log = logArg
	mux := &sync.Mutex{}
	options.Token = "Bearer " + options.Token
	options.Email = utils.GetStringOrFile(options.Email)
	emailList := strings.Split(options.Email, "\n")
	var wg sync.WaitGroup
	var validusers []string
	queue := make(chan string)

	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var valid = false
			for email := range queue {
				url := fmt.Sprintf(URL_TEAMS, email)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Add("Authorization", options.Token)
				req.Header.Add("x-ms-client-version", CLIENT_VERSION)

				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
						Proxy:           options.Proxy,
					},
				}
				resp, err := client.Do(req)
				if err != nil {
					log.Error("Error on response.\n[ERRO] - " + err.Error())
				}

				body, _ := ioutil.ReadAll(resp.Body)

				var jsonInterface interface{}
				var usefulInformation []struct {
					DisplayName string `json:"displayName"`
					Mri         string `json:"mri"`
				}

				json.Unmarshal([]byte(body), &jsonInterface)
				json.Unmarshal([]byte(body), &usefulInformation)

				log.Debug("Status code: " + strconv.Itoa(resp.StatusCode))

				bytes, _ := json.MarshalIndent(jsonInterface, "", " ")
				log.Debug("Response: " + string(bytes))

				if resp.StatusCode == 200 {
					if reflect.ValueOf(jsonInterface).Len() > 0 {
						presence, device, outOfOfficeNote := options.getPresence(usefulInformation[0].Mri, options.Token, log)
						log.Success(email + " - " + usefulInformation[0].DisplayName + " - " + presence + " - " + device + " - " + outOfOfficeNote)
						valid = true
					} else {
						log.Fail(email)
					}
					// If the status code is 403 it means the user exists but the organization did not enable connection from outside
				} else if resp.StatusCode == 403 {
					log.Success(email)
					valid = true
				} else if resp.StatusCode == 401 {
					log.Fail(email)
					log.Info("The token may be invalid or expired. The status code returned by the server is 401")
				} else {
					log.Fail(email)
					log.Error("Something went wrong. The status code returned by the server is " + strconv.Itoa(resp.StatusCode))
				}

				if valid {
					mux.Lock()
					validusers = append(validusers, email)
					mux.Unlock()
				}
			}
		}(i)
	}

	for _, email := range emailList {
		email = strings.ToValidUTF8(email, "")
		email = strings.Trim(email, "\r")
		email = strings.Trim(email, "\n")
		queue <- email
	}

	close(queue)
	wg.Wait()

	return validusers

}
