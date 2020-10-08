package twirpapi

import (
	"context"
	"log"
	"math/rand"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/internal/auth"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func (s *Server) TakeSeat(ctx context.Context, req *drafto.TakeSeatReq) (*drafto.TakeSeatResp, error) {
	userID := auth.UserID(ctx)
	if userID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "You must be logged in to take a seat.")
	}

	table, err := s.Datastore.GetTable(ctx, req.TableId)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error loading table")
	}

	openSeatIDs := []string{}
	for _, seat := range table.Seats {
		if seat.UserID == userID {
			return nil, twirp.NewError(twirp.InvalidArgument, "You are already in this draft.")
		}

		if seat.UserID == "" {
			openSeatIDs = append(openSeatIDs, seat.ID)
		}
	}

	if len(openSeatIDs) == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "There are no open seats at this table.")
	}

	seatID := openSeatIDs[rand.Intn(len(openSeatIDs))]

	err = s.Datastore.AssignUserToSeat(ctx, userID, seatID)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error assigning player to seat")
	}

	return &drafto.TakeSeatResp{
		TableId: req.TableId,
		SeatId:  seatID,
	}, nil
}
