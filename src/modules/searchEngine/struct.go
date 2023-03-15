package searchEngine

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
)

var log *logger.Logger

// Options for searchengine module
type Options struct {
	Format       string
	ExactMatch   bool
	SearchEngine string
	Company      string
	utils.BaseOptions
}
