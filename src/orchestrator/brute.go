package orchestrator

import (
	"GoMapEnum/src/utils"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (orchestrator *Orchestrator) Bruteforce(optionsModules Options) string {
	var usernameList, passwordList []string
	optionsInterface := reflect.ValueOf(optionsModules).Interface()
	options := optionsModules.GetBaseOptions()
	var wg sync.WaitGroup
	var validUsers []string
	mux := &sync.Mutex{}
	if orchestrator.PreActionBruteforce != nil {
		if !orchestrator.PreActionBruteforce(&optionsInterface) {
			return strings.Join(validUsers, "\n")
		}
	}
	if options.CheckIfValid {
		if orchestrator.CustomOptionsForCheckIfValid != nil {
			optionsEnum := orchestrator.CustomOptionsForCheckIfValid(&optionsInterface)
			usernameList = strings.Split(orchestrator.UserEnum(optionsEnum.(Options)), "\n")
		} else {
			usernameList = strings.Split(orchestrator.UserEnum(optionsModules), "\n")
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
				found := false
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
						found = true
						options.Log.Success(username + " / " + passwordList[j])
					} else {
						options.Log.Fail(username + " / " + passwordList[j])
					}
				} else {
					for _, password := range passwordList {
						if orchestrator.AuthenticationFunc(&optionsInterface, username, password) {
							mux.Lock()
							validUsers = append(validUsers, username+" / "+password)
							mux.Unlock()
							found = true
							options.Log.Success(username + " / " + password)
							break // No need to continue if password is valid
						} else {
							options.Log.Fail(username + " / " + password)
						}
					}
				}
				j++
				if !found {
					options.Log.Verbose("No password matched for " + username)
				}
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
	if orchestrator.PostActionBruteforce != nil {
		if !orchestrator.PostActionBruteforce(&optionsInterface) {
			return strings.Join(validUsers, "\n")
		}
	}
	return strings.Join(validUsers, "\n")
}
