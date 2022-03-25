package o365

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
	"errors"
	"reflect"
	"strconv"
)

func PrepareOptions(optionsInterface *interface{}) interface{} {
	options := (*optionsInterface).(*Options)

	options.Log.Debug("Validating the users")
	var optionsEnum = new(Options)
	*optionsEnum = *options
	var tmpLogger logger.Logger
	// Use office for enumeration (safer)
	optionsEnum.Mode = "office"
	optionsEnum.Log = &tmpLogger
	*optionsEnum.Log = *options.Log
	optionsEnum.Log.Type = "Enumeration"
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
	var valid bool
	var err error
	switch options.Mode {
	case "oauth2":
		valid, err = options.bruteOauth2(email, password)
		if err != nil && errors.Is(utils.ErrLockout, err) {
			options.Log.Error("The account %s is locked", email)
			options.lockoutCounter++
		}
		// Fail safe to avoid locking to many account
		if options.lockoutCounter >= options.LockoutThreshold {
			options.Log.Fatal("Too many lockout: " + strconv.Itoa(options.lockoutCounter) + " >= " + strconv.Itoa(options.LockoutThreshold))
		}

	case "autodiscover":
		// We have no information on account lockout
		valid = bruteAutodiscover(email, password)
	}

	return valid
}
