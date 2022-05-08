package main

import (
	"GoMapEnum/src/utils"
	"crypto/tls"
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/unicode"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/net/proxy"
)

const PASSWORD = "\"i3siLdA1se!\""
const (
	UAC_SCRIPT                         = 0x0001
	UAC_ACCOUNTDISABLE                 = 0x0002
	UAC_HOMEDIR_REQUIRED               = 0x0008
	UAC_LOCKOUT                        = 0x0010
	UAC_PASSWD_NOTREQD                 = 0x0020
	UAC_PASSWD_CANT_CHANGE             = 0x0040
	UAC_ENCRYPTED_TEXT_PWD_ALLOWED     = 0x0080
	UAC_TEMP_DUPLICATE_ACCOUNT         = 0x0100
	UAC_NORMAL_ACCOUNT                 = 0x0200
	UAC_INTERDOMAIN_TRUST_ACCOUNT      = 0x0800
	UAC_WORKSTATION_TRUST_ACCOUNT      = 0x1000
	UAC_SERVER_TRUST_ACCOUNT           = 0x2000
	UAC_DONT_EXPIRE_PASSWORD           = 0x10000
	UAC_MNS_LOGON_ACCOUNT              = 0x20000
	UAC_SMARTCARD_REQUIRED             = 0x40000
	UAC_TRUSTED_FOR_DELEGATION         = 0x80000
	UAC_NOT_DELEGATED                  = 0x100000
	UAC_USE_DES_KEY_ONLY               = 0x200000
	UAC_DONT_REQ_PREAUTH               = 0x400000
	UAC_PASSWORD_EXPIRED               = 0x800000
	UAC_TRUSTED_TO_AUTH_FOR_DELEGATION = 0x1000000
	UAC_PARTIAL_SECRETS_ACCOUNT        = 0x04000000
)

func main() {
	var wait string
	// Bind a LDAP connection
	ldapConn, err := authenticateNTLM("192.168.1.60", "pentest.lab", "vagrant", "vagrant", false)

	if err != nil {
		log.Fatal(err)
	}
	baseDN, _ := getDefaultNamingContext(ldapConn)

	createUser("gomapenumUser1", baseDN, ldapConn)
	createUserEmptyPassword("gomapenumUser2", baseDN, ldapConn)
	createUserWithoutPreAuth("gomapenumUser3", baseDN, ldapConn)
	createDisabledUser("gomapenumUser4", baseDN, ldapConn)
	createUserWithSPN("gomapenumUser5", baseDN, ldapConn)
	createUserPasswordExpired("gomapenumUser6", baseDN, ldapConn)

	fmt.Println("Enter to delete these entries ...")

	fmt.Scanf("%s", &wait)

	// Delete the user
	fmt.Println("Delete the entries")
	deleteUser("gomapenumUser1", baseDN, ldapConn)
	deleteUser("gomapenumUser2", baseDN, ldapConn)
	deleteUser("gomapenumUser3", baseDN, ldapConn)
	deleteUser("gomapenumUser4", baseDN, ldapConn)
	deleteUser("gomapenumUser5", baseDN, ldapConn)
	deleteUser("gomapenumUser6", baseDN, ldapConn)

}

func deleteUser(username, baseDN string, ldapConn *ldap.Conn) {
	del := ldap.NewDelRequest("CN="+username+",CN=Users,"+baseDN, nil)
	err := ldapConn.Del(del)
	if err != nil {
		fmt.Println(err)
	}
}

func createUserEmptyPassword(username, baseDN string, ldapConn *ldap.Conn) {
	// Create the user
	uac := strconv.Itoa(UAC_NORMAL_ACCOUNT | UAC_PASSWD_NOTREQD)
	add := ldap.NewAddRequest("CN="+username+",CN=Users,"+baseDN, nil)
	add.Attribute("description", []string{"GoMapEnum test"})
	add.Attribute("sAMAccountName", []string{username})
	add.Attribute("userAccountControl", []string{uac})
	add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})

	err := ldapConn.Add(add)
	if err != nil {
		fmt.Println(err)
	}
	// Reset the password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// According to the MS docs in the links above
	// The password needs to be enclosed in quotes
	pwdEncoded, _ := utf16.NewEncoder().String("\"\"")
	passReq := ldap.NewModifyRequest("CN="+username+",CN=Users,"+baseDN, nil)
	passReq.Replace("unicodePwd", []string{pwdEncoded})
	err = ldapConn.Modify(passReq)

	if err != nil {
		fmt.Printf("Password could not be changed: %s\n", err.Error())
	}
}

