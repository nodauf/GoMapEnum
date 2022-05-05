package o365

import (
	"GoMapEnum/src/utils"
	"encoding/xml"
)

// Options for o365 module
type Options struct {
	Mode        string
	DumpObjects string
	HTML        bool
	JSON        bool

	validTenants map[string]bool
	utils.BaseOptions
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
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

/* Structures used to dump azure data */

type Application struct {
	Odata_metadata string `json:"odata.metadata"`
	Value          []struct {
		AddIns                        []interface{} `json:"addIns"`
		AllowActAsForAllClients       interface{}   `json:"allowActAsForAllClients"`
		AllowPassthroughUsers         interface{}   `json:"allowPassthroughUsers"`
		AppBranding                   interface{}   `json:"appBranding"`
		AppCategory                   interface{}   `json:"appCategory"`
		AppData                       interface{}   `json:"appData"`
		AppID                         string        `json:"appId"`
		AppMetadata                   interface{}   `json:"appMetadata"`
		AppRoles                      []interface{} `json:"appRoles"`
		ApplicationTemplateID         interface{}   `json:"applicationTemplateId"`
		AvailableToOtherTenants       bool          `json:"availableToOtherTenants"`
		Certification                 interface{}   `json:"certification"`
		DeletionTimestamp             interface{}   `json:"deletionTimestamp"`
		DisabledByMicrosoftStatus     interface{}   `json:"disabledByMicrosoftStatus"`
		DisplayName                   string        `json:"displayName"`
		EncryptedMsiApplicationSecret interface{}   `json:"encryptedMsiApplicationSecret"`
		ErrorURL                      interface{}   `json:"errorUrl"`
		GroupMembershipClaims         interface{}   `json:"groupMembershipClaims"`
		Homepage                      string        `json:"homepage"`
		IdentifierUris                []string      `json:"identifierUris"`
		InformationalUrls             struct {
			Marketing      interface{} `json:"marketing"`
			Privacy        interface{} `json:"privacy"`
			Support        interface{} `json:"support"`
			TermsOfService interface{} `json:"termsOfService"`
		} `json:"informationalUrls"`
		IsDeviceOnlyAuthSupported      interface{}   `json:"isDeviceOnlyAuthSupported"`
		KeyCredentials                 []interface{} `json:"keyCredentials"`
		KnownClientApplications        []interface{} `json:"knownClientApplications"`
		LogoURL                        interface{}   `json:"logoUrl"`
		LogoutURL                      interface{}   `json:"logoutUrl"`
		Oauth2AllowIDTokenImplicitFlow bool          `json:"oauth2AllowIdTokenImplicitFlow"`
		Oauth2AllowImplicitFlow        bool          `json:"oauth2AllowImplicitFlow"`
		Oauth2AllowURLPathMatching     bool          `json:"oauth2AllowUrlPathMatching"`
		Oauth2Permissions              []struct {
			AdminConsentDescription string      `json:"adminConsentDescription"`
			AdminConsentDisplayName string      `json:"adminConsentDisplayName"`
			ID                      string      `json:"id"`
			IsEnabled               bool        `json:"isEnabled"`
			Lang                    interface{} `json:"lang"`
			Origin                  string      `json:"origin"`
			Type                    string      `json:"type"`
			UserConsentDescription  string      `json:"userConsentDescription"`
			UserConsentDisplayName  string      `json:"userConsentDisplayName"`
			Value                   string      `json:"value"`
		} `json:"oauth2Permissions"`
		Oauth2RequirePostResponse bool        `json:"oauth2RequirePostResponse"`
		ObjectID                  string      `json:"objectId"`
		ObjectType                string      `json:"objectType"`
		Odata_type                string      `json:"odata.type"`
		OptionalClaims            interface{} `json:"optionalClaims"`
		ParentalControlSettings   struct {
			CountriesBlockedForMinors []interface{} `json:"countriesBlockedForMinors"`
			LegalAgeGroupRule         string        `json:"legalAgeGroupRule"`
		} `json:"parentalControlSettings"`
		PasswordCredentials []struct {
			CustomKeyIdentifier interface{} `json:"customKeyIdentifier"`
			EndDate             string      `json:"endDate"`
			KeyID               string      `json:"keyId"`
			StartDate           string      `json:"startDate"`
			Value               interface{} `json:"value"`
		} `json:"passwordCredentials"`
		PublicClient            bool        `json:"publicClient"`
		PublisherDomain         string      `json:"publisherDomain"`
		RecordConsentConditions interface{} `json:"recordConsentConditions"`
		ReplyUrls               []string    `json:"replyUrls"`
		RequiredResourceAccess  []struct {
			ResourceAccess []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"resourceAccess"`
			ResourceAppID string `json:"resourceAppId"`
		} `json:"requiredResourceAccess"`
		SamlMetadataURL            interface{}   `json:"samlMetadataUrl"`
		SupportsConvergence        bool          `json:"supportsConvergence"`
		TokenEncryptionKeyID       interface{}   `json:"tokenEncryptionKeyId"`
		TrustedCertificateSubjects []interface{} `json:"trustedCertificateSubjects"`
		VerifiedPublisher          struct {
			AddedDateTime       interface{} `json:"addedDateTime"`
			DisplayName         interface{} `json:"displayName"`
			VerifiedPublisherID interface{} `json:"verifiedPublisherId"`
		} `json:"verifiedPublisher"`
	} `json:"value"`
}

type Devices struct {
	Odata_metadata string `json:"odata.metadata"`
	Odata_nextLink string `json:"odata.nextLink"`
	Value          []struct {
		AccountEnabled         bool `json:"accountEnabled"`
		AlternativeSecurityIds []struct {
			IdentityProvider interface{} `json:"identityProvider"`
			Key              string      `json:"key"`
			Type             int64       `json:"type"`
		} `json:"alternativeSecurityIds"`
		ApproximateLastLogonTimestamp string        `json:"approximateLastLogonTimestamp"`
		BitLockerKey                  []interface{} `json:"bitLockerKey"`
		Capabilities                  []interface{} `json:"capabilities"`
		ComplianceExpiryTime          interface{}   `json:"complianceExpiryTime"`
		CompliantApplications         []interface{} `json:"compliantApplications"`
		CompliantAppsManagementAppID  interface{}   `json:"compliantAppsManagementAppId"`
		DeletionTimestamp             interface{}   `json:"deletionTimestamp"`
		DeviceCategory                interface{}   `json:"deviceCategory"`
		DeviceID                      string        `json:"deviceId"`
		DeviceKey                     []struct {
			CreationTime         string `json:"creationTime"`
			CustomKeyInformation string `json:"customKeyInformation"`
			KeyIdentifier        string `json:"keyIdentifier"`
			KeyMaterial          string `json:"keyMaterial"`
			Usage                string `json:"usage"`
		} `json:"deviceKey"`
		DeviceManagementAppID interface{} `json:"deviceManagementAppId"`
		DeviceManufacturer    interface{} `json:"deviceManufacturer"`
		DeviceMetadata        interface{} `json:"deviceMetadata"`
		DeviceModel           interface{} `json:"deviceModel"`
		DeviceOSType          string      `json:"deviceOSType"`
		DeviceOSVersion       string      `json:"deviceOSVersion"`
		DeviceObjectVersion   int64       `json:"deviceObjectVersion"`
		DeviceOwnership       interface{} `json:"deviceOwnership"`
		DevicePhysicalIds     []string    `json:"devicePhysicalIds"`
		DeviceSystemMetadata  []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"deviceSystemMetadata"`
		DeviceTrustType              string        `json:"deviceTrustType"`
		DirSyncEnabled               interface{}   `json:"dirSyncEnabled"`
		DisplayName                  string        `json:"displayName"`
		DomainName                   interface{}   `json:"domainName"`
		EnrollmentProfileName        interface{}   `json:"enrollmentProfileName"`
		EnrollmentType               interface{}   `json:"enrollmentType"`
		ExchangeActiveSyncID         []interface{} `json:"exchangeActiveSyncId"`
		ExternalSourceName           interface{}   `json:"externalSourceName"`
		Hostnames                    []interface{} `json:"hostnames"`
		IsCompliant                  interface{}   `json:"isCompliant"`
		IsManaged                    interface{}   `json:"isManaged"`
		IsRooted                     interface{}   `json:"isRooted"`
		KeyCredentials               []interface{} `json:"keyCredentials"`
		LastDirSyncTime              interface{}   `json:"lastDirSyncTime"`
		LocalCredentials             interface{}   `json:"localCredentials"`
		ManagementType               interface{}   `json:"managementType"`
		ObjectID                     string        `json:"objectId"`
		ObjectType                   string        `json:"objectType"`
		Odata_type                   string        `json:"odata.type"`
		OnPremisesSecurityIdentifier interface{}   `json:"onPremisesSecurityIdentifier"`
		OrganizationalUnit           interface{}   `json:"organizationalUnit"`
		ProfileType                  string        `json:"profileType"`
		Reserved1                    interface{}   `json:"reserved1"`
		SourceType                   interface{}   `json:"sourceType"`
		SystemLabels                 []interface{} `json:"systemLabels"`
	} `json:"value"`
}

type DirectoryRoles struct {
	Odata_metadata string `json:"odata.metadata"`
	Value          []struct {
		CloudSecurityIdentifier string      `json:"cloudSecurityIdentifier"`
		DeletionTimestamp       interface{} `json:"deletionTimestamp"`
		Description             string      `json:"description"`
		DisplayName             string      `json:"displayName"`
		IsSystem                bool        `json:"isSystem"`
		ObjectID                string      `json:"objectId"`
		ObjectType              string      `json:"objectType"`
		Odata_type              string      `json:"odata.type"`
		RoleDisabled            bool        `json:"roleDisabled"`
		RoleTemplateID          string      `json:"roleTemplateId"`
	} `json:"value"`
}

type Groups struct {
	Odata_metadata string `json:"odata.metadata"`
	Odata_nextLink string `json:"odata.nextLink"`
	Value          []struct {
		AppMetadata struct {
			Data    []interface{} `json:"data"`
			Version int64         `json:"version"`
		} `json:"appMetadata"`
		Classification                interface{}   `json:"classification"`
		CloudSecurityIdentifier       string        `json:"cloudSecurityIdentifier"`
		CreatedByAppID                string        `json:"createdByAppId"`
		CreatedDateTime               string        `json:"createdDateTime"`
		CreationOptions               []string      `json:"creationOptions"`
		DeletionTimestamp             interface{}   `json:"deletionTimestamp"`
		Description                   string        `json:"description"`
		DirSyncEnabled                bool          `json:"dirSyncEnabled"`
		DisplayName                   string        `json:"displayName"`
		ExchangeResources             []string      `json:"exchangeResources"`
		ExpirationDateTime            interface{}   `json:"expirationDateTime"`
		ExternalGroupIds              []interface{} `json:"externalGroupIds"`
		ExternalGroupProviderID       interface{}   `json:"externalGroupProviderId"`
		ExternalGroupState            interface{}   `json:"externalGroupState"`
		GroupTypes                    []string      `json:"groupTypes"`
		InfoCatalogs                  []interface{} `json:"infoCatalogs"`
		IsAssignableToRole            interface{}   `json:"isAssignableToRole"`
		IsMembershipRuleLocked        interface{}   `json:"isMembershipRuleLocked"`
		IsPublic                      bool          `json:"isPublic"`
		LastDirSyncTime               string        `json:"lastDirSyncTime"`
		LicenseAssignment             []interface{} `json:"licenseAssignment"`
		Mail                          string        `json:"mail"`
		MailEnabled                   bool          `json:"mailEnabled"`
		MailNickname                  string        `json:"mailNickname"`
		MembershipRule                interface{}   `json:"membershipRule"`
		MembershipRuleProcessingState interface{}   `json:"membershipRuleProcessingState"`
		MembershipTypes               []interface{} `json:"membershipTypes"`
		ObjectID                      string        `json:"objectId"`
		ObjectType                    string        `json:"objectType"`
		Odata_type                    string        `json:"odata.type"`
		OnPremisesSecurityIdentifier  string        `json:"onPremisesSecurityIdentifier"`
		PreferredDataLocation         interface{}   `json:"preferredDataLocation"`
		PreferredLanguage             interface{}   `json:"preferredLanguage"`
		PrimarySMTPAddress            string        `json:"primarySMTPAddress"`
		ProvisioningErrors            []interface{} `json:"provisioningErrors"`
		ProxyAddresses                []string      `json:"proxyAddresses"`
		RenewedDateTime               string        `json:"renewedDateTime"`
		ResourceBehaviorOptions       []string      `json:"resourceBehaviorOptions"`
		ResourceProvisioningOptions   []string      `json:"resourceProvisioningOptions"`
		SecurityEnabled               bool          `json:"securityEnabled"`
		SharepointResources           []string      `json:"sharepointResources"`
		TargetAddress                 string        `json:"targetAddress"`
		Theme                         interface{}   `json:"theme"`
		Visibility                    string        `json:"visibility"`
		WellKnownObject               interface{}   `json:"wellKnownObject"`
	} `json:"value"`
}

type Oauth2PermissionGrants struct {
	Odata_metadata string `json:"odata.metadata"`
	Value          []struct {
		ClientID    string `json:"clientId"`
		ConsentType string `json:"consentType"`
		ExpiryTime  string `json:"expiryTime"`
		ObjectID    string `json:"objectId"`
		PrincipalID string `json:"principalId"`
		ResourceID  string `json:"resourceId"`
		Scope       string `json:"scope"`
		StartTime   string `json:"startTime"`
	} `json:"value"`
}

type Policies struct {
	Odata_metadata string `json:"odata.metadata"`
	Value          []struct {
		DeletionTimestamp   interface{}   `json:"deletionTimestamp"`
		DisplayName         string        `json:"displayName"`
		KeyCredentials      []interface{} `json:"keyCredentials"`
		ObjectID            string        `json:"objectId"`
		ObjectType          string        `json:"objectType"`
		Odata_type          string        `json:"odata.type"`
		PolicyDetail        []string      `json:"policyDetail"`
		PolicyIdentifier    interface{}   `json:"policyIdentifier"`
		PolicyType          int64         `json:"policyType"`
		TenantDefaultPolicy int64         `json:"tenantDefaultPolicy"`
	} `json:"value"`
}

type RoleDefinitions struct {
	Odata_metadata string `json:"odata.metadata"`
	Value          []struct {
		DeletionTimestamp       interface{} `json:"deletionTimestamp"`
		Description             string      `json:"description"`
		DisplayName             string      `json:"displayName"`
		InheritsPermissionsFrom []struct {
			ObjectID   string `json:"objectId"`
			Odata_type string `json:"odata.type"`
		} `json:"inheritsPermissionsFrom"`
		InheritsPermissionsFrom_odata_navigationLinkURL string   `json:"inheritsPermissionsFrom@odata.navigationLinkUrl"`
		IsBuiltIn                                       bool     `json:"isBuiltIn"`
		IsEnabled                                       bool     `json:"isEnabled"`
		ObjectID                                        string   `json:"objectId"`
		ObjectType                                      string   `json:"objectType"`
		Odata_type                                      string   `json:"odata.type"`
		ResourceScopes                                  []string `json:"resourceScopes"`
		RolePermissions                                 []struct {
			Condition       string `json:"condition"`
			ResourceActions struct {
				AllowedResourceActions []string `json:"allowedResourceActions"`
			} `json:"resourceActions"`
		} `json:"rolePermissions"`
		TemplateID string `json:"templateId"`
		Version    string `json:"version"`
	} `json:"value"`
}

type ServicePrincipals struct {
	Odata_metadata string `json:"odata.metadata"`
	Odata_nextLink string `json:"odata.nextLink"`
	Value          []struct {
		AccountEnabled   bool          `json:"accountEnabled"`
		AddIns           []interface{} `json:"addIns"`
		AlternativeNames []interface{} `json:"alternativeNames"`
		AppBranding      interface{}   `json:"appBranding"`
		AppCategory      interface{}   `json:"appCategory"`
		AppData          interface{}   `json:"appData"`
		AppDisplayName   string        `json:"appDisplayName"`
		AppID            string        `json:"appId"`
		AppMetadata      struct {
			Data []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"data"`
			Version int64 `json:"version"`
		} `json:"appMetadata"`
		AppOwnerTenantID          string `json:"appOwnerTenantId"`
		AppRoleAssignmentRequired bool   `json:"appRoleAssignmentRequired"`
		AppRoles                  []struct {
			AllowedMemberTypes []string    `json:"allowedMemberTypes"`
			Description        string      `json:"description"`
			DisplayName        string      `json:"displayName"`
			ID                 string      `json:"id"`
			IsEnabled          bool        `json:"isEnabled"`
			Lang               interface{} `json:"lang"`
			Origin             string      `json:"origin"`
			Value              string      `json:"value"`
		} `json:"appRoles"`
		ApplicationTemplateID     interface{} `json:"applicationTemplateId"`
		AuthenticationPolicy      interface{} `json:"authenticationPolicy"`
		Certification             interface{} `json:"certification"`
		DeletionTimestamp         interface{} `json:"deletionTimestamp"`
		DisabledByMicrosoftStatus interface{} `json:"disabledByMicrosoftStatus"`
		DisplayName               string      `json:"displayName"`
		ErrorURL                  interface{} `json:"errorUrl"`
		Homepage                  string      `json:"homepage"`
		InformationalUrls         struct {
			Marketing      interface{} `json:"marketing"`
			Privacy        string      `json:"privacy"`
			Support        interface{} `json:"support"`
			TermsOfService string      `json:"termsOfService"`
		} `json:"informationalUrls"`
		KeyCredentials             []interface{} `json:"keyCredentials"`
		LogoutURL                  string        `json:"logoutUrl"`
		ManagedIdentityResourceID  interface{}   `json:"managedIdentityResourceId"`
		MicrosoftFirstParty        bool          `json:"microsoftFirstParty"`
		NotificationEmailAddresses []interface{} `json:"notificationEmailAddresses"`
		Oauth2Permissions          []struct {
			AdminConsentDescription string      `json:"adminConsentDescription"`
			AdminConsentDisplayName string      `json:"adminConsentDisplayName"`
			ID                      string      `json:"id"`
			IsEnabled               bool        `json:"isEnabled"`
			Lang                    interface{} `json:"lang"`
			Origin                  string      `json:"origin"`
			Type                    string      `json:"type"`
			UserConsentDescription  string      `json:"userConsentDescription"`
			UserConsentDisplayName  string      `json:"userConsentDisplayName"`
			Value                   string      `json:"value"`
		} `json:"oauth2Permissions"`
		ObjectID                            string        `json:"objectId"`
		ObjectType                          string        `json:"objectType"`
		Odata_type                          string        `json:"odata.type"`
		PasswordCredentials                 []interface{} `json:"passwordCredentials"`
		PreferredSingleSignOnMode           interface{}   `json:"preferredSingleSignOnMode"`
		PreferredTokenSigningKeyEndDateTime interface{}   `json:"preferredTokenSigningKeyEndDateTime"`
		PreferredTokenSigningKeyThumbprint  interface{}   `json:"preferredTokenSigningKeyThumbprint"`
		PublisherName                       string        `json:"publisherName"`
		ReplyUrls                           []string      `json:"replyUrls"`
		SamlMetadataURL                     interface{}   `json:"samlMetadataUrl"`
		SamlSingleSignOnSettings            interface{}   `json:"samlSingleSignOnSettings"`
		ServicePrincipalNames               []string      `json:"servicePrincipalNames"`
		ServicePrincipalType                string        `json:"servicePrincipalType"`
		Tags                                []string      `json:"tags"`
		TokenEncryptionKeyID                interface{}   `json:"tokenEncryptionKeyId"`
		UseCustomTokenSigningKey            interface{}   `json:"useCustomTokenSigningKey"`
		VerifiedPublisher                   struct {
			AddedDateTime       string `json:"addedDateTime"`
			DisplayName         string `json:"displayName"`
			VerifiedPublisherID string `json:"verifiedPublisherId"`
		} `json:"verifiedPublisher"`
	} `json:"value"`
}

type TenantDetails struct {
	Odata_metadata string `json:"odata.metadata"`
	Value          []struct {
		AssignedPlans []struct {
			AssignedTimestamp string `json:"assignedTimestamp"`
			CapabilityStatus  string `json:"capabilityStatus"`
			Service           string `json:"service"`
			ServicePlanID     string `json:"servicePlanId"`
		} `json:"assignedPlans"`
		AuthorizedServiceInstance                 []string      `json:"authorizedServiceInstance"`
		City                                      string        `json:"city"`
		CloudRtcUserPolicies                      interface{}   `json:"cloudRtcUserPolicies"`
		CompanyLastDirSyncTime                    string        `json:"companyLastDirSyncTime"`
		CompanyTags                               []string      `json:"companyTags"`
		CompassEnabled                            interface{}   `json:"compassEnabled"`
		Country                                   interface{}   `json:"country"`
		CountryLetterCode                         string        `json:"countryLetterCode"`
		CreatedDateTime                           string        `json:"createdDateTime"`
		DeletionTimestamp                         interface{}   `json:"deletionTimestamp"`
		DirSyncEnabled                            bool          `json:"dirSyncEnabled"`
		DisplayName                               string        `json:"displayName"`
		IsMultipleDataLocationsForServicesEnabled interface{}   `json:"isMultipleDataLocationsForServicesEnabled"`
		MarketingNotificationEmails               []interface{} `json:"marketingNotificationEmails"`
		ObjectID                                  string        `json:"objectId"`
		ObjectType                                string        `json:"objectType"`
		Odata_type                                string        `json:"odata.type"`
		PostalCode                                string        `json:"postalCode"`
		PreferredLanguage                         string        `json:"preferredLanguage"`
		PrivacyProfile                            interface{}   `json:"privacyProfile"`
		ProvisionedPlans                          []struct {
			CapabilityStatus   string `json:"capabilityStatus"`
			ProvisioningStatus string `json:"provisioningStatus"`
			Service            string `json:"service"`
		} `json:"provisionedPlans"`
		ProvisioningErrors                   []interface{} `json:"provisioningErrors"`
		ReleaseTrack                         interface{}   `json:"releaseTrack"`
		ReplicationScope                     string        `json:"replicationScope"`
		SecurityComplianceNotificationMails  []interface{} `json:"securityComplianceNotificationMails"`
		SecurityComplianceNotificationPhones []interface{} `json:"securityComplianceNotificationPhones"`
		SelfServePasswordResetPolicy         interface{}   `json:"selfServePasswordResetPolicy"`
		State                                string        `json:"state"`
		Street                               string        `json:"street"`
		TechnicalNotificationMails           []string      `json:"technicalNotificationMails"`
		TelephoneNumber                      interface{}   `json:"telephoneNumber"`
		TenantType                           interface{}   `json:"tenantType"`
		VerifiedDomains                      []struct {
			Capabilities string `json:"capabilities"`
			Default      bool   `json:"default"`
			ID           string `json:"id"`
			Initial      bool   `json:"initial"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"verifiedDomains"`
		WindowsCredentialsEncryptionCertificate interface{} `json:"windowsCredentialsEncryptionCertificate"`
	} `json:"value"`
}

type Users struct {
	Odata_metadata string `json:"odata.metadata"`
	Odata_nextLink string `json:"odata.nextLink"`
	Value          []struct {
		AcceptedAs             string      `json:"acceptedAs"`
		AcceptedOn             string      `json:"acceptedOn"`
		AccountEnabled         bool        `json:"accountEnabled"`
		AgeGroup               interface{} `json:"ageGroup"`
		AlternativeSecurityIds []struct {
			IdentityProvider string `json:"identityProvider"`
			Key              string `json:"key"`
			Type             int64  `json:"type"`
		} `json:"alternativeSecurityIds"`
		AppMetadata      interface{} `json:"appMetadata"`
		AssignedLicenses []struct {
			DisabledPlans []string `json:"disabledPlans"`
			SkuID         string   `json:"skuId"`
		} `json:"assignedLicenses"`
		AssignedPlans []struct {
			AssignedTimestamp string `json:"assignedTimestamp"`
			CapabilityStatus  string `json:"capabilityStatus"`
			Service           string `json:"service"`
			ServicePlanID     string `json:"servicePlanId"`
		} `json:"assignedPlans"`
		City                               interface{} `json:"city"`
		CloudAudioConferencingProviderInfo interface{} `json:"cloudAudioConferencingProviderInfo"`
		CloudMSExchRecipientDisplayType    int64       `json:"cloudMSExchRecipientDisplayType"`
		CloudMSRtcIsSipEnabled             bool        `json:"cloudMSRtcIsSipEnabled"`
		CloudMSRtcOwnerUrn                 interface{} `json:"cloudMSRtcOwnerUrn"`
		CloudMSRtcPolicyAssignments        []string    `json:"cloudMSRtcPolicyAssignments"`
		CloudMSRtcPool                     string      `json:"cloudMSRtcPool"`
		CloudMSRtcServiceAttributes        struct {
			ApplicationOptions   int64  `json:"applicationOptions"`
			DeploymentLocator    string `json:"deploymentLocator"`
			HideFromAddressLists bool   `json:"hideFromAddressLists"`
			OptionFlags          int64  `json:"optionFlags"`
		} `json:"cloudMSRtcServiceAttributes"`
		CloudRtcUserPolicies     interface{}   `json:"cloudRtcUserPolicies"`
		CloudSecurityIdentifier  string        `json:"cloudSecurityIdentifier"`
		CloudSipLine             interface{}   `json:"cloudSipLine"`
		CloudSipProxyAddress     string        `json:"cloudSipProxyAddress"`
		CompanyName              string        `json:"companyName"`
		ConsentProvidedForMinor  interface{}   `json:"consentProvidedForMinor"`
		Country                  interface{}   `json:"country"`
		CreatedDateTime          string        `json:"createdDateTime"`
		CreationType             string        `json:"creationType"`
		DeletionTimestamp        interface{}   `json:"deletionTimestamp"`
		Department               string        `json:"department"`
		DirSyncEnabled           bool          `json:"dirSyncEnabled"`
		DisplayName              string        `json:"displayName"`
		EmployeeHireDate         interface{}   `json:"employeeHireDate"`
		EmployeeID               interface{}   `json:"employeeId"`
		EmployeeOrgData          interface{}   `json:"employeeOrgData"`
		EmployeeType             interface{}   `json:"employeeType"`
		ExtensionAttribute1      interface{}   `json:"extensionAttribute1"`
		ExtensionAttribute10     interface{}   `json:"extensionAttribute10"`
		ExtensionAttribute11     interface{}   `json:"extensionAttribute11"`
		ExtensionAttribute12     interface{}   `json:"extensionAttribute12"`
		ExtensionAttribute13     interface{}   `json:"extensionAttribute13"`
		ExtensionAttribute14     interface{}   `json:"extensionAttribute14"`
		ExtensionAttribute15     interface{}   `json:"extensionAttribute15"`
		ExtensionAttribute2      interface{}   `json:"extensionAttribute2"`
		ExtensionAttribute3      interface{}   `json:"extensionAttribute3"`
		ExtensionAttribute4      interface{}   `json:"extensionAttribute4"`
		ExtensionAttribute5      interface{}   `json:"extensionAttribute5"`
		ExtensionAttribute6      interface{}   `json:"extensionAttribute6"`
		ExtensionAttribute7      interface{}   `json:"extensionAttribute7"`
		ExtensionAttribute8      interface{}   `json:"extensionAttribute8"`
		ExtensionAttribute9      interface{}   `json:"extensionAttribute9"`
		FacsimileTelephoneNumber string        `json:"facsimileTelephoneNumber"`
		GivenName                string        `json:"givenName"`
		HasOnPremisesShadow      interface{}   `json:"hasOnPremisesShadow"`
		ImmutableID              string        `json:"immutableId"`
		InfoCatalogs             []interface{} `json:"infoCatalogs"`
		InviteReplyURL           []interface{} `json:"inviteReplyUrl"`
		InviteResources          []interface{} `json:"inviteResources"`
		InviteTicket             []struct {
			Ticket string `json:"ticket"`
			Type   string `json:"type"`
		} `json:"inviteTicket"`
		InvitedAsMail                     string      `json:"invitedAsMail"`
		InvitedOn                         string      `json:"invitedOn"`
		IsCompromised                     interface{} `json:"isCompromised"`
		IsResourceAccount                 interface{} `json:"isResourceAccount"`
		JobTitle                          string      `json:"jobTitle"`
		JrnlProxyAddress                  interface{} `json:"jrnlProxyAddress"`
		LastDirSyncTime                   string      `json:"lastDirSyncTime"`
		LastPasswordChangeDateTime        string      `json:"lastPasswordChangeDateTime"`
		LegalAgeGroupClassification       interface{} `json:"legalAgeGroupClassification"`
		Mail                              string      `json:"mail"`
		MailNickname                      string      `json:"mailNickname"`
		Mobile                            interface{} `json:"mobile"`
		MsExchMailboxGUID                 string      `json:"msExchMailboxGuid"`
		MsExchRecipientTypeDetails        string      `json:"msExchRecipientTypeDetails"`
		MsExchRemoteRecipientType         string      `json:"msExchRemoteRecipientType"`
		NetID                             string      `json:"netId"`
		ObjectID                          string      `json:"objectId"`
		ObjectType                        string      `json:"objectType"`
		Odata_type                        string      `json:"odata.type"`
		OnPremisesDistinguishedName       string      `json:"onPremisesDistinguishedName"`
		OnPremisesPasswordChangeTimestamp string      `json:"onPremisesPasswordChangeTimestamp"`
		OnPremisesSecurityIdentifier      string      `json:"onPremisesSecurityIdentifier"`
		OnPremisesUserPrincipalName       string      `json:"onPremisesUserPrincipalName"`
		OtherMails                        []string    `json:"otherMails"`
		PasswordPolicies                  string      `json:"passwordPolicies"`
		PasswordProfile                   struct {
			EnforceChangePasswordPolicy  bool        `json:"enforceChangePasswordPolicy"`
			ForceChangePasswordNextLogin bool        `json:"forceChangePasswordNextLogin"`
			Password                     interface{} `json:"password"`
		} `json:"passwordProfile"`
		PhysicalDeliveryOfficeName string      `json:"physicalDeliveryOfficeName"`
		PostalCode                 interface{} `json:"postalCode"`
		PreferredDataLocation      interface{} `json:"preferredDataLocation"`
		PreferredLanguage          string      `json:"preferredLanguage"`
		PrimarySMTPAddress         string      `json:"primarySMTPAddress"`
		ProvisionedPlans           []struct {
			CapabilityStatus   string `json:"capabilityStatus"`
			ProvisioningStatus string `json:"provisioningStatus"`
			Service            string `json:"service"`
		} `json:"provisionedPlans"`
		ProvisioningErrors             []interface{} `json:"provisioningErrors"`
		ProxyAddresses                 []string      `json:"proxyAddresses"`
		RefreshTokensValidFromDateTime string        `json:"refreshTokensValidFromDateTime"`
		ReleaseTrack                   interface{}   `json:"releaseTrack"`
		SearchableDeviceKey            []interface{} `json:"searchableDeviceKey"`
		SelfServePasswordResetData     interface{}   `json:"selfServePasswordResetData"`
		ShadowAlias                    string        `json:"shadowAlias"`
		ShadowDisplayName              string        `json:"shadowDisplayName"`
		ShadowLegacyExchangeDN         string        `json:"shadowLegacyExchangeDN"`
		ShadowMail                     string        `json:"shadowMail"`
		ShadowMobile                   interface{}   `json:"shadowMobile"`
		ShadowOtherMobile              []interface{} `json:"shadowOtherMobile"`
		ShadowProxyAddresses           []string      `json:"shadowProxyAddresses"`
		ShadowTargetAddress            string        `json:"shadowTargetAddress"`
		ShadowUserPrincipalName        string        `json:"shadowUserPrincipalName"`
		ShowInAddressList              bool          `json:"showInAddressList"`
		SignInNames                    []string      `json:"signInNames"`
		SignInNamesInfo                []interface{} `json:"signInNamesInfo"`
		SipProxyAddress                string        `json:"sipProxyAddress"`
		SMTPAddresses                  []string      `json:"smtpAddresses"`
		State                          interface{}   `json:"state"`
		StreetAddress                  interface{}   `json:"streetAddress"`
		StrongAuthenticationDetail     struct {
			EncryptedPinHash        interface{}   `json:"encryptedPinHash"`
			EncryptedPinHashHistory interface{}   `json:"encryptedPinHashHistory"`
			Methods                 []interface{} `json:"methods"`
			OathTokenMetadata       []interface{} `json:"oathTokenMetadata"`
			PhoneAppDetails         []interface{} `json:"phoneAppDetails"`
			ProofupTime             interface{}   `json:"proofupTime"`
			Requirements            []interface{} `json:"requirements"`
			VerificationDetail      interface{}   `json:"verificationDetail"`
		} `json:"strongAuthenticationDetail"`
		Surname                         string        `json:"surname"`
		TelephoneNumber                 string        `json:"telephoneNumber"`
		UsageLocation                   string        `json:"usageLocation"`
		UserPrincipalName               string        `json:"userPrincipalName"`
		UserState                       string        `json:"userState"`
		UserStateChangedOn              string        `json:"userStateChangedOn"`
		UserType                        string        `json:"userType"`
		WindowsInformationProtectionKey []interface{} `json:"windowsInformationProtectionKey"`
	} `json:"value"`
}
