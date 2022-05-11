package kerberos

func UserEnum(optionsInterface *interface{}, username string) bool {
	valid := false
	options := (*optionsInterface).(*Options)
	valid, err := options.testUsername(username)
	if valid {
		options.Log.Success(username)
	}
	if err != nil {
		// This is to determine if the error is "okay" or if we should abort everything
		options.Log.Debug(err.Error())
		ok, errorString := handleKerbError(err)
		if ok {
			options.Log.Debug("%s - %s", username, errorString)
		} else {
			options.Log.Fatal("%s - %s", username, errorString)
		}
	} else {
		options.Log.Debug("Unknown behavior for username %s", username)
	}

	return valid
}
