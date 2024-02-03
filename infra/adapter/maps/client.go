package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/frisk038/swipe_dungeon/business/models"
)

type Client struct {
	client *http.Client
	apiKey string
}

func New() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		apiKey: os.Getenv("MAPS_API_KEY"),
	}
}

func (ct *Client) GetCity(ctx context.Context, loc models.Location) (string, error) {
	reqStr := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%s,%s&result_type=political&key=%s", loc.Latitude, loc.Longitude, ct.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqStr, nil)
	if err != nil {
		return "", err
	}

	resp, err := ct.client.Do(req)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		st := struct {
			Result []struct {
				AddrCompo []struct {
					ShortName string `json:"short_name"`
				} `json:"address_components"`
			} `json:"results"`
		}{}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		if err = json.Unmarshal(body, &st); err != nil {
			return "", err
		}

		if len(st.Result) == 0 {
			return "", err
		}

		return st.Result[0].AddrCompo[0].ShortName, nil
	default:
		return resp.Status, err
	}
}
