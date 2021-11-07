package teams

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
)

var log *logger.Logger

type Options struct {
	Email  string
	Token  string
	Thread int
	utils.BaseOptions
}
