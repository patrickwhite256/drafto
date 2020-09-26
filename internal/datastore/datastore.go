package datastore

import (
	"context"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

type dErr string

func (e dErr) Error() string {
	return string(e)
}

const (
	CardNotInPack dErr = "card not in pack"
)

type Datastore struct {
}

type Table struct {
	ID          string   `json:"id"`
	SetCode     string   `json:"set_code"`
	CurrentPack int      `json:"current_pack"`
	SeatIDs     []string `json:"seat_ids"`
	Seats       []*Seat  `json:"-"`
}

type Seat struct {
	ID             string   `json:"id"`
	TableID        string   `json:"table_id"`
	PackIDs        []string `json:"pack_ids"`
	NonfoilCardIDs []string `json:"nonfoil_card_ids"`
	FoilCardIDs    []string `json:"foil_card_ids"`
}

type Pack struct {
	ID             string   `json:"id"`
	NonfoilCardIDs []string `json:"nonfoil_card_ids"`
	FoilCardIDs    []string `json:"foil_card_ids"`
}

func New() (*Datastore, error) {
	return &Datastore{}, nil
}

// NewTable must generate a table ID and nSeats seat IDs
func (d *Datastore) NewTable(ctx context.Context, nSeats int, setCode string) (*Table, error) {
	return nil, nil
}

func (d *Datastore) GetTable(ctx context.Context, tableID string) (*Table, error) {
	return nil, nil
}

func (d *Datastore) NewPack(ctx context.Context, cards []*drafto.Card) (string, error) {
	return "", nil
}

func (d *Datastore) GetPack(ctx context.Context, packID string) (*Pack, error) {
	return nil, nil
}

func (d *Datastore) GetSeat(ctx context.Context, seatID string) (*Seat, error) {
	return nil, nil
}

func (d *Datastore) GetSeats(ctx context.Context, seatIDs []string) ([]*Seat, error) {
	return nil, nil
}

func (d *Datastore) AddCardToPool(ctx context.Context, seatID, cardID string, foil bool) error {
	return nil
}

// RemoveCardFromPack will delete the pack if it is empty
func (d *Datastore) RemoveCardFromPack(ctx context.Context, packID, cardID string) (*drafto.Card, *Pack, error) {
	return nil, nil, nil
}

// MovePackToSeat accepts unset `oldSeatID`
func (d *Datastore) MovePackToSeat(ctx context.Context, packID, oldSeatID, newSeatID string) error {
	return nil
}
