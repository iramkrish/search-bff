package clients

import (
	"context"
	"time"

	model "github.com/iramkrish/search-bff/internal/modal"
)

type HotelClient struct{}

func NewHotelClient() *HotelClient {
	return &HotelClient{}
}

func (c *HotelClient) Search(ctx context.Context, q string) ([]model.Hotel, error) {
	select {
	case <-time.After(100 * time.Millisecond):
		return []model.Hotel{
			{ID: "h1", Name: "Grand Hotel"},
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
