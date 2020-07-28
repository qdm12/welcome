package display

import (
	"fmt"
	"github.com/fatih/color"
)

func Error(template string, arg ...interface{}) {
	message := fmt.Sprintf(template, arg...)
	header := color.HiRedString("ERROR:")
	message = color.HiWhiteString(message)
	fmt.Printf("%s %s\n", header, message)
}

func Warning(template string, arg ...interface{}) {
	message := fmt.Sprintf(template, arg...)
	header := color.YellowString("WARNING:")
	message = color.HiWhiteString(message)
	fmt.Printf("%s %s\n", header, message)
}
