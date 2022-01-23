package orchestrator

import "GoMapEnum/src/utils"

type Orchestrator struct {
	// User Enumeration
	PreActionUserEnum   func(options *interface{})
	CheckBeforeEnumFunc func(options *interface{}, username string) bool
	UserEnumFunc        func(options *interface{}, username string) bool
	PostActionUserEnum  func(options *interface{})

	// Password bruteforce / spraying
	PreActionBruteforce          func(options *interface{})
	CustomOptionsForCheckIfValid func(options *interface{}) interface{}
	AuthenticationFunc           func(options *interface{}, username, password string) bool
}

type Options interface {
	GetBaseOptions() *utils.BaseOptions
}
