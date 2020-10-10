package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

const (
	discordAPIBase = "https://discord.com/api/v8"

	createDMEndpoint      = "/users/@me/channels"
	createMessageEndpoint = "/channels/%s/messages"
)

type Discord struct {
	sess *discordgo.Session
}

func New(botToken string) (*Discord, error) {
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	if err = dg.Open(); err != nil {
		return nil, err
	}

	return &Discord{sess: dg}, nil
}

func (d *Discord) SendMessage(ctx context.Context, userID, message string) error {
	channel, err := d.sess.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	_, err = d.sess.ChannelMessageSend(channel.ID, message)
	return err
}
