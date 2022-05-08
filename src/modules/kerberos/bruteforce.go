package kerberos

import (
	"GoMapEnum/src/logger"
	"reflect"
	"strings"
)

func Authenticate(optionsInterface *interface{}, username, password string) bool {
	options := (*optionsInterface).(*Options)
	valid, client, err := options.authenticate(username, password)
	defer client.Destroy()
	if err != nil {
		ok, errorString := handleKerbError(err)
		if strings.Contains(errorString, "LOCKED OUT") && options.StopOnLockout {
			options.Log.Fatal("The user %s has been locked out. Abort the bruteforce", username)
		}
		if ok {
			options.Log.Error("%s - %s", username, errorString)
		} else {
			options.Log.Fatal("%s - %s", username, errorString)
		}
	}
	return valid

}

// PrepareOptions is called before checking if the users are valid. It update the logging options to avoid printing the success
func PrepareOptions(optionsInterface *interface{}) interface{} {
	options := (*optionsInterface).(*Options)

	var optionsEnum = new(Options)
	*optionsEnum = *options
	var tmpLogger logger.Logger
	optionsEnum.Log = &tmpLogger
	*optionsEnum.Log = *options.Log
	optionsEnum.Log.Mode = "Enumeration"
	// If debug or verbose use this level in userenum module otherwise do not show the valid user

	if options.Log.Level == logger.DebugLevel || options.Log.Level == logger.VerboseLevel {
		optionsEnum.Log.Level = options.Log.Level
	} else {
		optionsEnum.Log.Level = logger.ErrorLevel
	}
	return reflect.ValueOf(optionsEnum).Interface()
}
