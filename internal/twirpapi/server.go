package twirpapi

import (
	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/internal/notifications"
	"github.com/patrickwhite256/drafto/internal/packgen"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

type Server struct {
	Datastore     *datastore.Datastore
	CardLoader    *packgen.CardLoader
	Notifications *notifications.Notifications
}

func (s *Server) loadCards(nonfoilCardIDs, foilCardIDs []string) []*drafto.Card {
	cards := make([]*drafto.Card, 0, len(nonfoilCardIDs)+len(foilCardIDs))

	for _, cardID := range nonfoilCardIDs {
		cards = append(cards, s.CardLoader.CardByID(cardID))
	}

	for _, cardID := range foilCardIDs {
		card := s.CardLoader.CardByID(cardID)
		card.Foil = true
		cards = append(cards, card)
	}

	return cards
}
