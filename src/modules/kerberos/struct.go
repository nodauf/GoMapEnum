package kerberos

import (
	"GoMapEnum/src/utils"

	kconfig "github.com/nodauf/gokrb5/v8/config"
)

// Options for kerberos module
type Options struct {
	Timeout          int
	Hash             string
	Domain           string
	DomainController string

	kdcs           map[int]string
	kerberosConfig *kconfig.Config
	utils.BaseOptions
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}
