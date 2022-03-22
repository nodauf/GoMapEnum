package smtp

import (
	"GoMapEnum/src/utils"

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
	client, err := smtp.Dial(options.Target + ":25")
	if err != nil {
		options.Log.Error("Failed to establish a connection " + err.Error())
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
