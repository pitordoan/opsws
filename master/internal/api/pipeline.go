package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"opsws/master/internal/model"
)

// Migrate creates the pipelines table if it doesn't exist.
func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS pipelines (
	id TEXT PRIMARY KEY,
	name TEXT,
	agent TEXT,
	labels TEXT,
	stages TEXT
);
`)
	return err
}

func RegisterPipelineRoutes(r *mux.Router, db *sql.DB) {
	s := r.PathPrefix("/api/pipelines").Subrouter()
	s.HandleFunc("", listPipelines(db)).Methods("GET")
	s.HandleFunc("", createPipeline(db)).Methods("POST")
	s.HandleFunc("/{id}", getPipeline(db)).Methods("GET")
	s.HandleFunc("/{id}", updatePipeline(db)).Methods("PUT")
	s.HandleFunc("/{id}", deletePipeline(db)).Methods("DELETE")
}

func listPipelines(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, agent, labels, stages FROM pipelines")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var result []*model.Pipeline
		for rows.Next() {
			var p model.Pipeline
			var agentJSON, labelsJSON, stagesJSON string
			if err := rows.Scan(&p.ID, &p.Name, &agentJSON, &labelsJSON, &stagesJSON); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.Unmarshal([]byte(agentJSON), &p.Agent)
			json.Unmarshal([]byte(labelsJSON), &p.Labels)
			json.Unmarshal([]byte(stagesJSON), &p.Stages)
			result = append(result, &p)
		}
		json.NewEncoder(w).Encode(result)
	}
}

func createPipeline(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p model.Pipeline
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		agentJSON, _ := json.Marshal(p.Agent)
		labelsJSON, _ := json.Marshal(p.Labels)
		stagesJSON, _ := json.Marshal(p.Stages)
		_, err := db.Exec(
			"INSERT INTO pipelines(id, name, agent, labels, stages) VALUES (?, ?, ?, ?, ?)",
			p.ID, p.Name, string(agentJSON), string(labelsJSON), string(stagesJSON),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func getPipeline(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var p model.Pipeline
		var agentJSON, labelsJSON, stagesJSON string
		err := db.QueryRow("SELECT id, name, agent, labels, stages FROM pipelines WHERE id = ?", id).
			Scan(&p.ID, &p.Name, &agentJSON, &labelsJSON, &stagesJSON)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.Unmarshal([]byte(agentJSON), &p.Agent)
		json.Unmarshal([]byte(labelsJSON), &p.Labels)
		json.Unmarshal([]byte(stagesJSON), &p.Stages)
		json.NewEncoder(w).Encode(p)
	}
}

func updatePipeline(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var p model.Pipeline
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		agentJSON, _ := json.Marshal(p.Agent)
		labelsJSON, _ := json.Marshal(p.Labels)
		stagesJSON, _ := json.Marshal(p.Stages)
		res, err := db.Exec(
			"UPDATE pipelines SET name = ?, agent = ?, labels = ?, stages = ? WHERE id = ?",
			p.Name, string(agentJSON), string(labelsJSON), string(stagesJSON), id,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(p)
	}
}

func deletePipeline(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		res, err := db.Exec("DELETE FROM pipelines WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}