package entdb

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/go-faster/errors"
	"github.com/gotd/td/session"

	"github.com/cydev/cgbot/internal/ent"
)

type SessionStorage struct {
	id int
	db *ent.Client
}

func (s *SessionStorage) LoadSession(ctx context.Context) ([]byte, error) {
	ts, err := s.db.TelegramSession.Get(ctx, s.id)
	if ent.IsNotFound(err) {
		return nil, session.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	return ts.Data, nil
}

func (s *SessionStorage) StoreSession(ctx context.Context, data []byte) error {
	if err := s.db.TelegramSession.
		Create().
		SetID(s.id).
		SetData(data).
		OnConflict(
			sql.ConflictColumns("id"),
			sql.ResolveWith(func(set *sql.UpdateSet) {
				set.Set("data", data)
			}),
		).Exec(ctx); err != nil {
		return errors.Wrap(err, "create")
	}
	return nil
}

func NewSessionStorage(id int, db *ent.Client) *SessionStorage {
	return &SessionStorage{id: id, db: db}
}
