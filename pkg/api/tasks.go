package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

const tasksLimit = 50

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	tasks, err := db.Tasks(search, tasksLimit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, tasksResp{Tasks: tasks})
}
