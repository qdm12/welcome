package display

import (
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

var possibleFonts = []string{"3-d", "3x5", "5lineoblique", "acrobatic", "alligator", "alligator2", "alphabet", "avatar", "banner", "banner3-D", "banner3", "banner4", "barbwire", "basic", "bell", "big", "bigchief", "binary", "block", "bubble", "bulbhead", "calgphy2", "caligraphy", "catwalk", "chunky", "coinstak", "colossal", "computer", "contessa", "contrast", "cosmic", "cosmike", "cricket", "cursive", "cyberlarge", "cybermedium", "cybersmall", "diamond", "digital", "doh", "doom", "dotmatrix", "drpepper", "eftichess", "eftifont", "eftipiti", "eftirobot", "eftitalic", "eftiwall", "eftiwater", "epic", "fender", "fourtops", "fuzzy", "goofy", "gothic", "graffiti", "hollywood", "invita", "isometric1", "isometric2", "isometric3", "isometric4", "italic", "ivrit", "jazmine", "jerusalem", "katakana", "kban", "larry3d", "lcd", "lean", "letters", "linux", "lockergnome", "madrid", "marquee", "maxfour", "mike", "mini", "mirror", "mnemonic", "morse", "moscow", "nancyj-fancy", "nancyj-underlined", "nancyj", "nipples", "ntgreek", "o8", "ogre", "pawp", "peaks", "pebbles", "pepper", "poison", "puffy", "pyramid", "rectangles", "relief", "relief2", "rev", "roman", "rot13", "rounded", "rowancap", "rozzo", "runic", "runyc", "sblood", "script", "serifcap", "shadow", "short", "slant", "slide", "slscript", "small", "smisome1", "smkeyboard", "smscript", "smshadow", "smslant", "smtengwar", "speed", "stampatello", "standard", "starwars", "stellar", "stop", "straight", "tanja", "tengwar", "term", "thick", "thin", "threepoint", "ticks", "ticksslant", "tinker-toy", "tombstone", "trek", "tsalagi", "twopoint", "univers", "usaflag", "wavy", "weird"}

func GetAsciiArt(s, font string) string {
	if !fontExists(font) {
		Warning("Font " + font + " does not exist, picking a random one")
		font = pickRandomAsciiFont()
	}
	asciiArt := figure.NewFigure(s, font, false).String()
	asciiArt = color.HiCyanString(asciiArt)
	return asciiArt
}

func GetRandomAsciiArt(s string) string {
	font := pickRandomAsciiFont()
	asciiArt := figure.NewFigure(s, font, false).String()
	asciiArt = color.HiCyanString(asciiArt)
	asciiArt += color.HiBlueString("\n~ using font " + font + " ~")
	return asciiArt
}

func fontExists(font string) (exists bool) {
	for i := range possibleFonts {
		if font == possibleFonts[i] {
			exists = true
			break
		}
	}
	return exists
}

func pickRandomAsciiFont() (randomFont string) {
	randomFont = possibleFonts[time.Now().Unix()%int64(len(possibleFonts))]
	return randomFont
}
