package utils

import (
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"time"
)

type ULIDs interface {
	Create() string
}

type ULIDCreate struct {
	entropy *ulid.MonotonicEntropy
}

func NewULIDCreate() *ULIDCreate {
	// entropy ini untuk mencegah ULID duplikat saat membuat banyak ULID dalam waktu yang sangat singkat
	entropy := ulid.Monotonic(rand.Reader, 0)
	return &ULIDCreate{entropy: entropy}
}

func (u *ULIDCreate) Create() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), u.entropy).String()
}
