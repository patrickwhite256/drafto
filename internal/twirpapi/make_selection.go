package twirpapi

import (
	"context"
	"log"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

const FinalPack = 2

func (s *Server) MakeSelection(ctx context.Context, req *drafto.MakeSelectionReq) (*drafto.MakeSelectionResp, error) {
	seat, err := s.Datastore.GetSeat(ctx, req.SeatId)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error loading seat")
	}

	if len(seat.PackIDs) == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "seat has no packs")
	}

	table, err := s.Datastore.GetTable(ctx, seat.TableID)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error loading table")
	}

	packID := seat.PackIDs[0]

	foil, pack, err := s.Datastore.RemoveCardFromPack(ctx, packID, req.CardId)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error picking card")
	}

	err = s.Datastore.AddCardToPool(ctx, req.SeatId, req.CardId, foil)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error adding card to pool")
	}

	resp := &drafto.MakeSelectionResp{SeatId: seat.ID}

	// if pack's not empty, pass it and return
	if !pack.Empty() {
		err = s.Datastore.MovePackToSeat(ctx, packID, req.SeatId, nextSeatID(table, req.SeatId))
		if err != nil {
			log.Println(err)
			return nil, twirp.InternalError("error passing pack")
		}

		return resp, nil
	}

	// if pack is empty, clean it up
	err = s.Datastore.MovePackToSeat(ctx, packID, req.SeatId, "")
	if err != nil {
		return nil, twirp.InternalError("error cleaning up empty pack")
	}

	// check if the draft phase is done
	for _, seat := range table.Seats {
		if seat.ID == req.SeatId {
			if len(seat.PackIDs) > 1 { // table was loaded before pack was removed
				return resp, nil
			}
		} else {
			if len(seat.PackIDs) > 0 {
				return resp, nil
			}
		}
	}

	if table.CurrentPack == FinalPack {
		return resp, nil
	}

	if err := s.Datastore.IncrementTableCurrentPack(ctx, table.ID); err != nil {
		return nil, twirp.InternalError("error finishing draft phase")
	}

	return resp, s.distributeNewPacks(ctx, table)
}

// assumption: seatID is at table
func nextSeatID(table *datastore.Table, seatID string) string {
	for i, id := range table.SeatIDs {
		if id == seatID {
			nextSeatIndex := (i + 1) % len(table.SeatIDs)
			if table.CurrentPack%2 == 0 {
				nextSeatIndex = (i + len(table.SeatIDs) - 1) % len(table.SeatIDs)
			}

			return table.SeatIDs[nextSeatIndex]
		}
	}

	return ""
}
