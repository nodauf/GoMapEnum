package o365

import (
	"GoMapEnum/src/utils"
	"strings"
	"sync"
)

// UserEnum: Return a valid list of users according the provided options
func (options *Options) UserEnum() []string {
	mux := &sync.Mutex{}
	options.Users = utils.GetStringOrFile(options.Users)
	emailList := strings.Split(options.Users, "\n")
	var wg sync.WaitGroup
	var validusers []string
	queue := make(chan string)

	for i := 0; i < options.Thread; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for email := range queue {
				domain := strings.Split(email, "@")[1]
				if !options.validTenant(domain) {
					options.Log.Error("Tenant " + domain + " is not valid")
					return
				} else {
					options.Log.Verbose("Tenant " + domain + " is valid")
				}
				switch options.Mode {
				case "office":
					if options.enumOffice(email) {
						mux.Lock()
						validusers = append(validusers, email)
						mux.Unlock()

					}
				case "oauth2":
					if options.enumOauth2(email) {
						mux.Lock()
						validusers = append(validusers, email)
						mux.Unlock()
					}
				case "onedrive":
					if options.enumOnedrive(email) {
						mux.Lock()
						validusers = append(validusers, email)
						mux.Unlock()
					}
				}
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
	return validusers
}
