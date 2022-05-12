package orchestrator

import (
	"GoMapEnum/src/utils"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Bruteforce will call the functions according the orchestrator options to bruteforce.
// Firstly, PreActionBruteforce
// Then, CustomOptionsForCheckIfValid mostly to change the log level for enumeration and avoid printing success if the user a checked before the bruteforce
// Afterthat, AuthenticationFunc to authenticate with the username and password
// Finally, PostActionBruteforce
func (orchestrator *Orchestrator) Bruteforce(optionsModules Options) string {
	var usernameList, passwordList []string
	optionsInterface := reflect.ValueOf(optionsModules).Interface()
	options := optionsModules.GetBaseOptions()
	var wg sync.WaitGroup
	var validUsers []string
	mux := &sync.Mutex{}

	if orchestrator.PreActionBruteforce != nil {
		// If the PreActionUserEnum failed, just returned the list that is empty at this step
		if !orchestrator.PreActionBruteforce(&optionsInterface) {
			return strings.Join(validUsers, "\n")
		}
	}
	if options.CheckIfValid {
		options.Log.Debug("Validating the users")
		if orchestrator.CustomOptionsForCheckIfValid != nil {
			optionsEnum := optionsInterface
			if orchestrator.CustomOptionsForCheckIfValid != nil {
				// Execute the function CustomOptionsForCheckIfValid mainly to change the log level and print only errors
				optionsEnum = orchestrator.CustomOptionsForCheckIfValid(&optionsInterface)
			}
			// If wee check if the users are valid before, we retrieve the usernameList from the Userenum function
			usernameList = strings.Split(orchestrator.UserEnum(optionsEnum.(Options)), "\n")
		} else {
			// Switch the mode to enumeration by default
			previousMode := optionsInterface.(Options).GetBaseOptions().Log.Mode
			optionsInterface.(Options).GetBaseOptions().Log.Mode = "Enumeration"
			usernameList = strings.Split(orchestrator.UserEnum(optionsModules), "\n")
			optionsInterface.(Options).GetBaseOptions().Log.Mode = previousMode
		}
	} else {
		options.Users = utils.GetStringOrFile(options.Users)
		usernameList = strings.Split(options.Users, "\n")
	}
	options.Passwords = utils.GetStringOrFile(options.Passwords)
	passwordList = strings.Split(options.Passwords, "\n")

	queue := make(chan string)
	// Start the workers
	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for username := range queue {
				found := false
				options.Log.Verbose("Testing " + username)
				if options.Sleep != 0 {
					// Sleep to avoid detection and bypass rate-limiting
					options.Log.Debug("Sleeping " + strconv.Itoa(options.Sleep) + " seconds")
					time.Sleep(time.Duration(options.Sleep) * time.Second)
				}
				if options.NoBruteforce {
					// If no bruteforce we log with the first item in username list and password list, then the second and so one...
					index := utils.IndexInSlice(usernameList, username)
					if orchestrator.AuthenticationFunc(&optionsInterface, username, passwordList[index]) {
						mux.Lock()
						validUsers = append(validUsers, username+" / "+passwordList[index])
						mux.Unlock()
						found = true
						options.Log.Success(username + " / " + passwordList[index])
					} else {
						options.Log.Fail(username + " / " + passwordList[index])
					}
				} else {
					// If bruteforce we test for each password in the list
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
		if email != "" {
			queue <- email
		}
	}

	close(queue)
	// Wait all workers
	wg.Wait()

	// Doing the post action
	if orchestrator.PostActionBruteforce != nil {
		if !orchestrator.PostActionBruteforce(&optionsInterface) {
			return strings.Join(validUsers, "\n")
		}
	}
	return strings.Join(validUsers, "\n")
}
