package utils

import (
	"GoMapEnum/src/logger"
	"errors"
	"net/http"
	"net/url"
	"sync"
)

// BaseOptions is the common options for the module
type BaseOptions struct {
	Users            string
	UsernameList     []string
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
	Mutex            sync.Mutex
}

// ErrLockout is the error to returned when an account is locked
var ErrLockout = errors.New("account is locked")
