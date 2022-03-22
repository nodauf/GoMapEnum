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
	Level  Level
	Module string
	Target string
	Type   string
}

// New create a Logger object for the specified arguments
func New(kind, module, target string) *Logger {
	log := &Logger{}
	// If not set during the init of logger
	if log.Level == 0 {
		log.Level = InfoLevel
	}
	log.Module = module
	log.Target = target
	log.Type = kind
	return log
}

// SetLevel change the logging level
func (logger *Logger) SetLevel(level Level) {
	logger.Level = level
}

// Debug print with debug output
func (logger *Logger) Debug(format string, a ...interface{}) {
	if logger.Level >= DebugLevel {
		str := fmt.Sprintf(format, a...)
		output := color.HiMagentaString(str)
		logger.print(output)
	}
}

// Verbose print with verbose output
func (logger *Logger) Verbose(format string, a ...interface{}) {
	if logger.Level >= VerboseLevel {
		str := fmt.Sprintf(format, a...)
		output := color.HiYellowString(str)
		logger.print(output)
	}
}

// Info print with info output
func (logger *Logger) Info(format string, a ...interface{}) {
	if logger.Level >= InfoLevel {
		str := fmt.Sprintf(format, a...)
		output := color.HiYellowString(str)
		logger.print(output)
	}
}

// Error print with error output
func (logger *Logger) Error(format string, a ...interface{}) {
	if logger.Level >= ErrorLevel {
		str := fmt.Sprintf(format, a...)
		output := color.HiRedString("/!\\ " + str)
		logger.print(output)
	}
}

// Fatal print with fatal output
func (logger *Logger) Fatal(format string, a ...interface{}) {
	if logger.Level >= FatalLevel {
		str := fmt.Sprintf(format, a...)
		output := color.HiRedString("/!\\ " + str)
		logger.print(output)
	}
	os.Exit(1)
}

// Success print with success output (mostly when password is true or user exists)
func (logger *Logger) Success(format string, a ...interface{}) {
	if logger.Level >= InfoLevel {
		str := fmt.Sprintf(format, a...)
		output := color.HiGreenString("[+] " + str)
		logger.print(output)
	}
}

// Fail print with fail output (mostly when password is wrong or user does not exist)
func (logger *Logger) Fail(format string, a ...interface{}) {
	if logger.Level >= VerboseLevel {
		str := fmt.Sprintf(format, a...)
		output := color.HiRedString("[-] " + str)
		logger.print(output)
	}
}
