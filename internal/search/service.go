package search

import (
	"context"
	"log"
	"sync"

	"github.com/iramkrish/search-bff/internal/clients"
	model "github.com/iramkrish/search-bff/internal/modal"
)

type Service struct {
	hotels  *clients.HotelClient
	flights *clients.FlightClient
	logger  *log.Logger
}

func NewService(logger *log.Logger) *Service {
	return &Service{
		hotels:  clients.NewHotelClient(),
		flights: clients.NewFlightClient(),
		logger:  logger,
	}
}

func (s *Service) Search(ctx context.Context, q string) (*model.SearchResponse, error) {
	var (
		wg       sync.WaitGroup
		hotels   []model.Hotel
		flights  []model.Flight
		warnings []string
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		result, err := s.hotels.Search(ctx, q)
		if err != nil {
			s.logger.Println("hotel search failed:", err)
			return
		}
		hotels = result
	}()

	go func() {
		defer wg.Done()
		result, err := s.flights.Search(ctx, q)
		if err != nil {
			s.logger.Println("flight search failed:", err)
			return
		}
		flights = result
	}()

	wg.Wait()

	if len(hotels) == 0 && len(flights) == 0 {
		return nil, ErrUpstreamUnavailable
	}

	if len(hotels) == 0 {
		warnings = append(warnings, "hotels unavailable")
	}
	if len(flights) == 0 {
		warnings = append(warnings, "flights unavailable")
	}

	return &model.SearchResponse{
		Hotels:   hotels,
		Flights:  flights,
		Warnings: warnings,
	}, nil
}
