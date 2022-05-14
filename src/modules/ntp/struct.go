package ntp

import "GoMapEnum/src/utils"

// Options for ntp module
type Options struct {
	UTC bool
	utils.BaseOptions
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}
