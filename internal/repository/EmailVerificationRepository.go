package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/model"
)

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *model.EmailVerificationModel) error
	FindByToken(ctx context.Context, token string) (*model.EmailVerificationModel, error)
	MarkAsUsed(ctx context.Context, id string) error
}

type emailVerificationRepository struct {
	db *sql.DB
}

func NewEmailVerificationRepository(db *sql.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(ctx context.Context, ev *model.EmailVerificationModel) error {
	query := `INSERT INTO email_verifications (id, user_id, token, expires_at) VALUES(?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, ev.ID, ev.UserID, ev.Token, ev.ExpiresAt)
	if err != nil {
		return apperror.New(apperror.CodeDBError, "query insert email_verifications gagal", err)
	}
	return nil
}

func (r *emailVerificationRepository) FindByToken(ctx context.Context, token string) (*model.EmailVerificationModel, error) {
	query := `SELECT id, user_id, expires_at, is_used FROM email_verifications WHERE token = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, token)

	var ev model.EmailVerificationModel
	err := row.Scan(&ev.ID, &ev.UserID, &ev.ExpiresAt, &ev.IsUsed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.New("[TOKEN_NOT_FOUND]", "token tidak ditemukan", err, 404)
		}
		return nil, apperror.New(apperror.CodeDBError, "query select email_verifications gagal", err)
	}

	return &ev, nil
}

func (r *emailVerificationRepository) MarkAsUsed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE email_verifications SET is_used = true WHERE id = ?`, id)
	if err != nil {
		return apperror.New(apperror.CodeDBError, "query update is_used gagal", err)
	}
	return nil
}
