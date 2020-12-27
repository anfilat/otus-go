package memorystorage

import "github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"

func New() storage.Storage {
	result := store{}
	result.data = make(data)
	return &result
}
