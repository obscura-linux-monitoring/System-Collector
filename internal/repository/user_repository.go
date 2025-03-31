package repository

import (
	"database/sql"
	"system-collector/pkg/logger"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	sugar := logger.GetSugar()
	sugar.Info("UserRepository 초기화 중")

	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) ExistsUserByObscuraKey(ObscuraKey string) (bool, error) {
	sugar := logger.GetSugar()
	sugar.Infow("사용자 존재 여부 확인 시작", "ObscuraKey", ObscuraKey)

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE obscura_key = $1)`
	row := r.db.QueryRow(query, ObscuraKey)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		sugar.Errorw("사용자 존재 여부 확인 오류", "error", err)
		return false, err
	}

	sugar.Infow("사용자 존재 여부 확인 완료", "ObscuraKey", ObscuraKey, "exists", exists)
	return exists, nil
}
