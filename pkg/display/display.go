package display

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/qdm12/golibs/format"
)

type Display interface {
	Error(arg ...interface{})
	Warning(arg ...interface{})
	FormatRandomASCIIArt(s string) string
	FormatASCIIArt(s, font string) string
}

type display struct {
	randIntn func(n int) int
}

func New() Display {
	rand.Seed(time.Now().UnixNano())
	return &display{
		randIntn: rand.Intn,
	}
}

func (d *display) Error(args ...interface{}) {
	message := format.ArgsToString(args...)
	header := color.HiRedString("ERROR:")
	message = color.HiWhiteString(message)
	fmt.Printf("%s %s\n", header, message)
}

func (d *display) Warning(args ...interface{}) {
	message := format.ArgsToString(args...)
	header := color.YellowString("WARNING:")
	message = color.HiWhiteString(message)
	fmt.Printf("%s %s\n", header, message)
}
