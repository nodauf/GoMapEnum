package smtp

import (
	"GoMapEnum/src/utils"

	smtp "github.com/nodauf/net-smtp"
)

// Options for o365 module
type Options struct {
	Target string
	Domain string
	Mode   string
	utils.BaseOptions

	all             bool
	connectionsPool chan *smtp.Client
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}
