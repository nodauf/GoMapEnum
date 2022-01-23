package owa

import (
	"GoMapEnum/src/utils"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Azure/go-ntlmssp"
)

//harvestInternalDomain retrieve the internal domain name (netbios format)
func (options *Options) harvestInternalDomain(urlToHarvest string) string {
	options.Log.Verbose("Attempting to harvest internal domain:")

	timeout := time.Duration(3 * time.Second)

	client := &http.Client{
		Timeout:   timeout,
		Transport: options.tr,
	}
	req, _ := http.NewRequest("GET", urlToHarvest, nil)
	req.Header.Set("User-Agent", utils.GetUserAgent())
	req.Header.Set("Authorization", "NTLM TlRMTVNTUAABAAAAB4IIogAAAAAAAAAAAAAAAAAAAAAGAbEdAAAADw==")
	resp, err := client.Do(req)

	if err != nil {
		return ""
	}
	ntlmResponse := resp.Header.Get("WWW-Authenticate")
	data := strings.Split(ntlmResponse, " ")
	base64DecodedResp, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		options.Log.Error("Unable to parse NTLM response for internal domain name")
		return ""
	}

	var continueAppending bool
	var internalDomainDecimal []byte
	for _, decimalValue := range base64DecodedResp {
		if decimalValue == 0 {
			continue
		}
		if decimalValue == 2 {
			continueAppending = false
		}
		if continueAppending {
			internalDomainDecimal = append(internalDomainDecimal, decimalValue)
		}
		if decimalValue == 15 {
			continueAppending = true
			continue
		}
	}
	return string(internalDomainDecimal)
}

func (options *Options) getURIToAuthenticate(host string) string {
	// endpoint for NTLM authentication
	//url1 := "https://" + host + "/ews" // Disabled: Got some issue while testing for user enumeration, it was not working at all
	url2 := "https://" + host + "/autodiscover/autodiscover.xml"
	url3 := "https://" + host + "/rpc"
	url4 := "https://" + host + "/mapi"
	url5 := "https://" + host + "/oab"
	url6 := "https://autodiscover." + host + "/autodiscover/autodiscover.xml"
	var urlToHarvest string
	//if options.webRequestCodeResponse(url1) == 200 { // Disabled: Got some issue while testing for user enumeration, it was not working at all
	//	urlToHarvest = url1
	//} else if options.webRequestCodeResponse(url2) == 401 {
	if options.webRequestCodeResponse(url2) == 401 {
		urlToHarvest = url2
	} else if options.webRequestCodeResponse(url3) == 401 {
		urlToHarvest = url3
	} else if options.webRequestCodeResponse(url4) == 401 {
		urlToHarvest = url4
	} else if options.webRequestCodeResponse(url5) == 401 {
		urlToHarvest = url5
	} else if options.webRequestCodeResponse(url6) == 401 {
		urlToHarvest = url6
	} else {
		options.Log.Fatal("Unable to resolve host provided to harvest internal domain name")
	}
	options.Log.Verbose("OWA url that will be used: " + urlToHarvest)
	return urlToHarvest
}

// webRequestCodeResponse request an URI and return the status code
func (options *Options) webRequestCodeResponse(URI string) int {

	timeout := time.Duration(3 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: options.tr,
	}
	req, _ := http.NewRequest("GET", URI, nil)
	req.Header.Set("User-Agent", utils.GetUserAgent())
	resp, err := client.Do(req)
	if err != nil {
		options.Log.Error(err.Error())
	}
	return resp.StatusCode
}

// webRequestBasicAuth authenticate with basic auth on an URI
func (options *Options) webRequestBasicAuth(URI, user, pass string) bool {
	timeout := time.Duration(45 * time.Second)
	var client = &http.Client{}
	if options.Basic {
		client = &http.Client{
			Timeout:   timeout,
			Transport: options.tr,
		}
	} else {
		client = &http.Client{
			Timeout: timeout,
			Transport: ntlmssp.Negotiator{
				RoundTripper: options.tr,
			},
		}
	}
	req, _ := http.NewRequest("GET", URI, nil)
	req.Header.Set("User-Agent", utils.GetUserAgent())
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	if err != nil {
		options.Log.Error("Potential Timeout - " + user)
		options.Log.Error("One of your requests has taken longer than 45 seconds to respond.")
		options.Log.Error("Consider lowering amount of threads used for enumeration.")
		options.Log.Error(err.Error())
	}
	if resp.StatusCode == 500 {
		options.Log.Error("Something went wrong. Status code is 500")
		return false
	}
	if resp.StatusCode != 401 {
		return true
	}

	return false
}

// basicAuthAvgTime get an average response time for unknown users
func (options *Options) basicAuthAvgTime(urlToHarvest, internaldomain string) time.Duration {
	//We are determining sample auth response time for invalid users, the password used is irrelevant.
	pass := "Summer201823904"

	options.Log.Verbose("Collecting sample auth times...")

	var sliceOfTimes []float64
	var medianTime float64
	// randome users that probably do not exist
	usernamelist := []string{"sdfsdskljdfhkljhf", "ssdlfkjhgkjhdfsdfw", "sdfsdfdsfff", "sefsefsefsss", "lkjhlkjhiuyoiuy", "khiuoiuhohuio", "s2222dfs45g45gdf", "sdfseddf3333"}
	for i := 0; i < len(usernamelist)-1; i++ {
		startTime := time.Now()
		options.webRequestBasicAuth(urlToHarvest, internaldomain+"\\"+usernamelist[i], pass)
		elapsedTime := time.Since(startTime)
		if elapsedTime > time.Second*15 {
			options.Log.Error("Response taking longer than 15 seconds, setting time:")
			options.Log.Debug("Avg Response:" + fmt.Sprint(time.Duration(elapsedTime)))
			return time.Duration(elapsedTime)
		}
		// The first user has sometime an higher response time than the others
		if i != 0 {
			options.Log.Debug(fmt.Sprint(elapsedTime))
			sliceOfTimes = append(sliceOfTimes, float64(elapsedTime))
		}
	}
	sort.Float64s(sliceOfTimes)
	if len(sliceOfTimes)%2 == 0 {
		positionOne := len(sliceOfTimes)/2 - 1
		positionTwo := len(sliceOfTimes) / 2
		medianTime = (sliceOfTimes[positionTwo] + sliceOfTimes[positionOne]) / 2
	} else if len(sliceOfTimes)%2 != 0 {
		position := len(sliceOfTimes)/2 - 1
		medianTime = sliceOfTimes[position]
	} else {
		fmt.Println("Error determining whether length of times gathered is even or odd to obtain median value.")
	}
	options.Log.Debug("Avg Response:" + fmt.Sprint(time.Duration(medianTime)))
	return time.Duration(medianTime)
}
