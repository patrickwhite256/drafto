package datastore

import (
	"context"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

type Datastore struct {
}

type Table struct {
	ID          string   `json:"id"`
	CurrentPack int      `json:"current_pack"`
	SeatIDs     []string `json:"seat_ids"`
}

type Seat struct {
	ID      string   `json:"id"`
	PackIDs []string `json:"pack_ids"`
	CardIDs []string `json:"card_ids"`
}

type Pack struct {
	ID          string   `json:"id"`
	CardIDs     []string `json:"card_ids"`
	FoilIndices []int    `json:"foil_indices"`
}

func New() (*Datastore, error) {
	return &Datastore{}, nil
}

// NewTable must generate a table ID and eight seat IDs
func (d *Datastore) NewTable(ctx context.Context) (*Table, error) {
	return nil, nil
}

func (d *Datastore) SavePack(ctx context.Context, cards []*drafto.Card) (string, error) {
	return "", nil
}

func (d *Datastore) GetSeat(ctx context.Context, seatID string) (*Seat, error) {
	return nil, nil
}

func (d *Datastore) GetSeatsForTable(ctx context.Context, tableID string) ([]*Seat, error) {
	return nil, nil
}

// RemoveCardFromPack will delete the pack if it is empty
func (d *Datastore) RemoveCardFromPack(ctx context.Context, packID, cardID string) error {
	return nil
}

// MovePackToSeat accepts unset `oldSeatID`
func (d *Datastore) MovePackToSeat(ctx context.Context, packID, oldSeatID, newSeatID string) error {
	return nil
}
