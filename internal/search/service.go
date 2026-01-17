// package search

// import (
// 	"context"
// 	"log"
// 	"sync"

// 	"github.com/iramkrish/search-bff/internal/clients"
// 	model "github.com/iramkrish/search-bff/internal/modal"
// )

// type Service struct {
// 	hotels  *clients.HotelClient
// 	flights *clients.FlightClient
// 	logger  *log.Logger
// }

// func NewService(logger *log.Logger) *Service {
// 	return &Service{
// 		hotels:  clients.NewHotelClient(),
// 		flights: clients.NewFlightClient(),
// 		logger:  logger,
// 	}
// }

// func (s *Service) Search(ctx context.Context, q string) (*model.SearchResponse, error) {
// 	var (
// 		wg      sync.WaitGroup
// 		resp    model.SearchResponse
// 		errOnce error
// 	)

// 	wg.Add(2)

// 	go func() {
// 		defer wg.Done()
// 		hotels, err := s.hotels.Search(ctx, q)
// 		if err != nil {
// 			s.logger.Println("hotel search failed:", err)
// 			return
// 		}
// 		resp.Hotels = hotels
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		flights, err := s.flights.Search(ctx, q)
// 		if err != nil {
// 			s.logger.Println("flight search failed:", err)
// 			return
// 		}
// 		resp.Flights = flights
// 	}()

// 	wg.Wait()

// 	if len(resp.Hotels) == 0 && len(resp.Flights) == 0 {
// 		errOnce = ErrUpstreamUnavailable
// 	}

// 	return &resp, errOnce
// }

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
		wg      sync.WaitGroup
		hotels  []model.Hotel
		flights []model.Flight
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

	return &model.SearchResponse{
		Hotels:  hotels,
		Flights: flights,
	}, nil
}
