package ldap

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/net/proxy"
)

func (options *Options) initDumpMap() {
	options.queries = make(map[string]map[string]string)

	computers := make(map[string]string)
	computers["filter"] = "(objectClass=Computer)"
	computers["attributs"] = "cn,dNSHostName,operatingSystem,operatingSystemVersion,operatingSystemServicePack,whenCreated,lastLogon,objectSid,objectClass"
	options.queries["computers"] = computers

	users := make(map[string]string)
	users["filter"] = "(objectClass=user)"
	users["attributs"] = "cn,sAMAccountName,userPrincipalName,objectClass"
	options.queries["users"] = users
}

func (options *Options) authenticateSimple(username, password string) error {
	ldapConn, err := establisheConnection(options.Target, options.TLS, options.Timeout, options.ProxyTCP)
	if err != nil {
		return fmt.Errorf("fail to establish a connection to the target %s: %w", options.Target, err)
	}
	options.Log.Debug("Connection established to %s", options.Target)
	if username == "" {
		err = ldapConn.UnauthenticatedBind("")
	} else {
		err = ldapConn.Bind(username, password)
	}
	options.ldapConn = ldapConn
	return err
}

func (options *Options) authenticateNTLM(username, password string, isHash bool) error {

	ldapConn, err := establisheConnection(options.Target, options.TLS, options.Timeout, options.ProxyTCP)
	if err != nil {
		return fmt.Errorf("fail to establish a connection to the target %s: %w", options.Target, err)
	}
	options.Log.Debug("Connection established to %s", options.Target)
	if isHash {
		err = ldapConn.NTLMBindWithHash(options.Domain, username, password)
	} else {
		err = ldapConn.NTLMBind(options.Domain, username, password)
	}
	options.ldapConn = ldapConn
	return err
}

func (options *Options) getDefaultNamingContext() error {
	if options.BaseDN != "" {
		return nil
	}
	sr := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"defaultNamingContext"},
		nil)
	options.Log.Debug("LDAP request (objectClass=*) with attribut defaultNamingContext")
	res, err := options.ldapConn.Search(sr)
	if err != nil {
		return err
	}
	if len(res.Entries) == 0 {
		return fmt.Errorf("error getting metadata: No LDAP responses from server")
	}
	defaultNamingContext := res.Entries[0].GetAttributeValue("defaultNamingContext")
	if defaultNamingContext == "" {
		return fmt.Errorf("error getting metadata: attribute defaultNamingContext missing")
	}
	options.BaseDN = defaultNamingContext
	return nil

}

func (options *Options) dumpObject(object string) []*ldap.Entry {
	if _, ok := options.queries[object]; !ok {
		options.Log.Error("Not able to dump %s. The query is not implemented.", object)
		return nil
	}
	ldapResult := executeLdapQuery(options.ldapConn, options.BaseDN, options.queries[object])

	// Print the results
	for _, entry := range ldapResult.Entries {
		options.Log.Success(object + ": " + entry.DN)
	}
	return ldapResult.Entries
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
		conn, err = defaultDialer.Dial("tcp", fmt.Sprintf("%s:%s", target, ldap.DefaultLdapPort))
	}
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

func executeLdapQuery(ldapConn *ldap.Conn, baseDN string, filterAndAttributs map[string]string) *ldap.SearchResult {
	filter := filterAndAttributs["filter"]
	attributs := strings.Split(filterAndAttributs["attributs"], ",")
	ldapSearchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		attributs,
		nil)

	ldapResult, _ := ldapConn.Search(ldapSearchRequest)
	return ldapResult

}
