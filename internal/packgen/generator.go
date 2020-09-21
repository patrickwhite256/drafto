package packgen

import (
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

type Generator struct {
	standardSets map[string]*cardSet
}

func (g *Generator) CardByID(setCode, cardID string) *drafto.Card {
	return copyOf(g.standardSets[setCode].cardsByID[cardID])
}
