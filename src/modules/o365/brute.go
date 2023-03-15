package o365

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
	"errors"
	"reflect"
)

// PrepareOptions is called before checking if the users are valid. It update the logging options to avoid printing the success
func PrepareOptions(optionsInterface *interface{}) interface{} {
	options := (*optionsInterface).(*Options)

	var optionsEnum = new(Options)
	*optionsEnum = *options
	var tmpLogger logger.Logger
	// Use office for enumeration (safer)
	optionsEnum.Mode = "office"
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

// Authenticate will be called to test an authentication and use the specified mode and check the lockout
func Authenticate(optionsInterface *interface{}, email, password string) bool {
	options := (*optionsInterface).(*Options)
	options.Mode = "oauth2"
	var valid bool
	var err error
	switch options.Mode {
	case "oauth2":
		valid, err = options.bruteOauth2(email, password)
		if err != nil && errors.Is(utils.ErrLockout, err) && options.StopOnLockout {
			options.Log.Fatal("The account %s is locked", email)
		}

	case "autodiscover":
		// We have no information on account lockout
		valid = false
		options.Log.Fatal("Mode autodiscover not implemented yet")
		//valid = bruteAutodiscover(email, password)
	}

	return valid
}
