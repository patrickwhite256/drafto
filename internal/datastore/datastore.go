package datastore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

type dErr string

func (e dErr) Error() string {
	return string(e)
}

const (
	CardNotInPack dErr = "card not in pack"
	NotFound      dErr = "not found"

	tableTableName = "drafto-tables"
	seatTableName  = "drafto-seats"
	packTableName  = "drafto-packs"
)

type Datastore struct {
	ddb *dynamodb.DynamoDB
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

func (p *Pack) Empty() bool {
	return len(p.FoilCardIDs)+len(p.NonfoilCardIDs) == 0
}

func New() (*Datastore, error) {
	sess := session.Must(session.NewSession(&aws.Config{}))
	return &Datastore{
		ddb: dynamodb.New(sess),
	}, nil
}

// NewTable must generate a table ID and nSeats seat IDs
func (d *Datastore) NewTable(ctx context.Context, nSeats int, setCode string) (*Table, error) {
	table := &Table{
		ID:          uuid.New().String(),
		SetCode:     setCode,
		CurrentPack: 0,
	}

	seatWriteRequests := make([]*dynamodb.WriteRequest, 0, nSeats)

	for i := 0; i < nSeats; i++ {
		seat := &Seat{
			ID:      uuid.New().String(),
			TableID: table.ID,
		}

		table.SeatIDs = append(table.SeatIDs, seat.ID)
		table.Seats = append(table.Seats, seat)

		item, err := dynamodbattribute.MarshalMap(seat)
		if err != nil {
			return nil, fmt.Errorf("error marshaling seat: %w", err)
		}

		seatWriteRequests = append(seatWriteRequests, &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: item}})
	}

	item, err := dynamodbattribute.MarshalMap(table)
	if err != nil {
		return nil, fmt.Errorf("error marshalling table: %w", err)
	}

	if _, err = d.ddb.BatchWriteItemWithContext(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			tableTableName: {{PutRequest: &dynamodb.PutRequest{Item: item}}},
			seatTableName:  seatWriteRequests,
		},
	}); err != nil {
		return nil, fmt.Errorf("error writing to dynamo: %w", err)
	}

	return table, nil
}

func (d *Datastore) GetTable(ctx context.Context, tableID string) (*Table, error) {
	table := &Table{}

	var err error

	if err = d.loadItem(ctx, tableID, tableTableName, table); err != nil {
		return nil, err
	}

	table.Seats, err = d.GetSeats(ctx, table.SeatIDs)
	if err != nil {
		return nil, fmt.Errorf("unable to load table seats: %w", err)
	}

	return table, nil
}

func (d *Datastore) NewPack(ctx context.Context, cards []*drafto.Card) (string, error) {
	pack := &Pack{
		ID: uuid.New().String(),
	}

	for _, card := range cards {
		if card.Foil {
			pack.FoilCardIDs = append(pack.FoilCardIDs, card.Id)
			continue
		}

		pack.NonfoilCardIDs = append(pack.NonfoilCardIDs, card.Id)
	}

	if err := d.writeItem(ctx, pack, packTableName); err != nil {
		return "", fmt.Errorf("error writing pack: %w", err)
	}

	return pack.ID, nil
}

func (d *Datastore) GetPack(ctx context.Context, packID string) (*Pack, error) {
	pack := &Pack{}

	if err := d.loadItem(ctx, packID, packTableName, pack); err != nil {
		return nil, err
	}

	return pack, nil
}

func (d *Datastore) GetSeat(ctx context.Context, seatID string) (*Seat, error) {
	seat := &Seat{}

	if err := d.loadItem(ctx, seatID, seatTableName, seat); err != nil {
		return nil, err
	}

	return seat, nil
}

func (d *Datastore) GetSeats(ctx context.Context, seatIDs []string) ([]*Seat, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(seatIDs))
	for _, id := range seatIDs {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}})
	}

	resp, err := d.ddb.BatchGetItemWithContext(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			seatTableName: {
				Keys: keys,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error reading from dynamo: %w", err)
	}

	seats := []*Seat{}

	if err = dynamodbattribute.UnmarshalListOfMaps(resp.Responses[seatTableName], &seats); err != nil {
		return nil, fmt.Errorf("error unmarshalling seats: %w", err)
	}

	return seats, nil
}

