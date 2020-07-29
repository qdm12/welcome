package display

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func (d *display) FormatASCIIArt(s, font string) string {
	if !fontExists(font) {
		d.Warning("Font " + font + " does not exist, picking a random one")
		return d.FormatRandomASCIIArt(s)
	}
	asciiArt := figure.NewFigure(s, font, false).String()
	asciiArt = color.HiCyanString(asciiArt)
	return asciiArt
}

func (d *display) FormatRandomASCIIArt(s string) string {
	fonts := allPossibleFonts()
	font := fonts[d.randIntn(len(fonts))]
	asciiArt := d.FormatASCIIArt(s, font)
	asciiArt += color.HiBlueString("\n~ using font " + font + " ~")
	return asciiArt
}

func (d *display) FormatRandomBestASCIIArt(s string) string {
	fonts := bestFonts()
	font := fonts[d.randIntn(len(fonts))]
	return d.FormatASCIIArt(s, font)
}

func fontExists(font string) (exists bool) {
	for _, possibleFont := range allPossibleFonts() {
		if font == possibleFont {
			return true
		}
	}
	return false
}
