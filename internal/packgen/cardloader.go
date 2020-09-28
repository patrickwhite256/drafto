package packgen

import (
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

type CardLoader struct {
	standardSets map[string]*cardSet
	allCards     map[string]*drafto.Card
}

func (g *CardLoader) CardByID(cardID string) *drafto.Card {
	return copyOf(g.allCards[cardID])
}

func (g *CardLoader) PreloadSet(cardSet string) error {
	_, err := g.loadSet(cardSet)
	return err
}