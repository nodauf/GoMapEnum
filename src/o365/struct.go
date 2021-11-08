package o365

import (
	"GoMapEnum/src/utils"
	"encoding/xml"
)

// Options for o365 module
type Options struct {
	Mode             string
	LockoutThreshold int
	utils.BaseOptions
}

type officeData struct {
	OriginalRequest                string `json:"originalRequest"`
	IsOtherIdpSupported            bool   `json:"isOtherIdpSupported"`
	IsRemoteNGCSupported           bool   `json:"isRemoteNGCSupported"`
	IsAccessPassSupported          bool   `json:"isAccessPassSupported"`
	CheckPhones                    bool   `json:"checkPhones"`
	IsCookieBannerShown            bool   `json:"isCookieBannerShown"`
	IsFidoSupported                bool   `json:"isFidoSupported"`
	Forceotclogin                  bool   `json:"forceotclogin"`
	IsExternalFederationDisallowed bool   `json:"isExternalFederationDisallowed"`
	IsRemoteConnectSupported       bool   `json:"isRemoteConnectSupported"`
	IsSignup                       bool   `json:"isSignup"`
	FederationFlags                int    `json:"federationFlags"`
	Username                       string `json:"username"`
}

type officeResponse struct {
	Credentials struct {
		CertAuthParams  interface{} `json:"CertAuthParams"`
		FacebookParams  interface{} `json:"FacebookParams"`
		FidoParams      interface{} `json:"FidoParams"`
		GoogleParams    interface{} `json:"GoogleParams"`
		HasPassword     bool        `json:"HasPassword"`
		PrefCredential  int64       `json:"PrefCredential"`
		RemoteNgcParams interface{} `json:"RemoteNgcParams"`
		SasParams       interface{} `json:"SasParams"`
	} `json:"Credentials"`
	Display        string `json:"Display"`
	EstsProperties struct {
		DomainType         int64       `json:"DomainType"`
		UserTenantBranding interface{} `json:"UserTenantBranding"`
		DesktopSsoEnabled  *bool       `json:"DesktopSsoEnabled,omitempty"`
	} `json:"EstsProperties"`
	IfExistsResult     int64  `json:"IfExistsResult"`
	IsSignupDisallowed bool   `json:"IsSignupDisallowed"`
	IsUnmanaged        bool   `json:"IsUnmanaged"`
	ThrottleStatus     int64  `json:"ThrottleStatus"`
	Username           string `json:"Username"`
	APICanary          string `json:"apiCanary"`
}

type realmInfo struct {
	XMLName                xml.Name `xml:"RealmInfo"`
	Text                   string   `xml:",chardata"`
	Success                string   `xml:"Success,attr"`
	Script                 string   `xml:"script"`
	State                  string   `xml:"State"`
	UserState              string   `xml:"UserState"`
	Login                  string   `xml:"Login"`
	NameSpaceType          string   `xml:"NameSpaceType"`
	DomainName             string   `xml:"DomainName"`
	IsFederatedNS          string   `xml:"IsFederatedNS"`
	FederationBrandName    string   `xml:"FederationBrandName"`
	CloudInstanceName      string   `xml:"CloudInstanceName"`
	CloudInstanceIssuerUri string   `xml:"CloudInstanceIssuerUri"`
}

type oauth2Data struct {
	ClientID  string `form:"client_id"`
	GrantType string `form:"grant_type"`
	Password  string `form:"password"`
	Resource  string `form:"resource"`
	Scope     string `form:"scope"`
	Username  string `form:"username"`
}

type oauth2Output struct {
	CorrelationID    string  `json:"correlation_id"`
	Error            string  `json:"error"`
	ErrorCodes       []int64 `json:"error_codes"`
	ErrorDescription string  `json:"error_description"`
	ErrorURI         string  `json:"error_uri"`
	Timestamp        string  `json:"timestamp"`
	TraceID          string  `json:"trace_id"`
	AccessToken      string  `json:"access_token"`
	ExpiresIn        string  `json:"expires_in"`
	ExpiresOn        string  `json:"expires_on"`
	ExtExpiresIn     string  `json:"ext_expires_in"`
	IDToken          string  `json:"id_token"`
	NotBefore        string  `json:"not_before"`
	RefreshToken     string  `json:"refresh_token"`
	Resource         string  `json:"resource"`
	Scope            string  `json:"scope"`
	TokenType        string  `json:"token_type"`
}
