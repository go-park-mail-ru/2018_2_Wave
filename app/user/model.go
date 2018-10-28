package user

import (
	"github.com/jmoiron/sqlx"
)

type Model struct {
	db *sqlx.DB
}

func NewModel(db *sqlx.DB) *Model {
	return &Model{
		db: db,
	}
}

// ----------------|
