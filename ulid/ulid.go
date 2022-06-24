package ulid

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func NewBuilder() func() string {
	t := time.Now().UTC()
	unixNano := t.UnixNano()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(unixNano)), 0)
	return func() string {
		id := ulid.MustNew(ulid.Timestamp(t), entropy)
		return id.String()
	}
}
