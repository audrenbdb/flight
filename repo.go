package laresa

import "context"

type inMemRepo struct {
	reservations map[string]Reservation
}

type inMemRepoOption func(repo *inMemRepo)

func WithReservations(reservations map[string]Reservation) inMemRepoOption {
	return func(repo *inMemRepo) {
		repo.reservations = reservations
	}
}

func NewInMemRepo(options ...inMemRepoOption) *inMemRepo {
	repo := &inMemRepo{
		reservations: map[string]Reservation{},
	}
	for _, opt := range options {
		opt(repo)
	}
	return repo
}

func (r *inMemRepo) SaveReservation(ctx context.Context, resa Reservation) error {
	r.reservations[resa.ID] = resa
	return nil
}

func (r *inMemRepo) GetReservation(ctx context.Context, id string) (Reservation, error) {
	resa, ok := r.reservations[id]
	if !ok {
		return Reservation{}, ErrReservationNotFound
	}
	return resa, nil
}
