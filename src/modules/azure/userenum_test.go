package azure

import (
	"GoMapEnum/src/logger"
	"reflect"
	"testing"
)

func TestUserEnum(t *testing.T) {
	var results = make(map[string]bool)
	results["nodauf@gomapenum.onmicrosoft.com"] = true
	results["notExist@gomapenum.onmicrosoft.com"] = false
	results["notExist@tenantNotFound.com"] = false

	options := Options{}
	log := logger.New("User enumeration", "Azure", "https://autologon.microsoftazuread-sso.com")
	log.SetLevel(logger.FatalLevel)
	options.Log = log
	optionsInterface := reflect.ValueOf(&options).Interface()

	for username, wantedResults := range results {
		ok := UserEnum(&optionsInterface, username)
		if ok != wantedResults {
			t.Errorf("User enumeration for %s returned %t and was expected %t", username, ok, wantedResults)
		}
	}
}
