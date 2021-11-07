package owa

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/utils"
	"net/http"
)

var log *logger.Logger

type Options struct {
	tr *http.Transport
	utils.BaseOptions
}
