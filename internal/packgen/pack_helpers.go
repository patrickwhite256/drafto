package packgen

import (
	"math/rand"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

// replaceNonfoilCardOfRarity replaces the first nonfoil card of rarity `rarity` in
//   pack `pack` with card `card`. If no card of the appropriate rarity is found, no-op.
func replaceNonfoilCardOfRarity(pack []*drafto.Card, rarity drafto.Rarity, card *drafto.Card) {
	for i, card := range pack {
		if !card.Foil && card.Rarity == rarity {
			pack[i] = card
			return
		}
	}
}

func randomCardsFromCandidates(candidates []*drafto.Card, count int) []*drafto.Card {
	results := make([]*drafto.Card, count)

	perm := rand.Perm(len(candidates))

	for i := 0; i < count; i++ {
		results[i] = copyOf(candidates[perm[i]])
	}

	return results
}

func (s *cardSet) randomCardsOfRarity(rarity drafto.Rarity, count int) []*drafto.Card {
	return randomCardsFromCandidates(s.cardsByRarity[rarity], count)
}
