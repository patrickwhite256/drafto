package packgen

import (
	"time"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

// a full official magic set.
type cardSet struct {
	setCode         string
	releaseDate     time.Time
	cards           []*drafto.Card
	cardsByID       map[string]*drafto.Card
	cardsByColour   map[drafto.Colour][]*drafto.Card
	cardsByRarity   map[drafto.Rarity][]*drafto.Card
	dfcsByRarity    map[drafto.Rarity][]*drafto.Card
	nonDFCsByRarity map[drafto.Rarity][]*drafto.Card
}
