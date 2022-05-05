package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Level is a type representing the logging level
type Level int

const (
	FatalLevel = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	VerboseLevel
	DebugLevel
)

// Logger is the options for the logging module
type Logger struct {
	Level   Level
	NoColor bool
	Mode    string
	Module  string
	Target  string
	File    string
}

// New create a Logger object for the specified arguments
func New(mode, module, target string) *Logger {
	log := &Logger{}
	// If not set during the init of logger
	if log.Level == 0 {
		log.Level = InfoLevel
	}
	log.Mode = mode
	log.Module = module
	log.Target = target
	return log
}

// SetTarget update the target field
func (logger *Logger) SetTarget(target string) {
	logger.Target = target
}

// SetModule update the module field
func (logger *Logger) SetModule(module string) {
	logger.Module = module
}

// SetLevel change the logging level
func (logger *Logger) SetLevel(level Level) {
	logger.Level = level
}

// Debug print with debug output
func (logger *Logger) Debug(format string, a ...interface{}) {

	if logger.Level >= DebugLevel {
		var output string
		str := fmt.Sprintf(format, a...)
		if logger.NoColor {
			output = str
		} else {
			output = color.HiBlackString("[Debug] " + str)
		}

		logger.print(output)
	}
}

// Verbose print with verbose output
func (logger *Logger) Verbose(format string, a ...interface{}) {
	if logger.Level >= VerboseLevel {
		var output string
		str := fmt.Sprintf(format, a...)
		if logger.NoColor {
			output = str
		} else {
			output = color.HiYellowString(str)
		}
		logger.print(output)
	}
}

// Info print with info output
func (logger *Logger) Info(format string, a ...interface{}) {
	if logger.Level >= InfoLevel {
		var output string
		str := fmt.Sprintf(format, a...)
		if logger.NoColor {
			output = str
		} else {
			output = color.HiYellowString(str)
		}
		logger.print(output)
	}
}

// Error print with error output
func (logger *Logger) Error(format string, a ...interface{}) {
	if logger.Level >= ErrorLevel {
		var output string
		str := fmt.Sprintf(format, a...)
		if logger.NoColor {
			output = str
		} else {
			output = color.HiRedString("/!\\ Error: " + str)
		}
		logger.print(output)
	}
}

// Fatal print with fatal output
func (logger *Logger) Fatal(format string, a ...interface{}) {
	if logger.Level >= FatalLevel {
		var output string
		str := fmt.Sprintf(format, a...)
		if logger.NoColor {
			output = str
		} else {
			output = color.HiRedString("/!\\ " + str)
		}
		logger.print(output)
	}
	os.Exit(1)
}

// Success print with success output (mostly when password is true or user exists)
func (logger *Logger) Success(format string, a ...interface{}) {
	if logger.Level >= InfoLevel {
		var output string
		str := fmt.Sprintf(format, a...)
		if logger.NoColor {
			output = str
		} else {
			output = color.HiGreenString("[+] " + str)
		}
		logger.print(output)
	}
}

// Fail print with fail output (mostly when password is wrong or user does not exist)
func (logger *Logger) Fail(format string, a ...interface{}) {

	if logger.Level >= VerboseLevel {
		var output string
		str := fmt.Sprintf(format, a...)
		if logger.NoColor {
			output = str
		} else {
			output = color.HiRedString("[-] " + str)
		}
		logger.print(output)
	}
}
