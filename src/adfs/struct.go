package adfs

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
)

var log *logger.Logger

type Options struct {
	Domain string
	utils.BaseOptions
}

type userRealm struct {
	AuthNForwardType        int64  `json:"AuthNForwardType"`
	AuthURL                 string `json:"AuthURL"`
	CloudInstanceIssuerURI  string `json:"CloudInstanceIssuerUri"`
	CloudInstanceName       string `json:"CloudInstanceName"`
	DomainName              string `json:"DomainName"`
	FederationBrandName     string `json:"FederationBrandName"`
	FederationGlobalVersion int64  `json:"FederationGlobalVersion"`
	Login                   string `json:"Login"`
	NameSpaceType           string `json:"NameSpaceType"`
	State                   int64  `json:"State"`
	UserState               int64  `json:"UserState"`
}
