package apperror

import (
	"errors"
	"fmt"
)

type InitError struct {
	Code    string
	Message string
	Err     error
}

func (e *InitError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *InitError) Unwrap() error {
	return e.Err
}

func New(code, message string, err error) *InitError {
	return &InitError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is ngecek apakah error yang dikasih punya kode tertentu.
//
// Jadi gini... kadang kita dapet error dari berbagai lapisan (repo, service, dll),
// dan pengen tau: "Ini error karena apa sih? User gak ketemu? DB rusak?"
// Nah, fungsi ini bantuin kita buat ngecek apakah error itu punya kode tertentu.
//
// Cukup panggil aja: apperror.Is(err, apperror.CodeUserNotFound)
//
// Kalau error-nya dibungkus (wrapped) pake apperror.New, dia tetep bisa ke-detect kok.
//
// Param:
// - err: error yang mau dicek.
// - code: string kode error yang mau dicocokin.
//
// Return:
// - true: kalau error-nya memang punya kode itu.
// - false: kalau enggak.
//
// Contoh:
//
//	if apperror.Is(err, apperror.CodeDBError) {
//	    // mungkin log ke sentry, balikin 500, dll
//	}
func Is(err error, code string) bool {
	var e *InitError
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}
