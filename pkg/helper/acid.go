package helper

import (
	"database/sql"
	"github.com/gogaruda/auth/pkg/apperror"
)

func WithTx(db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return apperror.New(apperror.CodeDBTxFailed, "gagal memulai transaksi", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				err = apperror.New(apperror.CodeDBTxFailed, "gagal commit transaksi", err)
			}
		}
	}()

	err = fn(tx)
	return
}
