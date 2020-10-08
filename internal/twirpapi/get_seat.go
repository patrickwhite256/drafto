package twirpapi

import (
	"context"
	"log"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/internal/auth"
	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func (s *Server) GetSeat(ctx context.Context, req *drafto.GetSeatReq) (*drafto.GetSeatResp, error) {
	seat, err := s.Datastore.GetSeat(ctx, req.SeatId)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error loading seat")
	}

	if seat.UserID == "" || auth.UserID(ctx) != seat.UserID {
		return nil, twirp.NewError(twirp.PermissionDenied, "This is not your seat!")
	}

	resp := &drafto.GetSeatResp{
		SeatId:  seat.ID,
		TableId: seat.TableID,
	}

	if len(seat.PackIDs) > 0 {
		pack, err := s.Datastore.GetPack(ctx, seat.PackIDs[0])
		if err != nil {
			log.Println(err)
			return nil, twirp.InternalError("error loading current pack")
		}

		resp.CurrentPack = &drafto.Pack{
			Id:    pack.ID,
			Cards: s.loadCards(pack.NonfoilCardIDs, pack.FoilCardIDs),
		}
	}

	resp.Pool = s.loadCards(seat.NonfoilCardIDs, seat.FoilCardIDs)

	return resp, nil
}

func (s *Server) packFromDatastorePack(datastorePack datastore.Pack) *drafto.Pack {
	pack := &drafto.Pack{Id: datastorePack.ID}

	pack.Cards = s.loadCards(datastorePack.NonfoilCardIDs, datastorePack.FoilCardIDs)

	return pack
}
