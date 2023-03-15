package google

import (
	"GoMapEnum/src/utils"
)

// Options for teams module
type Options struct {
	utils.BaseOptions
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}
