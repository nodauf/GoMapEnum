package azure

import (
	"GoMapEnum/src/utils"
)

// Options for Azure module
type Options struct {
	utils.BaseOptions
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}

type azureResponse struct {
	Body struct {
		Fault struct {
			Detail struct {
				Text  string `xml:",chardata"`
				Error struct {
					Text          string `xml:",chardata"`
					Psf           string `xml:"psf,attr"`
					Value         string `xml:"value"`
					Internalerror struct {
						Chardata string `xml:",chardata"`
						Code     string `xml:"code"`
						Text     string `xml:"text"`
					} `xml:"internalerror"`
				} `xml:"error"`
			} `xml:"Detail"`
		} `xml:"Fault"`
	} `xml:"Body"`
}
