package owa

import (
	"crypto/tls"
	"net/http"
)

func PrepareOptions(optionsInterface *interface{}) interface{} {
	// Make a copy of the options
	return *optionsInterface
}

func PrepareBruteforce(optionsInterface *interface{}) bool {
	options := (*optionsInterface).(*Options)
	// Prepare the transport for all the requests
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           options.ProxyHTTP,
	}
	options.tr = tr
	options.urlToHarvest = options.getURIToAuthenticate(options.Target)
	options.internalDomain = options.harvestInternalDomain(options.urlToHarvest)
	options.Log.Info("Internal Domain: " + options.internalDomain)
	return true
}

func Authenticate(optionsInterface *interface{}, email, password string) bool {
	options := (*optionsInterface).(*Options)
	return options.webRequestBasicAuth(options.urlToHarvest, options.internalDomain+"\\"+email, password)
}
