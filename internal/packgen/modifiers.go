package packgen

import (
	"math/rand"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

// znrDFCModifier enforces that a pack has exactly one nonfoil DFC.
// - if the pack already has exactly one DFC, nothing will happen.
// - if the pack has more than one nonfoil DFC, the highest-rarity nonfoil DFC will be kept and the others will be replaced with equivalent-rarity non-DFCs
// - if the pack has zero nonfoil DFCs, one of the uncommons will be replaced.
func znrDFCModifier(set *cardSet, pack []*drafto.Card) {
	dfcIndices := []int{}
	highestRarityDFCIndex := -1

	for i, card := range pack {
		if card.Dfc {
			dfcIndices = append(dfcIndices, i)

			if highestRarityDFCIndex == -1 || card.Rarity > pack[highestRarityDFCIndex].Rarity {
				highestRarityDFCIndex = i
			}
		}
	}

	if len(dfcIndices) == 1 {
		return
	}

	if len(dfcIndices) == 0 {
		card := copyOf(set.dfcsByRarity[drafto.Rarity_UNCOMMON][rand.Intn(len(set.dfcsByRarity[drafto.Rarity_UNCOMMON]))])
		replaceNonfoilCardOfRarity(pack, drafto.Rarity_UNCOMMON, card)
		return
	}

	for _, idx := range dfcIndices {
		if idx == highestRarityDFCIndex {
			continue
		}

		pack[idx] = set.randomZNRSFC(pack[idx].Rarity, pack)
	}
}

// randomZNRDFC returns a card of rarity `rarity`, that
// is not a DFC, or any card already in `pack`, unless it is foil.
func (s *cardSet) randomZNRSFC(rarity drafto.Rarity, pack []*drafto.Card) *drafto.Card {
	invalidIDSet := make(map[string]struct{}, len(pack))

	for _, card := range pack {
		if card.Foil {
			continue
		}

		invalidIDSet[card.Id] = struct{}{}
	}

	candidates := make([]*drafto.Card, 0, len(s.cardsByRarity[rarity])-len(invalidIDSet))

	for _, card := range s.cardsByRarity[rarity] {
		if _, ok := invalidIDSet[card.Id]; ok || card.Dfc {
			continue
		}

		candidates = append(candidates, card)
	}

	return copyOf(candidates[rand.Intn(len(candidates))])
}
