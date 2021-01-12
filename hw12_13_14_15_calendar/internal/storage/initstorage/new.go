package initstorage

import (
	"context"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage/memorystorage"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage/sqlstorage"
)

func New(ctx context.Context, inmem bool, connect string) (storage.Storage, error) {
	var db storage.Storage
	if inmem {
		db = memorystorage.New()
	} else {
		db = sqlstorage.New()
	}
	err := db.Connect(ctx, connect)
	return db, err
}
