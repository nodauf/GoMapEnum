package smb

import (
	"GoMapEnum/src/utils"
)

// Options for smb module
type Options struct {
	Timeout  int
	Hash     string
	Domain   string
	Hostname string
	utils.BaseOptions
	IsHash bool

	lockoutCounter int
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}
