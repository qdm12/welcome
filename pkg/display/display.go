package display

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/qdm12/golibs/format"
)

type Display struct {
	randIntn func(n int) int
}

func New() *Display {
	rand.Seed(time.Now().UnixNano())
	return &Display{
		randIntn: rand.Intn,
	}
}

func (d *Display) Error(args ...interface{}) {
	message := format.ArgsToString(args...)
	header := color.HiRedString("ERROR:")
	message = color.HiWhiteString(message)
	fmt.Printf("%s %s\n", header, message)
}

func (d *Display) Warning(args ...interface{}) {
	message := format.ArgsToString(args...)
	header := color.YellowString("WARNING:")
	message = color.HiWhiteString(message)
	fmt.Printf("%s %s\n", header, message)
}