func (d *Datastore) AddCardToPool(ctx context.Context, seatID, cardID string, foil bool) error {
	seat, err := d.GetSeat(ctx, seatID)
	if err != nil {
		return fmt.Errorf("error getting seat: %w", err)
	}

	if foil {
		seat.FoilCardIDs = append(seat.FoilCardIDs, cardID)
	} else {
		seat.NonfoilCardIDs = append(seat.NonfoilCardIDs, cardID)
	}

	if err := d.writeItem(ctx, seat, seatTableName); err != nil {
		return fmt.Errorf("error writing seat: %w", err)
	}

	return nil
}

// RemoveCardFromPack will delete the pack if it is empty
// Returns NotFound if the card is not in the pack
func (d *Datastore) RemoveCardFromPack(ctx context.Context, packID, cardID string) (isFoil bool, pack *Pack, err error) {
	pack, err = d.GetPack(ctx, packID)
	if err != nil {
		return false, nil, fmt.Errorf("error getting pack: %w", err)
	}

	found := false
	for i, id := range pack.FoilCardIDs {
		if id == cardID {
			found = true
			isFoil = true
			pack.FoilCardIDs = append(pack.FoilCardIDs[:i], pack.FoilCardIDs[i+1:]...)
		}
	}

	if !found {
		for i, id := range pack.NonfoilCardIDs {
			if id == cardID {
				found = true
				pack.NonfoilCardIDs = append(pack.NonfoilCardIDs[:i], pack.NonfoilCardIDs[i+1:]...)
			}
		}
	}

	if !found {
		return false, nil, NotFound
	}

	if pack.Empty() {
		// delete pack
		if _, err := d.ddb.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(packTableName),
			Key: map[string]*dynamodb.AttributeValue{
				"id": {S: aws.String(pack.ID)},
			},
		}); err != nil {
			return false, nil, fmt.Errorf("error deleting pack: %w", err)
		}
	} else {
		// write pack back
		if err := d.writeItem(ctx, pack, packTableName); err != nil {
			return false, nil, fmt.Errorf("error writing pack: %w", err)
		}
	}

	return isFoil, pack, nil
}

// MovePackToSeat accepts unset `oldSeatID`
func (d *Datastore) MovePackToSeat(ctx context.Context, packID, oldSeatID, newSeatID string) error {
	if oldSeatID != "" {
		seat := &Seat{}

		if err := d.loadItem(ctx, oldSeatID, seatTableName, seat); err != nil {
			return err
		}

		if seat.PackIDs[0] != packID {
			return NotFound
		}

		seat.PackIDs = seat.PackIDs[1:]

		if err := d.writeItem(ctx, seat, seatTableName); err != nil {
			return err
		}
	}

	seat := &Seat{}

	if err := d.loadItem(ctx, newSeatID, seatTableName, seat); err != nil {
		return err
	}

	seat.PackIDs = append(seat.PackIDs, packID)

	return d.writeItem(ctx, seat, seatTableName)
}

func (d *Datastore) writeItem(ctx context.Context, item interface{}, tableName string) error {
	dynamoItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("error marshalling item: %w", err)
	}

	if _, err := d.ddb.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      dynamoItem,
	}); err != nil {
		return fmt.Errorf("error writing pack: %w", err)
	}

	return nil
}

func (d *Datastore) loadItem(ctx context.Context, id, tableName string, item interface{}) error {
	resp, err := d.ddb.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
	})
	if err != nil {
		return fmt.Errorf("error reading from dynamo: %w", err)
	}

	if resp.Item == nil {
		return NotFound
	}

	err = dynamodbattribute.UnmarshalMap(resp.Item, item)
	if err != nil {
		return fmt.Errorf("unable to unmarshal item: %w", err)
	}

	return nil
}
