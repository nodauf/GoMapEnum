package ldap

import (
	"GoMapEnum/src/logger"
	"strings"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

func TestEstablisheConnectionLDAP(t *testing.T) {
	ldapConn, err := establisheConnection("192.168.1.60", false, 5, nil)
	if err != nil {
		t.Errorf("establisheConnection return the error: %s", err.Error())
	}
	if ldapConn == nil {
		t.Error("ldapConn is nil")
	}
}

func TestEstablisheConnectionLDAPS(t *testing.T) {
	ldapConn, err := establisheConnection("192.168.1.60", true, 5, nil)
	if err != nil {
		t.Errorf("establisheConnection return the error: %s", err.Error())
	}
	if ldapConn == nil {
		t.Error("ldapConn is nil")
	}
}

func TestAuthenticateNTLMWithPassword(t *testing.T) {
	var results = make(map[string]int)
	results["gomapenumUser1/i3siLdA1se!"] = 0
	results["gomapenumUser2/"] = 0
	results["gomapenumUser3/i3siLdA1se!"] = 0
	results["gomapenumUser4/i3siLdA1se!"] = ldap.LDAPResultInvalidCredentials // The error is Invalid credentials if the account is disabled
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
	log := logger.New("Bruteforce", "SMB", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log

	for credential, wantedResults := range results {
		username := strings.Split(credential, "/")[0]
		password := strings.Split(credential, "/")[1]
		err := options.authenticateNTLM(username, password, false)
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

	ldapConn, err := establisheConnection("192.168.1.60", false, 5, nil)
	if err != nil {
		t.Errorf("establisheConnection return the error: %s", err.Error())
	}
	if ldapConn == nil {
		t.Error("ldapConn is nil")
	}
	options := Options{}
	options.Target = "192.168.1.60"
	log := logger.New("Bruteforce", "SMB", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log

	for credential, wantedResults := range results {
		username := strings.Split(credential, "/")[0]
		password := strings.Split(credential, "/")[1]
		err := options.authenticateNTLM(username, password, true)
		// Test for success
		if err == nil && wantedResults == 0 {
			continue
		}
		if !ldap.IsErrorWithCode(err, uint16(wantedResults)) {
			t.Errorf("Authentication for %s returned %v and was expected %v", credential, err, ldap.LDAPResultCodeMap[uint16(wantedResults)])
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
