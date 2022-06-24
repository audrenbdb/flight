package laresa

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type BookService interface {
	BookFlight(ctx context.Context, customer Customer, params BookingParams) (Reservation, error)
	GetReservation(ctx context.Context, id string) (Reservation, error)
}

func NewChiHTTPHandler(service BookService) http.Handler {
	router := chi.NewRouter()
	router.Post("/customers/{email}/reservations", handlePostReservation(service))
	router.Get("/customers/{email}/reservations/{id}", handleGetReservation(service))
	return router
}

func handleGetReservation(service BookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resa, err := service.GetReservation(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resa)
	}
}

func handlePostReservation(service BookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params BookingParams
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !params.IsComplete() {
			http.Error(w, "missing booking parameters", http.StatusBadRequest)
			return
		}
		resa, err := service.BookFlight(r.Context(), Customer{Email: chi.URLParam(r, "email")}, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resa)
	}
}
