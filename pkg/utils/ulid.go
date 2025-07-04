package utils

import (
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"time"
)

func (u *utils) GenerateULID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), ulid.Monotonic(rand.Reader, 0)).String()
}
