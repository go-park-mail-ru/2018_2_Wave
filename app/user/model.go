package user

import (
	"Wave/utiles/walhalla"

	"github.com/jmoiron/sqlx"
)

type Model struct {
	db 		*sqlx.DB
}

func NewModel(ctx *walhalla.Context) *Model {
	return &Model{
		db: ctx.DB,
	}
}

// ----------------|
