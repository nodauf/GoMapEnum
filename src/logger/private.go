package logger

import (
	"fmt"

	"github.com/fatih/color"
)

func (logger *Logger) print(str string) {
	output := color.HiCyanString(logger.Type+" - "+logger.Module) + "\t\t" + logger.Target + "\t\t" + str + "\n"
	//fmt.Println(output)
	// For Linux and Windows support of colored output
	fmt.Fprint(color.Output, output)

}
