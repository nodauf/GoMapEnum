package orchestrator

import "GoMapEnum/src/utils"

type Orchestrator struct {
	// User Enumeration
	PreActionUserEnum   func(options *interface{}) bool
	CheckBeforeEnumFunc func(options *interface{}, username string) bool
	UserEnumFunc        func(options *interface{}, username string) bool
	PostActionUserEnum  func(options *interface{}) bool

	// Password bruteforce / spraying
	PreActionBruteforce          func(options *interface{}) bool
	CustomOptionsForCheckIfValid func(options *interface{}) interface{}
	AuthenticationFunc           func(options *interface{}, username, password string) bool
	PostActionBruteforce         func(options *interface{}) bool
}

type Options interface {
	GetBaseOptions() *utils.BaseOptions
}
