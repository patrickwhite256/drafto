package twirpapi

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func (s *Server) NewDraft(ctx context.Context, req *drafto.NewDraftReq) (*drafto.NewDraftResp, error) {
	table, err := s.Datastore.NewTable(ctx)
	if err != nil {
		// TODO: log error
		return nil, twirp.InternalError("error starting draft")
	}

	for i := 0; i < 8; i++ {
		pack, err := s.Packgen.GenerateStandardPack(ctx, req.GetSetCode())
		if err != nil {
			// TODO: log error
			return nil, twirp.InternalError("error generating packs")
		}

		packID, err := s.Datastore.SavePack(ctx, pack.Cards)
		if err != nil {
			// TODO: log error
			return nil, twirp.InternalError("error generating packs")
		}

		err = s.Datastore.MovePackToSeat(ctx, packID, "", table.SeatIDs[i])
		if err != nil {
			// TODO: log error
			return nil, twirp.InternalError("error distributing packs")
		}
	}

	return &drafto.NewDraftResp{
		TableId: table.ID,
		SeatIds: table.SeatIDs,
	}, nil
}
