package laresa

import "errors"

var (
	ErrTooManyGuests       = errors.New("can't have more than three guests")
	ErrReservationNotFound = errors.New("reservation not found")
)
