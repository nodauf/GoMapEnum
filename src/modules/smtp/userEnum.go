package smtp

import (
	"GoMapEnum/src/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	smtp "github.com/nodauf/net-smtp"
)

func PrepareSMTPConnections(optionsInterface *interface{}) {
	options := (*optionsInterface).(*Options)
	options.connectionsPool = make(chan *smtp.Client, options.Thread)
	var nbConnectionsRequired int
	nbConnectionsRequired = options.Thread
	if len(options.UsernameList) < options.Thread {
		nbConnectionsRequired = len(options.UsernameList)
	}
	options.Log.Debug("Preparing a pool of " + strconv.Itoa(nbConnectionsRequired) + " connections")
	for i := 1; i <= nbConnectionsRequired; i++ {
		client, err := smtp.Dial(options.Target + ":25")
		if err != nil {
			options.Log.Error("Failed to establish a connection " + err.Error())
			continue
		}
		err = client.Hello(utils.RandomString(6))
		if err != nil {
			fmt.Println("hello" + err.Error())
		}
		err = client.Mail(utils.RandomString(6) + "@" + options.Domain)
		if err != nil {
			fmt.Println("mail" + err.Error())
		}
		options.connectionsPool <- client
	}
}

func UserEnum(optionsInterface *interface{}, username string) bool {
	options := (*optionsInterface).(*Options)
	valid := false
	smtpConnection := <-options.connectionsPool
	switch strings.ToLower(options.Mode) {
	case "rcpt", "":
		err := smtpConnection.Rcpt(username)
		if err == nil {
			options.Log.Success(username)
			valid = true
		} else {
			options.Log.Debug(username + " => " + err.Error())
			options.Log.Fail(username)
		}
	case "vrfy":
		err := smtpConnection.Verify(username)
		if err == nil {
			options.Log.Success(username)
			valid = true
		} else {
			options.Log.Debug(username + " => " + err.Error())
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
			options.Log.Fail(username)
			// If the command is not implemented no need to pursue
			if code == "502" {
				CloseSMTPConnections(optionsInterface)
				options.Log.Fatal("The command is not implemented. No need to pursue using this method.")
			}
			fmt.Println(code)
		}
	default:
		CloseSMTPConnections(optionsInterface)
		options.Log.Fatal("Unrecognised mode: " + options.Mode + ". Only RCPT, VRFY and EXPN are supported.")
	}

	options.connectionsPool <- smtpConnection
	return valid
}

func CloseSMTPConnections(optionsInterface *interface{}) {
	options := (*optionsInterface).(*Options)
	options.Log.Debug("Closing the pool of connections")
	for i := 1; i <= options.Thread; i++ {
		select {
		case smtpConnection := <-options.connectionsPool:
			smtpConnection.Close()
		case <-time.After(1 * time.Second):
			options.Log.Debug("Something went wrong0 A connection seems already closed")
		}

	}
	close(options.connectionsPool)
}
