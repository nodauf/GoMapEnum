package orchestrator

import (
	"GoMapEnum/src/utils"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (orchestrator *Orchestrator) Bruteforce(optionsModules Options) []string {
	var usernameList, passwordList []string
	optionsInterface := reflect.ValueOf(optionsModules).Interface()
	options := optionsModules.GetBaseOptions()
	var wg sync.WaitGroup
	var validUsers []string
	mux := &sync.Mutex{}
	if orchestrator.PreActionBruteforce != nil {
		orchestrator.PreActionBruteforce(&optionsInterface)
	}

	if options.CheckIfValid {
		if orchestrator.CustomOptionsForCheckIfValid != nil {
			optionsEnum := orchestrator.CustomOptionsForCheckIfValid(&optionsInterface)
			usernameList = orchestrator.UserEnum(optionsEnum.(Options))
		} else {
			usernameList = orchestrator.UserEnum(optionsModules)
		}
	} else {
		options.Users = utils.GetStringOrFile(options.Users)
		usernameList = strings.Split(options.Users, "\n")
	}
	options.Passwords = utils.GetStringOrFile(options.Passwords)
	passwordList = strings.Split(options.Passwords, "\n")

	queue := make(chan string)
	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var j = 0
			for username := range queue {
				options.Log.Verbose("Testing " + username)
				if options.Sleep != 0 {
					// Sleep to avoid detection and bypass rate-limiting
					options.Log.Debug("Sleep " + strconv.Itoa(options.Sleep) + " seconds")
					time.Sleep(time.Duration(options.Sleep) * time.Second)
				}
				if options.NoBruteforce {
					if orchestrator.AuthenticationFunc(&optionsInterface, username, passwordList[j]) {
						mux.Lock()
						validUsers = append(validUsers, username+" / "+passwordList[j])
						mux.Unlock()
					}
					//options.authenticate(email, passwordList[j], &nbLockout)
				} else {
					for _, password := range passwordList {
						if orchestrator.AuthenticationFunc(&optionsInterface, username, password) {
							mux.Lock()
							validUsers = append(validUsers, username+" / "+password)
							mux.Unlock()
							break // No need to continue if password is valid
						}
					}
				}
				j++
				options.Log.Verbose("No password matched for " + username)
			}

		}(i)
	}
	// Trim emails and send them to the pool of workers
	for _, email := range usernameList {
		email = strings.ToValidUTF8(email, "")
		email = strings.Trim(email, "\r")
		email = strings.Trim(email, "\n")
		queue <- email
	}

	close(queue)
	wg.Wait()
	return validUsers
}