func createDisabledUser(username, baseDN string, ldapConn *ldap.Conn) {
	// Create the user
	uac := strconv.Itoa(UAC_NORMAL_ACCOUNT | UAC_PASSWD_NOTREQD | UAC_ACCOUNTDISABLE)
	add := ldap.NewAddRequest("CN="+username+",CN=Users,"+baseDN, nil)
	add.Attribute("description", []string{"GoMapEnum test"})
	add.Attribute("sAMAccountName", []string{username})
	add.Attribute("userAccountControl", []string{uac})
	add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})

	err := ldapConn.Add(add)
	if err != nil {
		fmt.Println(err)
	}
	// Reset the password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// According to the MS docs in the links above
	// The password needs to be enclosed in quotes
	pwdEncoded, _ := utf16.NewEncoder().String(PASSWORD)
	passReq := ldap.NewModifyRequest("CN="+username+",CN=Users,"+baseDN, nil)
	passReq.Replace("unicodePwd", []string{pwdEncoded})
	err = ldapConn.Modify(passReq)

	if err != nil {
		fmt.Printf("Password could not be changed: %s\n", err.Error())
	}
}

func createUserWithoutPreAuth(username, baseDN string, ldapConn *ldap.Conn) {
	// Create the user
	uac := strconv.Itoa(UAC_DONT_REQ_PREAUTH | UAC_NORMAL_ACCOUNT | UAC_PASSWD_NOTREQD)
	add := ldap.NewAddRequest("CN="+username+",CN=Users,"+baseDN, nil)
	add.Attribute("description", []string{"GoMapEnum test"})
	add.Attribute("sAMAccountName", []string{username})
	add.Attribute("userAccountControl", []string{uac})
	add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})

	err := ldapConn.Add(add)
	if err != nil {
		fmt.Println(err)
	}

	// Reset the password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// According to the MS docs in the links above
	// The password needs to be enclosed in quotes
	pwdEncoded, _ := utf16.NewEncoder().String(PASSWORD)
	passReq := ldap.NewModifyRequest("CN="+username+",CN=Users,"+baseDN, nil)
	passReq.Replace("unicodePwd", []string{pwdEncoded})
	err = ldapConn.Modify(passReq)

	if err != nil {
		fmt.Printf("Password could not be changed: %s\n", err.Error())
	}
}

func createUserPasswordExpired(username, baseDN string, ldapConn *ldap.Conn) {
	// Create the user
	uac := strconv.Itoa(UAC_NORMAL_ACCOUNT | UAC_PASSWD_NOTREQD)
	add := ldap.NewAddRequest("CN="+username+",CN=Users,"+baseDN, nil)
	add.Attribute("description", []string{"GoMapEnum test"})
	add.Attribute("sAMAccountName", []string{username})
	add.Attribute("userAccountControl", []string{uac})
	add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})
	add.Attribute("pwdLastSet", []string{"0"})

	err := ldapConn.Add(add)
	if err != nil {
		fmt.Println(err)
	}
	// Reset the password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// According to the MS docs in the links above
	// The password needs to be enclosed in quotes
	pwdEncoded, _ := utf16.NewEncoder().String(PASSWORD)
	passReq := ldap.NewModifyRequest("CN="+username+",CN=Users,"+baseDN, nil)
	passReq.Replace("unicodePwd", []string{pwdEncoded})
	err = ldapConn.Modify(passReq)

	if err != nil {
		fmt.Printf("Password could not be changed: %s\n", err.Error())
	}

	passwordExpire := ldap.NewModifyRequest("CN="+username+",CN=Users,"+baseDN, nil)
	passwordExpire.Replace("pwdLastSet", []string{"0"})
	err = ldapConn.Modify(passwordExpire)

	if err != nil {
		fmt.Printf("Could not expire the user's password: %s\n", err.Error())
	}
}

