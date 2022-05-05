package smtp

import (
	"GoMapEnum/src/utils"
	"net"

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

	conn, err = utils.OpenConnectionWoProxy(options.Target, "25", options.Timeout, options.ProxyTCP)
	if err != nil {
		options.Log.Error("cannot open a connection to %s:%d : %v", options.Target, 25, err)
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
