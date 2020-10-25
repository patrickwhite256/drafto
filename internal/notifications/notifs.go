package notifications

import (
	"context"
	"fmt"

	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/internal/discord"
	"github.com/patrickwhite256/drafto/internal/socket"
)

const notifMsgTpl = `You've been passed a new pack! https://drafto.patrickwhite.io/seat/%s`

type Notifications struct {
	discord      *discord.Discord
	datastore    *datastore.Datastore
	socketServer *socket.SocketServer
}

func New(discord *discord.Discord, datastore *datastore.Datastore, socketServer *socket.SocketServer) *Notifications {
	return &Notifications{discord, datastore, socketServer}
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

func (n *Notifications) NotifyTable(ctx context.Context, tableID string) {
	n.socketServer.NotifyTopic(tableID)
}
