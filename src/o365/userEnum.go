package o365

import (
	"GoMapEnum/src/utils"
	"strings"
	"sync"
)

// UserEnum return a valid list of users according the provided options
func (options *Options) UserEnum() []string {
	mux := &sync.Mutex{}
	options.Users = utils.GetStringOrFile(options.Users)
	emailList := strings.Split(options.Users, "\n")
	var wg sync.WaitGroup
	var validusers []string
	domainValidated := make(map[string]bool)
	queue := make(chan string)

	for i := 0; i < options.Thread; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for email := range queue {
				domain := strings.Split(email, "@")[1]
				// If we didn't already checked the domain
				mux.Lock()
				if domainValid, ok := domainValidated[domain]; !ok {
					if !options.validTenant(domain) {
						options.Log.Error("Tenant " + domain + " is not valid")
						domainValidated[domain] = false
						mux.Unlock()
						continue
					}
					options.Log.Info("Tenant " + domain + " is valid")
					domainValidated[domain] = true
				} else {
					// If the domain was not valid, skip the email
					if !domainValid {
						options.Log.Debug("Tenant " + domain + " already checked and was not valid")
						mux.Unlock()
						continue
					}
				}
				mux.Unlock()

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
		if email == "" {
			continue
		}
		queue <- email
	}

	close(queue)
	wg.Wait()
	return validusers
}
