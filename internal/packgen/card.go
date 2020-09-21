package packgen

import (
	"fmt"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

var colourNames = map[drafto.Colour]string{
	drafto.Colour_WHITE: "W",
	drafto.Colour_BLUE:  "U",
	drafto.Colour_BLACK: "B",
	drafto.Colour_RED:   "R",
	drafto.Colour_GREEN: "G",
}

var ALL_COLOURS = []drafto.Colour{
	drafto.Colour_WHITE,
	drafto.Colour_BLUE,
	drafto.Colour_BLACK,
	drafto.Colour_RED,
	drafto.Colour_GREEN,
}

func copyOf(c *drafto.Card) *drafto.Card {
	return &drafto.Card{
		Id:       c.Id,
		Name:     c.Name,
		ImageUrl: c.ImageUrl,
		Colours:  c.Colours,
		Rarity:   c.Rarity,
		Foil:     c.Foil,
		Dfc:      c.Dfc,
	}
}

func CardString(c *drafto.Card) string {
	colourString := ""
	for _, colour := range c.Colours {
		colourString += colourNames[colour]
	}

	if c.Foil {
		return fmt.Sprintf("shiny %s (%s)", c.Name, colourString)
	}

	return fmt.Sprintf("%s (%s)", c.Name, colourString)
}
