package packgen

import (
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

type CardLoader struct {
	sets     map[string]*cardSet
	allCards map[string]*drafto.Card
}

func (g *CardLoader) CardByID(cardID string) *drafto.Card {
	return copyOf(g.allCards[cardID])
}

func (g *CardLoader) CardsByIDs(cardIDs []string) []*drafto.Card {
	cards := make([]*drafto.Card, len(cardIDs))
	for i, id := range cardIDs {
		cards[i] = g.CardByID(id)

	}

	return cards
}
