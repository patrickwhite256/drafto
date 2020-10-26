package packgen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const cubeListURL = "https://cubecobra.com/cube/list/"

func LoadCardIDsForCube(ctx context.Context, cubeID string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", cubeListURL+cubeID, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status from cubecobra: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading cubecobra response: %w", err)
	}

	// it's not quite JSON, so we need to massage it a little
	re := regexp.MustCompile(`new Date\([^\)]+\)`)

	parts := bytes.Split(bytes.Split(body, []byte(";"))[0], []byte("reactProps = "))
	if len(parts) != 2 {
		return nil, fmt.Errorf("bad response from cubecobra")
	}

	processed := re.ReplaceAllLiteral(parts[1], []byte("\"\""))

	format := struct {
		Cube struct {
			Cards []struct {
				ID string `json:"cardID"`
			} `json:"cards"`
		} `json:"cube"`
	}{}

	if err := json.Unmarshal(processed, &format); err != nil {
		return nil, fmt.Errorf("could not process json from cubecobra, %w", err)
	}

	ids := make([]string, len(format.Cube.Cards))
	for i, card := range format.Cube.Cards {
		ids[i] = card.ID
	}

	return ids, nil
}
