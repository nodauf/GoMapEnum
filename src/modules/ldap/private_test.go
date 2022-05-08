package ldap

import (
	"GoMapEnum/src/logger"
	"strings"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

func TestEstablisheConnectionLDAPNoTLS(t *testing.T) {
	ldapConnNoTLS, err := establisheConnection("192.168.1.60", "NoTLS", 5, nil)
	if err != nil {
		t.Errorf("establisheConnection for NoTLS mode return the error: %s", err.Error())
	}
	if ldapConnNoTLS == nil {
		t.Error("ldapConnNoTLS is nil")
	}
	ldapConnNoTLS.Close()

	ldapConnStartTLS, err := establisheConnection("192.168.1.60", "StartTLS", 5, nil)
	if err != nil {
		t.Errorf("establisheConnection for StartTLS mode return the error: %s", err.Error())
	}
	if ldapConnStartTLS == nil {
		t.Error("ldapConnStartTLS is nil")
	}

	ldapConnUnknowTLS, err := establisheConnection("192.168.1.60", "NotSupportedTLSMode", 5, nil)
	expectedError := "invalid TLSMode NotSupportedTLSMode"
	if err == nil || err != nil && err.Error() != expectedError {
		t.Errorf("establisheConnection for NotSupportedTLSMode return the error %v and was expected %s", err, expectedError)
	}
	if ldapConnUnknowTLS != nil {
		t.Error("ldapConnUnknowTLS is not nil")
	}
	ldapConnStartTLS.Close()
}

func TestEstablisheConnectionLDAPS(t *testing.T) {
	ldapConn, err := establisheConnection("192.168.1.60", "TLS", 5, nil)
	if err != nil {
		t.Errorf("establisheConnection return the error: %s", err.Error())
	}
	if ldapConn == nil {
		t.Error("ldapConn is nil")
	}
	ldapConn.Close()
}

func TestAuthenticateNTLMWithPassword(t *testing.T) {
	var results = make(map[string]int)
	results["gomapenumUser1/i3siLdA1se!"] = 0
	results["gomapenumUser2/"] = 0
	results["gomapenumUser3/i3siLdA1se!"] = 0
	results["gomapenumUser4/i3siLdA1se!"] = ldap.LDAPResultInvalidCredentials // The error is Invalid credentials if the account is disabled
	results["gomapenumUser1/wrongPassword"] = ldap.LDAPResultInvalidCredentials
	results["wrongUser/wrongPassword"] = ldap.LDAPResultInvalidCredentials

	options := Options{}
	options.Target = "192.168.1.60"
	options.TLS = "NoTLS"
	log := logger.New("Bruteforce", "SMB", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log

	for credential, wantedResults := range results {
		username := strings.Split(credential, "/")[0]
		password := strings.Split(credential, "/")[1]
		err := options.authenticateNTLM(username, password, false)
		options.ldapConn.Close()
		// Test for success
		if err == nil && wantedResults == 0 {
			continue
		}
		if !ldap.IsErrorWithCode(err, uint16(wantedResults)) {
			t.Errorf("Authentication for %s returned %v and was expected %v", credential, err, ldap.LDAPResultCodeMap[uint16(wantedResults)])
		}
	}

}

func TestAuthenticateNTLMWithHash(t *testing.T) {
	var results = make(map[string]int)
	results["gomapenumUser1/bea05617f843e9971f7233c975b2ffc1"] = 0
	results["gomapenumUser2/31d6cfe0d16ae931b73c59d7e0c089c0"] = 0
	results["gomapenumUser3/bea05617f843e9971f7233c975b2ffc1"] = 0
	results["gomapenumUser4/bea05617f843e9971f7233c975b2ffc1"] = ldap.LDAPResultInvalidCredentials
	results["gomapenumUser1/bea05617f843e9971f7233c975b2ff"] = ldap.LDAPResultInvalidCredentials
	results["wrongUser/bea05617f843e9971f7233c975b2ffc1"] = ldap.LDAPResultInvalidCredentials

	options := Options{}
	options.Target = "192.168.1.60"
	options.TLS = "NoTLS"
	log := logger.New("Bruteforce", "SMB", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log

	for credential, wantedResults := range results {
		username := strings.Split(credential, "/")[0]
		password := strings.Split(credential, "/")[1]
		err := options.authenticateNTLM(username, password, true)
		options.ldapConn.Close()
		// Test for success
		if err == nil && wantedResults == 0 {
			continue
		}
		if !ldap.IsErrorWithCode(err, uint16(wantedResults)) {
			t.Errorf("Authentication for %s returned %v and was expected %v", credential, err, ldap.LDAPResultCodeMap[uint16(wantedResults)])
		}
	}

}

func TestAuthenticate(t *testing.T) {
	type testUser struct {
		username string
		password string
		err      string
		valid    bool
	}

	var results []testUser
	results = append(results, testUser{username: "gomapenumUser1", password: "i3siLdA1se!", err: "", valid: true})
	results = append(results, testUser{username: "gomapenumUser2", password: "", err: "", valid: true})
	results = append(results, testUser{username: "gomapenumUser3", password: "i3siLdA1se!", err: "", valid: true})
	results = append(results, testUser{username: "gomapenumUser4", password: "i3siLdA1se!", err: "account disabled", valid: true})
	results = append(results, testUser{username: "gomapenumUser6", password: "i3siLdA1se!", err: "user must reset password", valid: true})
	results = append(results, testUser{username: "gomapenumUser1", password: "wrongPassword", err: "Invalid Credentials", valid: false})
	results = append(results, testUser{username: "wrongUser", password: "wrongPassword", err: "Invalid Credentials", valid: false})

	options := Options{}
	options.Target = "192.168.1.60"
	options.TLS = "NoTLS"
	options.UseNTLM = true
	log := logger.New("Bruteforce", "SMB", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log

	for _, wantedResults := range results {
		valid, err := options.authenticate(wantedResults.username, wantedResults.password)
		options.ldapConn.Close()
		if !(err == nil && wantedResults.err == "") && !strings.Contains(err.Error(), wantedResults.err) {
			t.Errorf("Authentication for %s/%s returned %v and was expected %v", wantedResults.username, wantedResults.password, err, wantedResults.err)
		}

		if valid != wantedResults.valid {
			t.Errorf("Authentication for %s/%s returned %v and was expected %v", wantedResults.username, wantedResults.password, valid, wantedResults.valid)
		}
	}

}

/*func TestAuthenticateSimpleWithPassword(t *testing.T) {
	var results = make(map[string]int)
	results["gomapenumUser1/i3siLdA1se!"] = 0
	results["gomapenumUser2/"] = 0
	results["gomapenumUser1/wrongPassword"] = ldap.LDAPResultInvalidCredentials
	results["wrongUser/wrongPassword"] = ldap.LDAPResultInvalidCredentials

	ldapConn, err := establisheConnection("192.168.1.60", false, 5, nil)
	if err != nil {
		t.Errorf("establisheConnection return the error: %s", err.Error())
	}
	if ldapConn == nil {
		t.Error("ldapConn is nil")
	}
	options := Options{}
	options.Target = "192.168.1.60"
	options.LockoutThreshold = 1
	log := logger.New("Bruteforce", "SMB", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log

	for credential, wantedResults := range results {
		username := strings.Split(credential, "/")[0]
		password := strings.Split(credential, "/")[1]
		err := options.authenticateSimple(username, password)
		// Test for success
		if err == nil && wantedResults == 0 {
			continue
		}
		if !ldap.IsErrorWithCode(err, uint16(wantedResults)) {
			t.Errorf("Authentication for %s returned %v and was expected %v", credential, err, ldap.LDAPResultCodeMap[uint16(wantedResults)])
		}
	}

}*/
