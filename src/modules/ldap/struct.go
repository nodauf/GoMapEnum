package ldap

import (
	"GoMapEnum/src/utils"

	"github.com/go-ldap/ldap/v3"
)

// Options for ldap module
type Options struct {
	Timeout  int
	Hash     string
	Domain   string
	Hostname string
	BaseDN   string
	TLS      string
	utils.BaseOptions
	DumpObjects []string
	UseNTLM     bool
	HTML        bool
	JSON        bool
	IsHash      bool

	ldapConn *ldap.Conn
	queries  map[string]map[string]string
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}
