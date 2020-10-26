package twirpapi

import (
	"context"
	"fmt"
	"log"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/internal/packgen"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func (s *Server) NewDraft(ctx context.Context, req *drafto.NewDraftReq) (*drafto.NewDraftResp, error) {
	table, err := s.Datastore.NewTable(ctx, int(req.PlayerCount), req.SetCode)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error starting draft")
	}

	table.DraftMode = int(req.DraftMode)

	// if cube: load cube, make sure count is >= n * 15 * 3
	if req.DraftMode == drafto.DraftMode_CUBE {
		cubeCardIDs, err := packgen.LoadCardIDsForCube(ctx, req.CubeId)
		if err != nil {
			log.Println(err)
			return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("failed to load cube %s", req.CubeId))
		}

		minCardCount := int(req.PlayerCount * 45)
		if len(cubeCardIDs) < minCardCount {
			return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("not enough cards in cube %s (need %d)", req.CubeId, minCardCount))
		}

		// no need to persist this - it will be saved as part of distributeNewPacks
		table.CubeUnusedIDs = cubeCardIDs
	}

	if err = s.distributeNewPacks(ctx, table); err != nil {
		return nil, err
	}

	return &drafto.NewDraftResp{
		TableId: table.ID,
		SeatIds: table.SeatIDs,
	}, nil
}
