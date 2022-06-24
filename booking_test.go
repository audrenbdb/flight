package laresa_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"laresa"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestBookFlight(t *testing.T) {
	t.Run("A flight cannot be booked with more than three guests", func(t *testing.T) {
		booker := laresa.Booker{}
		ctx := context.Background()

		_, err := booker.BookFlight(ctx, jon, laresa.BookingParams{
			DepartureCity:   parisCity,
			DestinationCity: newYorkCity,
			Guests: []laresa.Guest{
				{"Jean"}, {"Benoît"}, {"Pierre"}, {"Marie"},
			},
		})
		if !errors.Is(err, laresa.ErrTooManyGuests) {
			t.Errorf("wantResa too many guests error, got: %v", err)
		}
	})

	t.Run("Successfull registrations", func(t *testing.T) {
		fakeID := fakeIDGenerator("xyz")

		tests := []struct {
			name string

			params laresa.BookingParams

			wantResa laresa.Reservation
		}{
			{
				name: `Monday flight departure price without guests should
					equal to the distance between the two cities multiplied by 50.7`,

				params: laresa.BookingParams{
					// monday date
					DepartureDate:   monday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
				},

				wantResa: laresa.Reservation{
					ID:              fakeID(),
					Customer:        jon,
					DepartureDate:   monday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
					Price:           parisToNewYorkDistance * 50.7,
				},
			},
			{
				name: `Tuesday flight departure price without guests should
					equal to the distance between the two cities multiplied by
					the city distance minus 10`,

				params: laresa.BookingParams{
					DepartureDate:   tuesday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
				},

				wantResa: laresa.Reservation{
					ID:              fakeID(),
					Customer:        jon,
					DepartureDate:   tuesday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
					Price:           parisToNewYorkDistance * (parisToNewYorkDistance - 10),
				},
			},
			{
				name: `Wednesday flight departure price without guests should
					equal to the distance between the two cities multiplied by
					the sum of the city's name length`,

				params: laresa.BookingParams{
					DepartureDate:   wednesday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
				},

				wantResa: laresa.Reservation{
					ID:              fakeID(),
					Customer:        jon,
					DepartureDate:   wednesday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
					Price:           parisToNewYorkDistance * float64(len(parisCity.Name)+len(newYorkCity.Name)),
				},
			},
			{
				name: `Thursday flight departure price without guests should
					equal to the distance between the two cities`,

				params: laresa.BookingParams{
					DepartureDate:   thursday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
				},

				wantResa: laresa.Reservation{
					ID:              fakeID(),
					Customer:        jon,
					DepartureDate:   thursday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
					Price:           parisToNewYorkDistance,
				},
			},
			{
				name: `Adding 2 guests to a regular flight
					should discount the end price by 20%`,

				params: laresa.BookingParams{
					DepartureDate:   thursday,
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
					Distance:        parisToNewYorkDistance,
					Guests: []laresa.Guest{
						{Name: "Bob Ross"},
						{Name: "Jon Doe"},
					},
				},

				wantResa: laresa.Reservation{
					ID:              fakeID(),
					Customer:        jon,
					DepartureDate:   thursday,
					DepartureCity:   parisCity,
					Distance:        parisToNewYorkDistance,
					DestinationCity: newYorkCity,
					Price:           parisToNewYorkDistance * 0.8,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				ctx := context.Background()
				inMemRepo := laresa.NewInMemRepo()
				booker := laresa.Booker{
					Repo:      inMemRepo,
					NewID:     fakeID,
					Publisher: laresa.NewInMemPublisher(),
				}

				// publisher is an in memory that logs reservation.
				// below catches what is logged.
				var buf bytes.Buffer
				log.SetOutput(&buf)
				defer func() {
					log.SetOutput(os.Stdin)
				}()

				booker.BookFlight(ctx, jon, test.params)
				resa, _ := inMemRepo.GetReservation(ctx, fakeID())
				if resa != test.wantResa {
					t.Errorf("want reservation: %#v, got: %#v", test.wantResa, resa)
				}

				out := buf.String()
				if want := fmt.Sprintf("%#v", resa); !strings.Contains(out, want) {
					t.Errorf("want: %s, got: %s", want, out)
				}
			})
		}
	})
}

func TestGetReservation(t *testing.T) {
	jonResa := laresa.Reservation{
		ID:              "xyz",
		Customer:        jon,
		Distance:        5000,
		Price:           5000,
		DepartureDate:   time.Now().UTC().UnixNano(),
		DepartureCity:   parisCity,
		DestinationCity: newYorkCity,
	}

	tests := []struct {
		name string

		booker laresa.Booker

		id string

		wantResa laresa.Reservation
		wantErr  error
	}{
		{
			name: "A reservation with unkown ID cannot be found",

			booker: laresa.Booker{
				Repo: laresa.NewInMemRepo(),
			},

			id: "123",

			wantErr: laresa.ErrReservationNotFound,
		},
		{
			name: "Reservation matching ID is returned",

			booker: laresa.Booker{
				Repo: laresa.NewInMemRepo(laresa.WithReservations(map[string]laresa.Reservation{
					jonResa.ID: jonResa,
				})),
			},

			id: jonResa.ID,

			wantResa: jonResa,
			wantErr:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resa, err := test.booker.GetReservation(context.Background(), test.id)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("want err: %v, got: %v", test.wantErr, err)
			}
			if resa != test.wantResa {
				t.Errorf("want resa: %#v, got: %#v", test.wantResa, resa)
			}
		})
	}
}