func createUser(username, baseDN string, ldapConn *ldap.Conn) {
	// Create the user
	uac := strconv.Itoa(UAC_NORMAL_ACCOUNT | UAC_PASSWD_NOTREQD)
	add := ldap.NewAddRequest("CN="+username+",CN=Users,"+baseDN, nil)
	add.Attribute("description", []string{"GoMapEnum test"})
	add.Attribute("sAMAccountName", []string{username})
	add.Attribute("userAccountControl", []string{uac})
	add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})

	err := ldapConn.Add(add)
	if err != nil {
		fmt.Println(err)
	}

	// Reset the password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// According to the MS docs in the links above
	// The password needs to be enclosed in quotes
	pwdEncoded, _ := utf16.NewEncoder().String(PASSWORD)
	passReq := ldap.NewModifyRequest("CN="+username+",CN=Users,"+baseDN, nil)
	passReq.Replace("unicodePwd", []string{pwdEncoded})
	err = ldapConn.Modify(passReq)

	if err != nil {
		fmt.Printf("Password could not be changed: %s\n", err.Error())
	}
}

func createUserWithSPN(username, baseDN string, ldapConn *ldap.Conn) {
	// Create the user
	uac := strconv.Itoa(UAC_NORMAL_ACCOUNT | UAC_PASSWD_NOTREQD)
	add := ldap.NewAddRequest("CN="+username+",CN=Users,"+baseDN, nil)
	add.Attribute("description", []string{"GoMapEnum test"})
	add.Attribute("sAMAccountName", []string{username})
	add.Attribute("userAccountControl", []string{uac})
	add.Attribute("servicePrincipalName", []string{"whatever/random.xx.lan"})
	add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})

	err := ldapConn.Add(add)
	if err != nil {
		fmt.Println(err)
	}

	// Reset the password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// According to the MS docs in the links above
	// The password needs to be enclosed in quotes
	pwdEncoded, _ := utf16.NewEncoder().String(PASSWORD)
	passReq := ldap.NewModifyRequest("CN="+username+",CN=Users,"+baseDN, nil)
	passReq.Replace("unicodePwd", []string{pwdEncoded})
	err = ldapConn.Modify(passReq)

	if err != nil {
		fmt.Printf("Password could not be changed: %s\n", err.Error())
	}
}

func authenticateNTLM(target, domain, username, password string, isHash bool) (*ldap.Conn, error) {

	ldapConn, err := establisheConnection(target, "starttls", 5, nil)
	if err != nil || ldapConn == nil {
		return nil, fmt.Errorf("fail to establish a connection to the target %s: %w", target, err)
	}
	fmt.Printf("Connection established to %s\n", target)
	if isHash {
		err = ldapConn.NTLMBindWithHash(domain, username, password)
	} else {
		err = ldapConn.NTLMBind(domain, username, password)
	}
	return ldapConn, err
}
func establisheConnection(target string, TLSMode string, timeout int, proxyTCP proxy.Dialer) (*ldap.Conn, error) {

	var port string
	switch strings.ToLower(TLSMode) {
	case "tls":
		port = ldap.DefaultLdapsPort
	case "starttls", "notls":
		port = ldap.DefaultLdapPort
	default:
		return nil, fmt.Errorf("invalid TLSMode %s", TLSMode)
	}
	conn, err := utils.OpenConnectionWoProxy(target, port, timeout, proxyTCP)

	fmt.Printf("connect to %s:%s\n", target, port)
	// Check if connection is successful
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the target " + target + ":" + ldap.DefaultLdapPort + ": " + err.Error())
	}

	var ldapConnection *ldap.Conn
	switch strings.ToLower(TLSMode) {
	case "tls":
		tlsConn := tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
		ldapConnection = ldap.NewConn(tlsConn, true)
		ldapConnection.Start()
	case "starttls":
		ldapConnection = ldap.NewConn(conn, false)
		ldapConnection.Start()
		err = ldapConnection.StartTLS(&tls.Config{InsecureSkipVerify: true})

	case "notls":
		ldapConnection = ldap.NewConn(conn, false)
		ldapConnection.Start()
	default:
		return nil, fmt.Errorf("invalid TLSMode %s", TLSMode)
	}
	return ldapConnection, err

}
func getDefaultNamingContext(ldapConn *ldap.Conn) (string, error) {

	sr := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"defaultNamingContext"},
		nil)
	fmt.Println("LDAP request (objectClass=*) with attribut defaultNamingContext")
	res, err := ldapConn.Search(sr)
	if err != nil {
		return "", err
	}
	if len(res.Entries) == 0 {
		return "", fmt.Errorf("error getting metadata: No LDAP responses from server")
	}
	defaultNamingContext := res.Entries[0].GetAttributeValue("defaultNamingContext")
	if defaultNamingContext == "" {
		return "", fmt.Errorf("error getting metadata: attribute defaultNamingContext missing")
	}
	return defaultNamingContext, nil

}
