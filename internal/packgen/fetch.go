package packgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

const (
	// we're being good scryfall API citizens
	scryfallBackoffTime = 100 * time.Millisecond
	// see https://golang.org/pkg/time/#Parse
	scryfallDateFormat = "2006-01-02"
)

var rarityByName = map[string]drafto.Rarity{
	"common":   drafto.Rarity_COMMON,
	"uncommon": drafto.Rarity_UNCOMMON,
	"rare":     drafto.Rarity_RARE,
	"mythic":   drafto.Rarity_MYTHIC,
}

var colourByName = map[string]drafto.Colour{
	"W": drafto.Colour_WHITE,
	"U": drafto.Colour_BLUE,
	"B": drafto.Colour_BLACK,
	"R": drafto.Colour_RED,
	"G": drafto.Colour_GREEN,
}

func coloursFromStrings(colourStrings []string) []drafto.Colour {
	colourSet := map[drafto.Colour]struct{}{}

	for _, c := range colourStrings {
		colourSet[colourByName[c]] = struct{}{}
	}

	colours := make([]drafto.Colour, 0, len(colourSet))

	for c := range colourSet {
		colours = append(colours, c)
	}

	return colours
}

var basicNames = map[string]struct{}{
	"Plains":   {},
	"Island":   {},
	"Swamp":    {},
	"Mountain": {},
	"Forest":   {},
}

// TODO: load price data for the moneybags bot
type scryfallCard struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	RarityString string             `json:"rarity"`
	ImageURIs    map[string]string  `json:"image_uris"`
	Colours      []string           `json:"colors"`
	Faces        []scryfallCardFace `json:"card_faces"`
	ReleaseDate  string             `json:"released_at"`
}

func (c scryfallCard) Rarity() drafto.Rarity {
	if _, ok := basicNames[c.Name]; ok {
		return drafto.Rarity_BASIC
	}

	return rarityByName[c.RarityString]
}

type scryfallCardFace struct {
	Colours   []string          `json:"colors"`
	ImageURIs map[string]string `json:"image_uris"`
}

func (c scryfallCard) toCard() *drafto.Card {
	card := &drafto.Card{
		Id:       c.ID,
		Name:     c.Name,
		ImageUrl: c.ImageURIs["normal"],
		Colours:  coloursFromStrings(c.Colours),
		Rarity:   c.Rarity(),
	}

	if len(c.Faces) == 0 {
		return card
	}

	// this is potentially destructive to the first face's colours - we don't use those for anything else though
	card.Colours = coloursFromStrings(append(c.Faces[0].Colours, c.Faces[1].Colours...))
	card.ImageUrl = c.Faces[0].ImageURIs["normal"]
	card.Dfc = true

	return card
}

type scryfallSearchResponse struct {
	Data        []*scryfallCard `json:"data"`
	NextPageURL string          `json:"next_page"`
}

func (s *cardSet) addScryfallCard(scryCard *scryfallCard) *drafto.Card {
	if s.releaseDate.IsZero() {
		var err error
		if s.releaseDate, err = time.Parse(scryfallDateFormat, scryCard.ReleaseDate); err != nil {
			// TODO: log parse error
			s.releaseDate = time.Time{}
		}
	}

	card := scryCard.toCard()
	s.cards = append(s.cards, card)

	for _, colour := range card.Colours {
		s.cardsByColour[colour] = append(s.cardsByColour[colour], card)
	}

	s.cardsByRarity[card.Rarity] = append(s.cardsByRarity[card.Rarity], card)

	if card.Dfc {
		s.dfcsByRarity[card.Rarity] = append(s.dfcsByRarity[card.Rarity], card)
	} else {
		s.nonDFCsByRarity[card.Rarity] = append(s.nonDFCsByRarity[card.Rarity], card)
	}

	s.cardsByID[card.Id] = card

	return card
}

func scryfallSetSearchURL(setCode string) string {
	return "https://api.scryfall.com/cards/search?order=set&q=e%3A" + setCode + "+in%3Abooster&unique=cards"
}

func getScryfallPage(url string) (*scryfallSearchResponse, error) {
	var scryResp *scryfallSearchResponse

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching cards from scryfall: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("scryfall error response: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading scryfall response: %w", err)
	}

	if err := json.Unmarshal(body, &scryResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling scryfall response: %w", err)
	}

	return scryResp, nil
}

func (g *CardLoader) loadSet(setCode string) (*cardSet, error) {
	setCode = strings.ToLower(setCode)
	if g.standardSets == nil {
		g.standardSets = map[string]*cardSet{}
		g.allCards = map[string]*drafto.Card{}
	}

	if set, ok := g.standardSets[setCode]; ok {
		return set, nil
	}

	// TODO: handle unknown set code
	url := scryfallSetSearchURL(setCode)

	set := &cardSet{
		setCode:         setCode,
		cardsByID:       map[string]*drafto.Card{},
		cardsByColour:   map[drafto.Colour][]*drafto.Card{},
		cardsByRarity:   map[drafto.Rarity][]*drafto.Card{},
		dfcsByRarity:    map[drafto.Rarity][]*drafto.Card{},
		nonDFCsByRarity: map[drafto.Rarity][]*drafto.Card{},
	}

	for url != "" {
		scryResp, err := getScryfallPage(url)
		if err != nil {
			return nil, err
		}

		for _, scryCard := range scryResp.Data {
			card := set.addScryfallCard(scryCard)
			g.allCards[card.Id] = card
		}

		url = scryResp.NextPageURL

		time.Sleep(scryfallBackoffTime)
	}

	g.standardSets[setCode] = set

	return set, nil
}
