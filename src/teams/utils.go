package teams

import (
	"GoMapEnum/src/logger"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// getPresence request the Teams API to get additional details about the user with its mri
func (options *Options) getPresence(mri, bearer string, log *logger.Logger) (string, string, string) {

	var jsonData = []byte(`[{"mri":"` + mri + `"}]`)
	req, _ := http.NewRequest("POST", URL_PRESENCE_TEAMS, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("x-ms-client-version", CLIENT_VERSION)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           options.Proxy,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error on response.\n[ERRO] - " + err.Error())
		return "", "", ""
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var status []struct {
		Mri      string `json:"mri"`
		Presence struct {
			Availability string `json:"availability"`
			DeviceType   string `json:"deviceType"`
			CalendarData struct {
				OutOfOfficeNote struct {
					Message     string    `json:"message"`
					PublishTime time.Time `json:"publishTime"`
					Expiry      time.Time `json:"expiry"`
				} `json:"outOfOfficeNote"`
				IsOutOfOffice bool `json:"isOutOfOffice"`
			} `json:"calendarData"`
		} `json:"presence"`
	}

	json.Unmarshal([]byte(body), &status)

	if len(status) > 0 {
		return status[0].Presence.Availability, status[0].Presence.DeviceType, status[0].Presence.CalendarData.OutOfOfficeNote.Message
	}
	return "", "", ""

}
