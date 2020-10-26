package twirpapi

import (
	"context"
	"errors"
	"log"
	"math/rand"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func (s *Server) distributeNewPacks(ctx context.Context, table *datastore.Table) error {
	packs := make([]*drafto.Pack, len(table.SeatIDs))
	for i := 0; i < len(table.SeatIDs); i++ {
		log.Println("generating pack")
		pack, err := s.generatePack(ctx, table)
		if err != nil {
			log.Println(err)
			return twirp.InternalError("error generating packs")
		}

		packs[i] = pack
	}

	if drafto.DraftMode(table.DraftMode) == drafto.DraftMode_CUBE {
		if err := s.Datastore.WriteTable(ctx, table); err != nil {
			return twirp.InternalError("error generating cube packs")
		}
	}

	for i := 0; i < len(table.SeatIDs); i++ {
		packID, err := s.Datastore.NewPack(ctx, packs[i].Cards)
		if err != nil {
			log.Println(err)
			return twirp.InternalError("error creating packs")
		}

		err = s.Datastore.MovePackToSeat(ctx, packID, "", table.SeatIDs[i])
		if err != nil {
			log.Println(err)
			return twirp.InternalError("error distributing packs")
		}
	}

	return nil
}

func (s *Server) generatePack(ctx context.Context, table *datastore.Table) (*drafto.Pack, error) {
	switch drafto.DraftMode(table.DraftMode) {
	case drafto.DraftMode_PACK:
		return s.CardLoader.GenerateStandardPack(ctx, table.SetCode)
	case drafto.DraftMode_CUBE:
		log.Println("shuffling IDs")
		rand.Shuffle(len(table.CubeUnusedIDs), func(i, j int) {
			table.CubeUnusedIDs[i], table.CubeUnusedIDs[j] = table.CubeUnusedIDs[j], table.CubeUnusedIDs[i]
		})

		log.Println("shuffled IDs")

		packIDs := table.CubeUnusedIDs[:15]
		table.CubeUnusedIDs = table.CubeUnusedIDs[15:]
		return &drafto.Pack{Cards: s.CardLoader.CardsByIDs(packIDs)}, nil
	}

	return nil, errors.New("unknown draft type")
}
