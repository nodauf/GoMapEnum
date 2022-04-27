package smb

import (
	"testing"
)

func TestRetrieveTargetInfo(t *testing.T) {
	var excpectedDomain = "PENTEST"
	var excpectedHostname = "DC"
	var target = "192.168.1.60"

	domainName, hostname, err := GetTargetInfo(target, 5, nil)
	if err == nil {
		if domainName != excpectedDomain {
			t.Errorf("The detected domain is %s and was expected %s", domainName, excpectedDomain)
		}
		if hostname != excpectedHostname {
			t.Errorf("The detected hostname is %s and was expected %s", hostname, excpectedHostname)
		}
	} else {
		t.Errorf("GetTargetInfo return the error: %s", err.Error())
	}
}
