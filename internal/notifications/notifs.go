package notifications

import (
	"context"
	"fmt"

	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/internal/discord"
)

const notifMsgTpl = `You've been passed a new pack! https://drafto.patrickwhite.io/seat/%s`

type Notifications struct {
	discord   *discord.Discord
	datastore *datastore.Datastore
}

func New(discord *discord.Discord, datastore *datastore.Datastore) *Notifications {
	return &Notifications{discord, datastore}
}

func (n *Notifications) NotifyUserOfPendingPick(ctx context.Context, userID, seatID string) error {
	user, err := n.datastore.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error loading user: %w", err)
	}

	err = n.discord.SendMessage(ctx, user.DiscordID, fmt.Sprintf(notifMsgTpl, seatID))
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
