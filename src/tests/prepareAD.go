package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/text/encoding/unicode"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/net/proxy"
)

var PASSWORD = "i3siLdA1se!"

func main() {
	var wait string
	// Bind a LDAP connection
	ldapConn, err := authenticateNTLM("192.168.1.60", "pentest.lab", "vagrant", "vagrant", false)
	if err != nil {
		log.Fatal(err)
	}
	baseDN, _ := getDefaultNamingContext(ldapConn)

	createUser("gomapenumUser1", baseDN, ldapConn)

	fmt.Println("Enter to contiue ...")

	fmt.Scanf("%s", &wait)

	// Delete the user
	fmt.Println("Delete the entry")
	del := ldap.NewDelRequest("CN=gomapenumUser1,CN=Users,"+baseDN, nil)
	err = ldapConn.Del(del)
	if err != nil {
		fmt.Println(err)
	}
}

func createUser(username, baseDN string, ldapConn *ldap.Conn) {
	// Create the user
	add := ldap.NewAddRequest("CN=gomapenumUser1,CN=Users,"+baseDN, nil)
	add.Attribute("description", []string{"GoMapEnum test"})
	add.Attribute("sAMAccountName", []string{"gomapenumUser1"})
	add.Attribute("userAccountControl", []string{"544"})
	//add.Attribute("memberOf", []string{"CN=Users,CN=Builtin," + baseDN})
	add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})

	err := ldapConn.Add(add)
	if err != nil {
		fmt.Println(err)
	}

	// Reset the password
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// According to the MS docs in the links above
	// The password needs to be enclosed in quotes
	pwdEncoded, _ := utf16.NewEncoder().String("\"Oooooo12!\"")
	passReq := ldap.NewModifyRequest("CN=gomapenumUser1,CN=Users,"+baseDN, nil)
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
