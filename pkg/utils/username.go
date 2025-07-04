package utils

import (
	"strings"
	"unicode"
)

func (u *utils) GenerateUsernameFromName(name string) string {
	// ubah ke lowercase dan ganti spasi menjadi titik
	username := strings.ToLower(name)
	username = strings.ReplaceAll(username, " ", ".")

	// Menghapus karakter non-alfabet dan non-pemakaian biasa
	username = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
			return r
		}
		return -1
	}, username)

	return username
}
