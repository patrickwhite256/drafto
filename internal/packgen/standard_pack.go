// package packgen provides utilities for generating packs
// until it's shown to be necessary, this package will prioritize readbility over efficiency.
package packgen

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

type packModifier func(set *cardSet, pack []*drafto.Card)

var (
	setRules = map[string][]packModifier{
		"znr": {znrDFCModifier},
	}

	mythicChangeDate = time.Date(2020, 9, 24, 0, 0, 0, 0, time.UTC)
)

// GenerateStandardPack creates a pack from a WotC official set with the following rules:
// - the pack has 15 cards
// - the pack has three uncommons
// - each colour must be represented in the pack
//   - this is not an "official" requirement, but most draft simulators do the same
// - the pack has either a mythic or a rare
//   - if released on or after 25 September 2020 (ZNR), the mythic chance is 1/7.4
//   - if released before 25 September 2020 (ZNR), the mythic chance is 1/8
// - the pack has either nine or ten commons
//   - 1/3 packs will have a foil instead of the tenth common.
//   - this is based on the M21 forward rate.
// - the pack has one basic land
// - excluding foils, the pack may have no duplicates
//
// Pack generation does NOT take into account official collation.
// Pack collation is a WotC trade secret and is not publicly available.
//
// The special rules for the following sets are implemented:
// - ZNR: A DFC replaces one card of the appropriate rarity.
func (g *CardLoader) GenerateStandardPack(ctx context.Context, setCode string) (*drafto.Pack, error) {
	set, err := g.loadSet(setCode)
	if err != nil {
		return nil, fmt.Errorf("error loading set %s: %w", setCode, err)
	}

	pack := make([]*drafto.Card, 0, 15)

	// rare slot
	if packShouldHaveMythic(set) {
		pack = append(pack, set.randomCardsOfRarity(drafto.Rarity_MYTHIC, 1)[0])
	} else {
		pack = append(pack, set.randomCardsOfRarity(drafto.Rarity_RARE, 1)[0])
	}

	// uncommons
	pack = append(pack, set.randomCardsOfRarity(drafto.Rarity_UNCOMMON, 3)...)

	// common 1 - possible foil
	if packShouldHaveFoil() {
		card := randomCardsFromCandidates(set.cards, 1)[0]
		card.Foil = true
		pack = append(pack, card)
	} else {
		pack = append(pack, set.randomCardsOfRarity(drafto.Rarity_COMMON, 1)[0])
	}

	// commons 2-10
	pack = append(pack, set.randomCardsOfRarity(drafto.Rarity_COMMON, 9)...)

	// basic
	pack = append(pack, set.randomCardsOfRarity(drafto.Rarity_BASIC, 1)[0])

	// apply any set-specific rules
	for _, mod := range setRules[strings.ToLower(setCode)] {
		mod(set, pack)
	}

	// adjust pack if necessary to fix colour balance
	fixColourBalance(set, pack)

	return &drafto.Pack{Cards: pack}, nil
}

func packShouldHaveMythic(set *cardSet) bool {
	if set.releaseDate.After(mythicChangeDate) {
		return rand.Float64() < 1.0/7.4
	}

	return rand.Float64() < 1.0/8.0
}

func packShouldHaveFoil() bool {
	return rand.Float64() < 1.0/3.0
}

// fixColourBalance ensures a pack has at least one card of every colour by replacing commons if necessary.
func fixColourBalance(set *cardSet, pack []*drafto.Card) {
	// don't replace the first common of each colour
	representedColourSet := map[drafto.Colour]struct{}{}
	replaceableCommonIndices := []int{}

	for i, card := range pack {
		replaceable := card.Rarity == drafto.Rarity_COMMON && !card.Foil
		for _, colour := range card.Colours {
			if _, ok := representedColourSet[colour]; !ok {
				representedColourSet[colour] = struct{}{}
				replaceable = false
			}
		}

		if replaceable {
			replaceableCommonIndices = append(replaceableCommonIndices, i)
		}
	}

	for i, colour := range ALL_COLOURS {
		if _, ok := representedColourSet[colour]; ok {
			continue
		}

		replaceIndexWithColourBalancingCard(replaceableCommonIndices[i], set, pack, colour)
	}
}

func replaceIndexWithColourBalancingCard(idx int, set *cardSet, pack []*drafto.Card, colour drafto.Colour) {
	invalidIDSet := make(map[string]struct{}, len(pack))

	for _, card := range pack {
		if card.Foil {
			continue
		}

		invalidIDSet[card.Id] = struct{}{}
	}

	candidates := make([]*drafto.Card, 0, len(set.cardsByColour[colour])-len(invalidIDSet))

	for _, card := range set.cardsByColour[colour] {
		if card.Rarity == drafto.Rarity_COMMON {
			candidates = append(candidates, card)
		}
	}

	pack[idx] = randomCardsFromCandidates(candidates, 1)[0]
}
