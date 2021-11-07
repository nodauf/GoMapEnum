package owa

import (
	"GoMapEnum/src/utils"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// UserEnum: Return a valid list of users according the provided options
func (options *Options) UserEnum() []string {
	log = options.Log

	options.Users = utils.GetStringOrFile(options.Users)

	return options.determineValidUsers()
}

func (options *Options) determineValidUsers() []string {
	var wg sync.WaitGroup
	mux := &sync.Mutex{}
	queue := make(chan string)
	userList := strings.Split(options.Users, "\n")

	/*Keep in mind you, nothing has been added to handle successful auths
	  so the password for auth attempts has been hardcoded to something
	  that is not likely to be correct.
	*/
	pass := "Summer2018978"
	// Prepare the transport for all the requests
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           options.Proxy,
	}
	options.tr = tr
	urlToHarvest := options.getURIToAuthenticate(options.Target)
	internaldomain := options.harvestInternalDomain(urlToHarvest)
	avgResponse := options.basicAuthAvgTime(urlToHarvest, internaldomain)
	log.Info("Internal Domain: " + internaldomain)
	var validusers []string

	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for user := range queue {
				startTime := time.Now()
				webRequestBasicAuth(urlToHarvest, internaldomain+"\\"+user, pass, options.tr)
				elapsedTime := time.Since(startTime)

				if float64(elapsedTime) < float64(avgResponse)*0.77 {
					mux.Lock()
					log.Success(user + " - " + elapsedTime.String())
					validusers = append(validusers, user)
					mux.Unlock()
				} else {
					mux.Lock()
					log.Fail(user + " - " + fmt.Sprint(elapsedTime))
					mux.Unlock()
				}
			}
		}(i)
	}

	// Send the users to the pool of workers
	for i := 0; i < len(userList); i++ {
		queue <- userList[i]
	}

	close(queue)
	wg.Wait()
	return validusers
}
