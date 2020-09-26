package twirpapi

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func (s *Server) GetDraftStatus(ctx context.Context, req *drafto.GetDraftStatusReq) (*drafto.GetDraftStatusResp, error) {
	table, err := s.Datastore.GetTable(ctx, req.TableId)
	if err != nil {
		return nil, twirp.InternalError("failed to load table")
	}

	status := &drafto.GetDraftStatusResp{
		TableId:     req.TableId,
		CurrentPack: int32(table.CurrentPack),
		Seats:       make([]*drafto.SeatData, len(table.Seats)),
	}

	for i, seat := range table.Seats {
		status.Seats[i] = &drafto.SeatData{
			SeatId:            seat.ID,
			PackCount:         int32(len(seat.PackIDs)),
			PoolCount:         int32(len(seat.FoilCardIDs) + len(seat.NonfoilCardIDs)),
			PoolRevealedCards: []*drafto.Card{},
			PackRevealedCards: []*drafto.Card{},
		}

		if len(seat.PackIDs) > 0 {
			pack, err := s.Datastore.GetPack(ctx, seat.PackIDs[0])
			if err != nil {
				return nil, twirp.InternalError("error loading pack")
			}

			packCards := s.loadCards(pack.NonfoilCardIDs, pack.FoilCardIDs)
			for _, card := range packCards {
				if card.Dfc {
					status.Seats[i].PackRevealedCards = append(status.Seats[i].PackRevealedCards, card)
				}
			}
		}

		poolCards := s.loadCards(seat.NonfoilCardIDs, seat.FoilCardIDs)
		for _, card := range poolCards {
			if card.Dfc {
				status.Seats[i].PoolRevealedCards = append(status.Seats[i].PoolRevealedCards, card)
			}
		}
	}

	return status, nil
}
