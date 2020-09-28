package twirpapi

import (
	"context"
	"log"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/internal/datastore"
)

func (s *Server) distributeNewPacks(ctx context.Context, table *datastore.Table) error {
	for i := 0; i < len(table.SeatIDs); i++ {
		pack, err := s.CardLoader.GenerateStandardPack(ctx, table.SetCode)
		if err != nil {
			log.Println(err)
			return twirp.InternalError("error generating packs")
		}

		packID, err := s.Datastore.NewPack(ctx, pack.Cards)
		if err != nil {
			log.Println(err)
			return twirp.InternalError("error generating packs")
		}

		err = s.Datastore.MovePackToSeat(ctx, packID, "", table.SeatIDs[i])
		if err != nil {
			log.Println(err)
			return twirp.InternalError("error distributing packs")
		}
	}

	return nil
}
