package httpserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
)

func handleUpdate(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req := Event{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		change := httpEventToStorageEvent(req)
		err = app.Update(r.Context(), req.ID, change)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, OkResult{Ok: true})
	}
}
