package laresa_test

import (
	"encoding/json"
	"laresa"
	"testing"
	"time"
)

var (
	parisCity              = laresa.City{Name: "Paris"}
	newYorkCity            = laresa.City{Name: "New-York"}
	parisToNewYorkDistance = 5836.5

	jon = laresa.Customer{Email: "jon@doe.org"}

	monday    = time.Date(2022, 06, 20, 0, 0, 0, 0, time.UTC).UnixNano()
	tuesday   = time.Date(2022, 06, 21, 0, 0, 0, 0, time.UTC).UnixNano()
	wednesday = time.Date(2022, 06, 22, 0, 0, 0, 0, time.UTC).UnixNano()
	thursday  = time.Date(2022, 06, 23, 0, 0, 0, 0, time.UTC).UnixNano()
)

func fakeIDGenerator(id string) func() string {
	return func() string {
		return id
	}
}

func bytesFrom(t *testing.T, i any) []byte {
	b, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
