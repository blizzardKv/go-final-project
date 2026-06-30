package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}
