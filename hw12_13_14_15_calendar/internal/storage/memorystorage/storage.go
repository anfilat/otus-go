package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type data map[int]storage.Event

type store struct {
	mu     sync.Mutex
	lastID int
	data   data
}

func (s *store) Connect(_ context.Context, _ string) error {
	return nil
}

func (s *store) Close(_ context.Context) error {
	return nil
}

func (s *store) Create(_ context.Context, event storage.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	event.ID = id
	s.data[id] = storage.Event{
		ID:           id,
		Title:        event.Title,
		Start:        event.Start,
		Stop:         event.Stop,
		Description:  event.Description,
		UserID:       event.UserID,
		Notification: event.Notification,
	}
	return id, nil
}

func (s *store) Update(_ context.Context, id int, change storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.data[id]
	if !ok {
		return storage.ErrNotExistsEvent
	}

	event.Title = change.Title
	event.Start = change.Start
	event.Stop = change.Stop
	event.Description = change.Description
	event.Notification = change.Notification
	s.data[id] = event

	return nil
}

func (s *store) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}

func (s *store) DeleteAll(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(data)
	return nil
}

func (s *store) ListAll(_ context.Context) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]storage.Event, 0, len(s.data))
	for _, event := range s.data {
		result = append(result, event)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result, nil
}

func (s *store) ListDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, month, day := date.Date()
	for _, event := range s.data {
		eventYear, eventMonth, eventDay := event.Start.Date()
		if eventYear == year && eventMonth == month && eventDay == day {
			result = append(result, event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result, nil
}

func (s *store) ListWeek(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, week := date.ISOWeek()
	for _, event := range s.data {
		eventYear, eventWeek := event.Start.ISOWeek()
		if eventYear == year && eventWeek == week {
			result = append(result, event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result, nil
}

func (s *store) ListMonth(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, month, _ := date.Date()
	for _, event := range s.data {
		eventYear, eventMonth, _ := event.Start.Date()
		if eventYear == year && eventMonth == month {
			result = append(result, event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result, nil
}

func (s *store) IsTimeBusy(_ context.Context, userID int, start, stop time.Time, excludeID int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.data {
		if event.UserID == userID && event.ID != excludeID && event.Start.Before(stop) && event.Stop.After(start) {
			return true, nil
		}
	}
	return false, nil
}

func (s *store) newID() int {
	s.lastID++
	return s.lastID
}
