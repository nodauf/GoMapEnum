package kerberos

import (
	"GoMapEnum/src/logger"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nodauf/gokrb5/v8/iana/errorcode"
	"github.com/nodauf/gokrb5/v8/messages"
)

func TestTestUsername(t *testing.T) {
	var results = make(map[string]string)
	results["gomapenumUser1"] = ""
	results["gomapenumUser2"] = ""
	results["gomapenumUser3"] = ""
	results["gomapenumUser4"] = "KRB Error: " + errorcode.Lookup(errorcode.KDC_ERR_CLIENT_REVOKED)
	results["wrongUser"] = "KRB Error: " + errorcode.Lookup(errorcode.KDC_ERR_C_PRINCIPAL_UNKNOWN)

	options := Options{}
	options.Target = "192.168.1.60"
	options.Domain = "pentest.lab"
	options.DomainController = "192.168.1.60"
	log := logger.New("UserEnumeration", "Kerberos", options.Target)
	log.SetLevel(logger.FatalLevel)
	options.Log = log
	optionsInterface := reflect.ValueOf(&options).Interface()

	KerberosSession(&optionsInterface)

	for username, wantedResults := range results {
		ok, err := options.TestUsername(username)
		if (err == nil && wantedResults != "") || (err != nil && err.Error() != wantedResults) {
			t.Errorf("Authentication for %s returned %t with error %v and was expected %v", username, ok, err, wantedResults)
		}
	}
}

func TestASRepToHascat(t *testing.T) {
	var wantedResults = "$krb5asrep$18$gomapenumUser3@PENTEST.LAB:5992dc2c126d17faf428d584115c7850$7ca1357909085c1673e7bff6f6c2c4a857d558d3334877f6580dd0ba4badefd6d3b70c4b63187ad7a526a524791c93f42916e23e39b6eaeb0ce01d0540e3e5053a63da01c0ec9bf17f1e0465b10016ff9d13bf08cdfd2b1fffa478ac4f92f38569261db4e81e2eaf8c08c9204a45cd2b1958cfc0e9050bbf57ac4e639898bd376eb2bf84583726b5bebaac3abf6b5391bc62d5c878f55033d06d2d55b4e2b29490f0a4bc80b7a57f7cd4d44906e8fe7b9759587eeb5286ef35fc98b7fb6e53e182c500d56967c99cb172c16cd5847060720c1df6ce00ee468f5b5f36a206463e673f973ccd4827b3ffa5be2fc1906c57e8aefbb0e85f9fee71b4a621c8a2"
	var asrep messages.ASRep
	messageJSON := `{"PVNO":5,"MsgType":11,"PAData":[{"PADataType":19,"PADataValue":"MCQwIqADAgESoRsbGVBFTlRFU1QuTEFCZ29tYXBlbnVtVXNlcjM="}],"CRealm":"PENTEST.LAB","CName":{"NameType":1,"NameString":["gomapenumUser3"]},"Ticket":{"TktVNO":5,"Realm":"PENTEST.LAB","SName":{"NameType":2,"NameString":["krbtgt","PENTEST.LAB"]},"EncPart":{"EType":18,"KVNO":2,"Cipher":"UeviMn92vVacaIeICIa5KsU/PThpXyZVgobCAxRhjMBwTPgJR5O5Q5ZEDXENlNuNV+YxL51Do+Yq853PD6Et0k+MAHrQEbAX9U0CVw0v4EbHyKuImaU7Nl7GCsetRuDD6f2zJGe0BjQXLp/AupzmBStYIUHS2tZpdwXWCvP9qPKLV9Co6XFSZ/i4DzZVvyGT7/SsHy472saBQwaEVz/e0ijcJ7HcJnRjwcuh7vL1AiDtY5bkSr9Capii5eZo6gzp+rOU/Z7nFox0QZp23XNZczcx8FcKjgOoFZQDzF0="},"DecryptedEncPart":{"Flags":{"Bytes":null,"BitLength":0},"Key":{"KeyType":0},"CRealm":"","CName":{"NameType":0,"NameString":null},"Transited":{"TRType":0,"Contents":null},"AuthTime":"0001-01-01T00:00:00Z","StartTime":"0001-01-01T00:00:00Z","EndTime":"0001-01-01T00:00:00Z","RenewTill":"0001-01-01T00:00:00Z","CAddr":null,"AuthorizationData":null}},"EncPart":{"EType":18,"KVNO":2,"Cipher":"WZLcLBJtF/r0KNWEEVx4UHyhNXkJCFwWc+e/9vbCxKhX1VjTM0h39lgN0LpLre/W07cMS2MYetelJqUkeRyT9CkW4j45turrDOAdBUDj5QU6Y9oBwOyb8X8eBGWxABb/nRO/CM39Kx//pHisT5LzhWkmHbToHi6vjAjJIEpFzSsZWM/A6QULv1esTmOYmL03brK/hFg3JrW+uqw6v2tTkbxi1ch49VAz0G0tVbTispSQ8KS8gLelf3zU1EkG6P57l1lYfutShu81/Ji3+25T4YLFANVpZ8mcsXLBbNWEcGByDB32zgDuRo9bXzaiBkY+Zz+XPM1IJ7P/pb4vwZBsV+iu+7DoX5/ucbSmIcii"},"DecryptedEncPart":{"Key":{"KeyType":0},"LastReqs":null,"Nonce":0,"KeyExpiration":"0001-01-01T00:00:00Z","Flags":{"Bytes":null,"BitLength":0},"AuthTime":"0001-01-01T00:00:00Z","StartTime":"0001-01-01T00:00:00Z","EndTime":"0001-01-01T00:00:00Z","RenewTill":"0001-01-01T00:00:00Z","SRealm":"","SName":{"NameType":0,"NameString":null},"CAddr":null,"EncPAData":null}}`
	json.Unmarshal([]byte(messageJSON), &asrep)
	hash, err := asRepToHashcat(asrep)
	if err != nil {
		t.Errorf("Error while decoding ASRep message. Got error %v expected nil", err)
	}
	if hash != wantedResults {
		t.Errorf("ASRep message decoded to %s and was expected %s", hash, wantedResults)
	}
}
