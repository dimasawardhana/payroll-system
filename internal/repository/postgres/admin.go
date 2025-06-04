package postgres

import (
	"context"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	pool *pgxpool.Pool
}

func NewAdminRepository(pool *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{
		pool: pool,
	}
}

func (r *AdminRepository) GetAdmin(ctx context.Context, credential domain.Admin) (domain.Admin, error) {
	var admin domain.Admin
	if credential.Email == "" {
		return domain.Admin{}, error_const.ErrInvalidCredentials
	}
	err := r.pool.
		QueryRow(ctx, "SELECT id, email, password_hash, role FROM admins WHERE email = $1", credential.Email).
		Scan(&admin.ID, &admin.Email, &admin.Password_hash, &admin.Role)

	if err != nil {
		return domain.Admin{}, err
	}
	if admin.ID == 0 {
		return domain.Admin{}, error_const.ErrUserNotFound
	}
	return admin, nil
}
