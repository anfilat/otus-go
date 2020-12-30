package sqlstorage

import "github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"

func New() storage.Storage {
	return &store{}
}
