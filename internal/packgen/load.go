package packgen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/patrickwhite256/drafto/rpc/drafto"
)

const (
	allCardsBulkURL = "https://api.scryfall.com/bulk-data/default-cards"
	allCardsFile    = "all_cards.json"
)

func (l *CardLoader) Preload() error {
	log.Println("preloading cards")
	cardsFile, err := os.Open(allCardsFile)
	if err != nil {
		p := &os.PathError{}
		if errors.As(err, &p) {
			log.Println("card cache not found - downloading from scryfall...")
			return l.loadCardsFromScryfall()
		} else {
			return fmt.Errorf("error opening bulk card file: %w", err)
		}
	}

	log.Println("local card cache found...")

	defer cardsFile.Close()
	body, err := ioutil.ReadAll(cardsFile)
	if err != nil {
		return fmt.Errorf("error reading bulk card data from file: %w", err)
	}

	return l.loadCardsByData(body)
}

func (l *CardLoader) loadCardsFromScryfall() error {
	cardsFile, err := os.Create(allCardsFile)
	if err != nil {
		return fmt.Errorf("error opening cards file for writing: %w", err)
	}
	defer cardsFile.Close()

	allCardsURI, err := getAllCardsURI()
	if err != nil {
		return fmt.Errorf("error getting all cards URI: %w", err)
	}

	resp, err := http.Get(allCardsURI)
	if err != nil {
		return fmt.Errorf("error loading bulk data from scryfall: %w", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading scryfall response: %w", err)
	}

	log.Println("finished downloading cards from scryfall. writing to local cache...")

	_, err = cardsFile.Write(body)
	if err != nil {
		return fmt.Errorf("error writing cards to file: %w", err)
	}

	return l.loadCardsByData(body)
}

func (l *CardLoader) loadCardsByData(data []byte) error {
	log.Println("loaded card data, processing...")
	scryCards := []*scryfallCard{}
	if err := json.Unmarshal(data, &scryCards); err != nil {
		return fmt.Errorf("error unmarshaling bulk data as json: %w", err)
	}

	log.Println("finished converting cards...")

	l.sets = make(map[string]*cardSet)
	l.allCards = make(map[string]*drafto.Card)

	for _, scryCard := range scryCards {
		setCode := strings.ToLower(scryCard.Set)
		set := l.sets[setCode]
		if set == nil {
			set = &cardSet{
				setCode:         setCode,
				cardsByID:       map[string]*drafto.Card{},
				cardsByColour:   map[drafto.Colour][]*drafto.Card{},
				cardsByRarity:   map[drafto.Rarity][]*drafto.Card{},
				dfcsByRarity:    map[drafto.Rarity][]*drafto.Card{},
				nonDFCsByRarity: map[drafto.Rarity][]*drafto.Card{},
			}

			l.sets[setCode] = set
		}

		card := set.addScryfallCard(scryCard)
		l.allCards[card.Id] = card
	}

	return nil
}

func getAllCardsURI() (string, error) {
	resp, err := http.Get(allCardsBulkURL)
	if err != nil {
		return "", fmt.Errorf("error loading bulk data from scryfall: %w", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading scryfall response: %w", err)
	}

	var scryResp = struct {
		DownloadURI string `json:"download_uri"`
	}{}

	if err := json.Unmarshal(body, &scryResp); err != nil {
		return "", fmt.Errorf("error unmarshalling scryfall response: %w", err)
	}

	return scryResp.DownloadURI, nil
}
