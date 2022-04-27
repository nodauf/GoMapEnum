package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

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

	fmt.Println("Enter to delete these entries ...")

	fmt.Scanf("%s", &wait)

	// Delete the user
	fmt.Println("Delete the entries")
	deleteUser("gomapenumUser1", baseDN, ldapConn)
	deleteUser("gomapenumUser2", baseDN, ldapConn)
	deleteUser("gomapenumUser3", baseDN, ldapConn)
	deleteUser("gomapenumUser4", baseDN, ldapConn)

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

func authenticateNTLM(target, domain, username, password string, isHash bool) (*ldap.Conn, error) {

	ldapConn, err := establisheConnection(target, true, 5, nil)
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
func establisheConnection(target string, TLS bool, timeout int, proxyTCP proxy.Dialer) (*ldap.Conn, error) {
	var conn net.Conn
	var err error
	var port string
	if TLS {
		port = ldap.DefaultLdapsPort
	} else {
		port = ldap.DefaultLdapPort
	}
	if proxyTCP != nil {
		conn, err = proxyTCP.Dial("tcp", fmt.Sprintf("%s:%s", target, port))
	} else {
		defaultDialer := &net.Dialer{Timeout: time.Duration(timeout * int(time.Second))}
		conn, err = defaultDialer.Dial("tcp", fmt.Sprintf("%s:%s", target, port))
	}
	fmt.Printf("connect to %s:%s\n", target, port)
	// Check if connection is successful
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the target " + target + ":" + ldap.DefaultLdapPort + ": " + err.Error())
	}

	var ldapConnection *ldap.Conn
	if TLS {
		tlsConn := tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
		ldapConnection = ldap.NewConn(tlsConn, TLS)
	} else {
		ldapConnection = ldap.NewConn(conn, TLS)
	}
	ldapConnection.Start()

	return ldapConnection, nil
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
