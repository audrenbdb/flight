package laresa

import (
	"context"
	"time"
)

type Repo interface {
	SaveReservation(ctx context.Context, r Reservation) error
	GetReservation(ctx context.Context, id string) (Reservation, error)
}

type Publisher interface {
	PublishReservation(ctx context.Context, r Reservation) error
}

type Booker struct {
	Repo      Repo
	NewID     func() string
	Publisher Publisher
}

func (b *Booker) BookFlight(ctx context.Context, customer Customer, params BookingParams) (Reservation, error) {
	if isTooManyGuestsRegistered(params.Guests) {
		return Reservation{}, ErrTooManyGuests
	}
	price := calculateFlightPrice(params)
	resa := Reservation{
		ID:              b.NewID(),
		Customer:        customer,
		DepartureDate:   params.DepartureDate,
		Distance:        params.Distance,
		DepartureCity:   params.DepartureCity,
		DestinationCity: params.DestinationCity,
		Price:           price,
	}
	err := b.Repo.SaveReservation(ctx, resa)
	if err != nil {
		return Reservation{}, err
	}
	return resa, b.Publisher.PublishReservation(ctx, resa)
}

func (b *Booker) GetReservation(ctx context.Context, id string) (Reservation, error) {
	return b.Repo.GetReservation(ctx, id)
}

func calculateFlightPrice(params BookingParams) float64 {
	coefficient := calculateFlightPriceCoefficient(params)
	basePrice := params.Distance * coefficient
	return applyFlightDiscount(basePrice, len(params.Guests))
}

func applyFlightDiscount(basePrice float64, guestsInvited int) float64 {
	return basePrice * (1 - (float64(guestsInvited) * .1))
}

func calculateFlightPriceCoefficient(params BookingParams) float64 {
	switch time.Unix(0, params.DepartureDate).Weekday() {
	// monday starts at 1
	case 1:
		return 50.7
	case 2:
		return params.Distance - 10
	case 3:
		return float64(len(params.DepartureCity.Name) + len(params.DestinationCity.Name))
	default:
		return 1
	}
}

func isTooManyGuestsRegistered(guests []Guest) bool {
	return len(guests) > 3
}

type Customer struct {
	Email string `json:"email"`
}

type Guest struct {
	Name string `json:"name"`
}

type City struct {
	Name string `json:"name"`
}

type Reservation struct {
	ID       string   `json:"id"`
	Customer Customer `json:"customer"`
	// unix nano
	DepartureDate   int64   `json:"departureDate"`
	DepartureCity   City    `json:"departureCity"`
	DestinationCity City    `json:"destinationCity"`
	Distance        float64 `json:"distance"`
	Price           float64 `json:"price"`
}

type BookingParams struct {
	// unix nano
	DepartureDate   int64   `json:"departureDate"`
	DepartureCity   City    `json:"departureCity"`
	DestinationCity City    `json:"destinationCity"`
	Distance        float64 `json:"distance"`
	Guests          []Guest `json:"guests"`
}

func (p BookingParams) IsComplete() bool {
	switch {
	case p.DepartureDate == 0, p.DepartureCity.Name == "", p.DestinationCity.Name == "":
		return false
	}
	return true
}
