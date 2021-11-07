package utils

import (
	"GoMapEnum/src/logger"
	"errors"
	"net/http"
	"net/url"
)

// BaseOptions is the common options for the module
type BaseOptions struct {
	Users            string
	Passwords        string
	Thread           int
	Log              *logger.Logger
	NoBruteforce     bool
	LockoutThreshold int
	Sleep            int
	Target           string
	CheckIfValid     bool
	Company          string
	Proxy            func(*http.Request) (*url.URL, error)
}

var ErrLockout = errors.New("account is locked")
