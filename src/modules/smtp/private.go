package smtp

import (
	"GoMapEnum/src/utils"
	"fmt"
	"net"
	"time"

	smtp "github.com/nodauf/net-smtp"
)

func (options *Options) prepareOneConnection(client *smtp.Client) error {

	mailFrom := utils.RandomString(6) + "@" + utils.RandomString(4) + "." + utils.RandomString(2)
	err := client.Mail(mailFrom)
	if err != nil {
		//options.Log.Error("Mail From command failed with email:  " + mailFrom + " and error " + err.Error())
		return err
	}
	options.Log.Debug("One connection has been successfully prepared")

	return nil
}

func (options *Options) createNewConnection() *smtp.Client {
	var conn net.Conn
	var err error
	if options.ProxyTCP != nil {
		conn, err = options.ProxyTCP.Dial("tcp", fmt.Sprintf("%s:%d", options.Target, 25))
	} else {
		defaultDailer := &net.Dialer{Timeout: time.Duration(options.Timeout * int(time.Second))}
		conn, err = defaultDailer.Dial("tcp", fmt.Sprintf("%s:%d", options.Target, 25))
	}

	if err != nil {
		options.Log.Error("Failed to establish a connection " + err.Error())
		return nil
	}

	client, err := smtp.NewClient(conn, options.Target)
	if err != nil {
		options.Log.Error("Failed to open the SMTP connection " + err.Error())
		return nil
	}

	hello := utils.RandomString(6)
	err = client.Hello(hello)
	if err != nil {
		options.Log.Error("helo command failed with helo " + hello + " and error " + err.Error())
		return nil
	}
	return client
}
