package twirpapi

import (
	"context"
	"log"

	"github.com/twitchtv/twirp"

	"github.com/patrickwhite256/drafto/internal/auth"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func (s *Server) GetCurrentUser(ctx context.Context, req *drafto.GetCurrentUserReq) (*drafto.GetCurrentUserResp, error) {
	userID := auth.UserID(ctx)
	if userID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "not logged in")
	}

	user, err := s.Datastore.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, twirp.InternalError("error loading user")
	}

	return &drafto.GetCurrentUserResp{
		Id:        user.ID,
		Name:      user.Name,
		AvatarUrl: user.AvatarURL,
	}, nil
}
