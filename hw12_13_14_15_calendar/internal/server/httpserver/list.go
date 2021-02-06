package httpserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
)

func handleListDay(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleList(w, r, app.ListDay)
	}
}

func handleListWeek(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleList(w, r, app.ListWeek)
	}
}

func handleListMonth(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleList(w, r, app.ListMonth)
	}
}

func handleList(w http.ResponseWriter, r *http.Request, fn app.ListEvents) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := ListRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := fn(r.Context(), req.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := make(ListResult, 0, len(events))
	for _, event := range events {
		result = append(result, storageEventToHTTPEvent(event))
	}
	writeJSON(w, result)
}
