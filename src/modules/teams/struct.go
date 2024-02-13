package teams

import (
	"GoMapEnum/src/utils"
)

// Options for teams module
type Options struct {
	Token           string
	TenantIsPrivate utils.LockedMapBool
	utils.BaseOptions
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}
