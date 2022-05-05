package ldap

import (
	"GoMapEnum/src/utils"
	"crypto/tls"
	"fmt"
	"net"
	"reflect"
	"strings"

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

	kerberoastableAccounts := make(map[string]string)
	kerberoastableAccounts["filter"] = "(&(servicePrincipalName=*)(UserAccountControl:1.2.840.113556.1.4.803:=512)(!(UserAccountControl:1.2.840.113556.1.4.803:=2))(!(objectCategory=computer)))"
	kerberoastableAccounts["attributs"] = "cn,sAMAccountName,servicePrincipalName,PasswordLastSet,LastLogon"
	options.queries["kerberoastableaccounts"] = kerberoastableAccounts
}

func (options *Options) authenticateSimple(username, password string) error {
	ldapConn, err := establisheConnection(options.Target, options.TLS, options.Timeout, options.ProxyTCP)
	if err != nil {
		return fmt.Errorf("fail to establish a connection to the target %s: %w", options.Target, err)
	}
	options.Log.Debug("Connection established to %s", options.Target)
	if password == "" {
		//err = ldapConn.UnauthenticatedBind(username)
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
		if password == "" {
			err = ldapConn.NTLMUnauthenticatedBind(options.Domain, username)
		} else {
			err = ldapConn.NTLMBind(options.Domain, username, password)
		}
	}
	options.ldapConn = ldapConn
	return err
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

func (options *Options) authenticate(username, password string) (bool, error) {
	options.Log.Verbose("Using domain " + options.Domain + " for authentication. Hostname: " + options.Hostname)
	var err error
	if options.UseNTLM {
		err = options.authenticateNTLM(username, password, options.IsHash)
	} else {
		err = options.authenticateSimple(username, password)
	}
	if err != nil {
		if !strings.Contains(err.Error(), "Invalid Credentials") {
			options.Log.Error("fail to authenticate: %s", err.Error())
			return false, err
		}
		return false, nil
	}
	return true, err
}

// FindLDAPServers attempts to find LDAP servers in a domain via DNS. First it attempts looking up LDAP via SRV records,
// if that fails, it will just resolve the domain to an IP and return that.
// credits: https://github.com/ropnop/go-windapsearch/blob/ed05587fd70bbe787d5f69a545a18a0371ce30e8/pkg/dns/dns.go
func findLDAPServers(domain string) (servers []string, err error) {
	_, srvs, err := net.LookupSRV("ldap", "tcp", domain)
	if err != nil {
		if strings.Contains(err.Error(), "No records found") {
			return net.LookupHost(domain)
		}
	}

	for _, s := range srvs {
		servers = append(servers, s.Target)
	}
	// also resolve the domain itself and return that IP
	domainIPs, _ := net.LookupHost(domain)
	servers = append(servers, domainIPs...)

	if len(servers) == 0 {
		err = fmt.Errorf("no LDAP servers found for domain: %s", domain)
		return
	}
	return servers, nil
}

func ParseLDAPData(allData interface{}, columns []string) [][]string {
	var data [][]string
	v := reflect.ValueOf(allData)
	// for each item in slice ( = for each row of the table)
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		var row []string
		dataAttributes := item.Elem().FieldByName("Attributes")
		for j := 0; j < dataAttributes.Len(); j++ {
			if utils.StringInSlice(columns, dataAttributes.Index(j).Elem().FieldByName("Name").String()) {
				values := fmt.Sprint(dataAttributes.Index(j).Elem().FieldByName("Values"))
				values = strings.ReplaceAll(values, "[", "")
				values = strings.ReplaceAll(values, "]", "")
				row = append(row, values)
			}

		}
		data = append(data, row)

	}

	return data
}
