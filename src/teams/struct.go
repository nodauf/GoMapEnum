package teams

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
)

var log *logger.Logger

// Options for teams module
type Options struct {
	Email  string
	Token  string
	Thread int
	utils.BaseOptions
}
