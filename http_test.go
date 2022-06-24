package laresa_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"laresa"
	"laresa/ulid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewChiHTTPHandler(t *testing.T) {
	testReservation(t, laresa.NewChiHTTPHandler(&laresa.Booker{
		NewID:     ulid.NewBuilder(),
		Repo:      laresa.NewInMemRepo(),
		Publisher: laresa.NewInMemPublisher(),
	}))
}

func testReservation(t *testing.T, handler http.Handler) {
	t.Run("Unsuccessful request to book a flight", func(t *testing.T) {
		tests := []struct {
			name string

			customerEmail string

			body []byte

			expectedStatus int
		}{
			{
				name: "Request with invalid body is a bad request",

				customerEmail: "bruno@gacio.org",

				body: []byte("`"),

				expectedStatus: http.StatusBadRequest,
			},
			{
				name: "Request to book a flight without departure city is a bad request",

				customerEmail: "jon@doe.org",

				body: bytesFrom(t, laresa.BookingParams{
					DestinationCity: newYorkCity,
					DepartureDate:   time.Now().UnixNano(),
				}),

				expectedStatus: http.StatusBadRequest,
			},
			{
				name: "Request to book a flight without destination city is a bad request",

				customerEmail: "bob@sap.org",

				body: bytesFrom(t, laresa.BookingParams{
					DepartureCity: parisCity,
					DepartureDate: time.Now().UnixNano(),
				}),

				expectedStatus: http.StatusBadRequest,
			},
			{
				name: "Request to book a flight without departure date is a bad request",

				customerEmail: "bob@sap.org",

				body: bytesFrom(t, laresa.BookingParams{
					DepartureCity:   parisCity,
					DestinationCity: newYorkCity,
				}),

				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				r := httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("/customers/%s/reservations", test.customerEmail),
					bytes.NewReader(test.body),
				)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				if test.expectedStatus != w.Code {
					t.Errorf("want status: %d, got: %d", test.expectedStatus, w.Code)
				}
			})
		}
	})

	t.Run("A newly reserved flight should be accessible via a GET endpoint", func(t *testing.T) {
		customer := jon
		body := laresa.BookingParams{
			DepartureCity:   parisCity,
			DestinationCity: newYorkCity,
			DepartureDate:   time.Now().UnixNano(),
			Guests:          []laresa.Guest{{Name: "Bob Sap"}, {Name: "Mike Tyson"}},
		}

		r := httptest.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/customers/%s/reservations", customer.Email),
			bytes.NewReader(bytesFrom(t, body)),
		)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		var resa laresa.Reservation
		err := json.NewDecoder(w.Body).Decode(&resa)
		if err != nil {
			t.Fatal(err)
		}

		r = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/customers/%s/reservations/%s", resa.Customer.Email, resa.ID),
			nil,
		)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want code: %d, got: %d", http.StatusOK, w.Code)
		}

		if !bytes.Equal(append(bytesFrom(t, resa), "\n"...), w.Body.Bytes()) {
			t.Errorf("want resa: %s, got: %s", string(bytesFrom(t, resa)), w.Body.String())
		}

	})

}
