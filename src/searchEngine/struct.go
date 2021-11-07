package searchEngine

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
)

var log *logger.Logger

type Options struct {
	Format     string
	ExactMatch bool
	utils.BaseOptions
}
