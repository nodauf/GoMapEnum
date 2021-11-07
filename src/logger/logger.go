package logger

import (
	"os"

	"github.com/fatih/color"
)

type Level int

const (
	FatalLevel = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	VerboseLevel
	DebugLevel
)

type Logger struct {
	Level  Level
	Module string
	Target string
	Type   string
}

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
func (logger *Logger) Debug(str string) {
	if logger.Level >= DebugLevel {
		output := color.HiMagentaString(str)
		logger.print(output)
	}
}

// Verbose print with verbose output
func (logger *Logger) Verbose(str string) {
	if logger.Level >= VerboseLevel {
		output := color.HiYellowString(str)
		logger.print(output)
	}
}

// Info print with info output
func (logger *Logger) Info(str string) {
	if logger.Level >= InfoLevel {
		output := color.HiYellowString(str)
		logger.print(output)
	}
}

// Error print with error output
func (logger *Logger) Error(str string) {
	if logger.Level >= ErrorLevel {
		output := color.HiRedString("/!\\ " + str)
		logger.print(output)
	}
}

// Fatal print with fatal output
func (logger *Logger) Fatal(str string) {
	if logger.Level >= FatalLevel {
		output := color.HiRedString("/!\\ " + str)
		logger.print(output)
	}
	os.Exit(1)
}

// Success print with success output (mostly when password is true or user exists)
func (logger *Logger) Success(str string) {
	if logger.Level >= InfoLevel {
		output := color.HiGreenString("[+] " + str)
		logger.print(output)
	}
}

// Fail print with fail output (mostly when password is wrong or user does not exist)
func (logger *Logger) Fail(str string) {

	if logger.Level >= VerboseLevel {
		output := color.HiRedString("[-] " + str)
		logger.print(output)
	}
}
