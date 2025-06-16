package utils

import (
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"time"
)

func NewULID() string {
	t := time.Now()
	return ulid.MustNew(ulid.Timestamp(t), ulid.Monotonic(rand.Reader, 0)).String()
}
