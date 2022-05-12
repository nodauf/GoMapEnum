package o365

import (
	"strings"
)

// UserEnum return a valid list of users according the provided options
func UserEnum(optionsInterface *interface{}, username string) bool {
	options := (*optionsInterface).(*Options)

	switch options.Mode {
	case "office":
		return options.enumOffice(username)
	case "oauth2":
		return options.enumOauth2(username)
	case "onedrive":
		return options.enumOnedrive(username)

	}
	options.Log.Error("Invalid mode. It should be office, oauth2 or onedrive")
	return false

}

func CheckTenant(optionsInterface *interface{}, username string) bool {
	options := (*optionsInterface).(*Options)
	// If it's empty we initialize the map
	if len(options.validTenants) == 0 {
		options.validTenants = make(map[string]bool)
	}
	if len(strings.Split(username, "@")) == 1 {
		options.Log.Error("User should be in format user@tenant.tld")
		return false
	}
	domain := strings.Split(username, "@")[1]
	// If we didn't already checked the domain
	options.Mutex.Lock()
	if domainValid, ok := options.validTenants[domain]; !ok {
		if !options.validTenant(domain) {
			options.Log.Error("Tenant " + domain + " is not valid")
			options.validTenants[domain] = false
			options.Mutex.Unlock()
			return false
		}
		options.Log.Info("Tenant " + domain + " is valid")
		options.validTenants[domain] = true
	} else if !domainValid {
		// If the domain was not valid, skip the email
		options.Log.Debug("Tenant " + domain + " already checked and was not valid")
		options.Mutex.Unlock()
		return false

	}
	options.Mutex.Unlock()
	return true
}
