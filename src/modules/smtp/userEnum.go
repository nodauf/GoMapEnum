package smtp

import (
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	smtp "github.com/nodauf/net-smtp"
)

func PrepareSMTPConnections(optionsInterface *interface{}) {
	options := (*optionsInterface).(*Options)
	options.connectionsPool = make(chan *smtp.Client, options.Thread)

	if options.Target == "" {
		mxrecords, err := net.LookupMX(options.Domain)
		if err != nil {
			options.Log.Fatal("Not able to retrieve the MX for the domain " + options.Domain)
		}
		options.Target = strings.TrimRight(mxrecords[0].Host, ".")
	}
	options.Log.Target = options.Target

	var nbConnectionsRequired int
	nbConnectionsRequired = options.Thread
	if (options.Mode != "" && len(options.UsernameList) < options.Thread) || (options.Mode == "" && len(options.UsernameList)*3 < options.Thread) {
		nbConnectionsRequired = len(options.UsernameList)
	}
	options.Log.Debug("Preparing a pool of " + strconv.Itoa(nbConnectionsRequired) + " connections")
	for i := 1; i <= nbConnectionsRequired; i++ {
		client := options.createNewConnection()
		if client != nil {
			options.connectionsPool <- client
		}
	}
}

func UserEnum(optionsInterface *interface{}, username string) bool {
	options := (*optionsInterface).(*Options)
	valid := false
	smtpConnection := <-options.connectionsPool
	smtpConnection.Reset()
	err := options.prepareOneConnection(smtpConnection)
	if err != nil && strings.Contains(err.Error(), "connection reset by peer") {
		options.Log.Debug("Connection reset. Generating new one")
		smtpConnection = options.createNewConnection()
		err = options.prepareOneConnection(smtpConnection)
	}
	if err != nil {
		options.Log.Fatal("Failed to prepare a connection. " + err.Error())
	}
	switch strings.ToLower(options.Mode) {
	case "rcpt":
		err := smtpConnection.Rcpt(username)
		if err == nil {
			options.Log.Success(username)
			valid = true
		} else {
			options.Log.Debug(username + " => " + err.Error())
			if strings.Contains(err.Error(), "connection reset by peer") {
				smtpConnection.Close()
				options.createNewConnection()
				return UserEnum(optionsInterface, username)
			}
			options.Log.Fail(username)
		}
	case "vrfy":
		err := smtpConnection.Verify(username)
		if err == nil {
			options.Log.Success(username)
			valid = true
		} else {
			options.Log.Debug(username + " => " + err.Error())
			if strings.Contains(err.Error(), "connection reset by peer") {
				smtpConnection.Close()
				options.createNewConnection()
				return UserEnum(optionsInterface, username)
			}
			options.Log.Fail(username)
		}
	case "expn":
		err := smtpConnection.Expand(username)
		if err == nil {
			options.Log.Success(username)
			valid = true
		} else {
			code := strings.Split(err.Error(), " ")[0]
			options.Log.Debug(username + " => " + err.Error())
			if strings.Contains(err.Error(), "connection reset by peer") {
				smtpConnection.Close()
				options.createNewConnection()
				return UserEnum(optionsInterface, username)
			}
			options.Log.Fail(username)
			// If the command is not implemented no need to pursue
			if code == "502" && !options.all {
				CloseSMTPConnections(optionsInterface)
				options.Log.Fatal("The command EXPN is not implemented. No need to pursue using this method.")
			}
			if code == "502" && options.all {
				options.expnNotRecognized = true
			}
		}
	case "", "all":

		optionsCopy := *options
		options.connectionsPool <- smtpConnection
		// Execute the 3 enumeration methods
		optionsCopy.all = true
		// RCPT request
		options.Log.Debug("No enumeration method specify. Executing enumeration with RCPT, VRFY, EXPN and ALL")
		options.Log.Debug("Enumerate with RCPT")
		optionsCopy.Mode = "rcpt"
		newOptionsInterface := reflect.ValueOf(&optionsCopy).Interface()
		valid = UserEnum(&newOptionsInterface, username)
		if valid {
			return true
		}
		// VRFY
		options.Log.Debug("Enumerate with VRFY")
		optionsCopy.Mode = "vrfy"
		newOptionsInterface = reflect.ValueOf(&optionsCopy).Interface()
		valid = UserEnum(&newOptionsInterface, username)
		if valid {
			return true
		}
		// EXPN
		if !options.expnNotRecognized {
			options.Log.Debug("Enumerate with EXPN")
			optionsCopy.Mode = "expn"
			newOptionsInterface = reflect.ValueOf(&optionsCopy).Interface()
			valid = UserEnum(&newOptionsInterface, username)
		}
		return valid
	default:
		options.Log.Fatal("Unrecognised mode: " + options.Mode + ". Only RCPT, VRFY and EXPN are supported.")
	}
	options.connectionsPool <- smtpConnection
	return valid
}

func CloseSMTPConnections(optionsInterface *interface{}) {
	options := (*optionsInterface).(*Options)
	options.Log.Debug("Closing the pool of connections")
	for i := 1; i <= len(options.connectionsPool); i++ {
		select {
		case smtpConnection := <-options.connectionsPool:
			smtpConnection.Close()
		case <-time.After(1 * time.Second):
			options.Log.Debug("Something went wrong. A connection seems already closed")
		}

	}
	close(options.connectionsPool)
}
