package smb

import (
	"GoMapEnum/src/logger"
	"strings"
	"testing"
)

func TestBruteforceWithPassword(t *testing.T) {
	var results = make(map[string]string)
	results["gomapenumUser1/i3siLdA1se!"] = ""
	results["gomapenumUser2/"] = ""
	results["gomapenumUser3/i3siLdA1se!"] = ""
	results["gomapenumUser4/i3siLdA1se!"] = "response error: The referenced account is currently disabled and may not be logged on to."
	results["gomapenumUser1/wrongPassword"] = "response error: The attempted logon is invalid. This is either due to a bad username or authentication information."
	results["wrongUser/wrongPassword"] = "response error: The attempted logon is invalid. This is either due to a bad username or authentication information."

	options := Options{}
	options.Target = "192.168.1.60"
	log := logger.New("Bruteforce", "SMB", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log

	for credential, wantedResults := range results {
		username := strings.Split(credential, "/")[0]
		password := strings.Split(credential, "/")[1]
		_, err := options.authenticate(username, password)
		if (err == nil && wantedResults != "") || (err != nil && err.Error() != wantedResults) {
			t.Errorf("Authentication for %s returned %v and was expected %v", credential, err, wantedResults)
		}
	}
}

func TestBruteforceWithHash(t *testing.T) {
	var results = make(map[string]string)
	results["gomapenumUser1/bea05617f843e9971f7233c975b2ffc1"] = ""
	results["gomapenumUser2/31d6cfe0d16ae931b73c59d7e0c089c0"] = ""
	results["gomapenumUser3/bea05617f843e9971f7233c975b2ffc1"] = ""
	results["gomapenumUser4/bea05617f843e9971f7233c975b2ffc1"] = "response error: The referenced account is currently disabled and may not be logged on to."
	results["gomapenumUser1/20cc650a5ac276a1cfc22fbc23beada1"] = "response error: The attempted logon is invalid. This is either due to a bad username or authentication information."
	results["wrongUser/bea05617f843e9971f7233c975b2ffc1"] = "response error: The attempted logon is invalid. This is either due to a bad username or authentication information."
	results["wrongUser/wronghash"] = "Cannot decode the hash wronghash from hex to byte: encoding/hex: invalid byte: U+0077 'w'"
	results["wrongUser/bea05617f843e9971f7233c975b2ffc"] = "Cannot decode the hash bea05617f843e9971f7233c975b2ffc from hex to byte: encoding/hex: odd length hex string"

	options := Options{}
	options.Target = "192.168.1.60"
	log := logger.New("Bruteforce", "SMB", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log
	options.IsHash = true

	for credential, wantedResults := range results {
		username := strings.Split(credential, "/")[0]
		password := strings.Split(credential, "/")[1]
		_, err := options.authenticate(username, password)
		if (err == nil && wantedResults != "") || (err != nil && err.Error() != wantedResults) {
			t.Errorf("Authentication for %s returned %v and was expected %v", credential, err, wantedResults)
		}
	}
}
