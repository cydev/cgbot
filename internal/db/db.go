package entdb

import (
	"github.com/cydev/cgbot/internal/ent"
)

type DB struct {
	ent *ent.Client
}

func New(ent *ent.Client) *DB {
	return &DB{
		ent: ent,
	}
}
