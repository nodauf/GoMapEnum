package owa

import (
	"GoMapEnum/src/utils"
	"net/http"
	"time"
)

// Options for owa module
type Options struct {
	Basic          bool
	tr             *http.Transport
	internalDomain string
	urlToHarvest   string
	avgResponse    time.Duration
	utils.BaseOptions
}

func (options *Options) GetBaseOptions() *utils.BaseOptions {
	return &options.BaseOptions
}
