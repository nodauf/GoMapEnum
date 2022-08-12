package orchestrator

import (
	"GoMapEnum/src/utils"
	"reflect"
	"strings"
	"sync"
)

// UserEnum will call the functions according the orchestrator options to enumerate valid users.
// Firstly, PreActionUserEnum
// Then for each user, the function CheckBeforeEnumFunc
// After that, UserEnumFunc
// Finally, PostActionUserEnum
func (orchestrator *Orchestrator) UserEnum(optionsModules Options) string {
	optionsInterface := reflect.ValueOf(optionsModules).Interface()
	options := optionsModules.GetBaseOptions()
	options.Users = utils.GetStringOrFile(options.Users)
	options.UsernameList = strings.Split(options.Users, "\n")
	mux := &sync.Mutex{}
	var wg sync.WaitGroup
	var validUsers []string
	queue := make(chan string)
	if orchestrator.PreActionUserEnum != nil {
		// If the PreActionUserEnum failed, just returned the list that is empty at this step
		if !orchestrator.PreActionUserEnum(&optionsInterface) {
			return strings.Join(validUsers, "\n")
		}
	}
	// Start the workers
	for i := 0; i < options.Thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for username := range queue {
				options.Log.Verbose("Testing " + username)
				if orchestrator.CheckBeforeEnumFunc != nil {
					// If the check did not pass for that user skip it
					if !orchestrator.CheckBeforeEnumFunc(&optionsInterface, username) {
						continue
					}
				}
				if orchestrator.UserEnumFunc(&optionsInterface, username) {
					options.Log.Debug(username + " exists")
					mux.Lock()
					validUsers = append(validUsers, username)
					mux.Unlock()
				}

				options.Log.Debug(username + " does not exist")
			}
		}(i)
	}

	// Trim usernames and send them to the pool of workers
	for _, username := range options.UsernameList {
		username = strings.ToValidUTF8(username, "")
		username = strings.Trim(username, "\r")
		username = strings.Trim(username, "\n")
		if username == "" {
			continue
		}
		queue <- username
	}

	close(queue)
	// Wait all the workers
	wg.Wait()

	// Doing the post action
	if orchestrator.PostActionUserEnum != nil {
		if !orchestrator.PostActionUserEnum(&optionsInterface) {
			return strings.Join(validUsers, "\n")
		}
	}
	return strings.Join(validUsers, "\n")
}
