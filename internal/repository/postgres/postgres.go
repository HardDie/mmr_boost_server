package postgres

import (
	"github.com/HardDie/mmr_boost_server/internal/db"
)

type Postgres struct {
	txManager *txManager

	accessToken
	history
	password
	user
}

func NewPostgres(db *db.DB) *Postgres {
	return &Postgres{
		txManager: newTxManager(db),

		accessToken: newAccessToken(db),
		history:     newHistory(db),
		password:    newPassword(db),
		user:        newUser(db),
	}
}

func (r *Postgres) TxManager() *txManager {
	return r.txManager
}
