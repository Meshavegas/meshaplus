package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Task represents a single task in the system.
type Task struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Completed bool        `json:"completed"`
	CreatedAt time.Time   `json:"createdAt"`
	Reminders []time.Time `json:"reminders"`
	IsExotic  bool        `json:"isExotic"`
}

// FinancialRecord represents a single financial record.
type FinancialRecord struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	IsRevenue   bool      `json:"isRevenue"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
}

// Goal represents a user's life goal.
type Goal struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tasks       []Task   `json:"tasks"`
	Milestones  []string `json:"milestones"`
	Progress    float64  `json:"progress"`
}

// AIPersona represents the AI's management style.
type AIPersona struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// --- In-Memory Data Stores ---
var tasks = make(map[string]Task)
var financialRecords = make(map[string]FinancialRecord)
var goals = make(map[string]Goal)
var personas = []AIPersona{
	{ID: "1", Name: "The Nurturing Mentor", Description: "Gentle, encouraging, and supportive."},
	{ID: "2", Name: "The Strict Coach", Description: "Firm, disciplined, and focused on results."},
	{ID: "3", Name: "The Zen Guide", Description: "Mindful, calm, and focused on balance."},
}

func main() {
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/financial-records", financialRecordsHandler)
	http.HandleFunc("/goals", goalsHandler)
	http.HandleFunc("/personas", personasHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(tasks)
	case http.MethodPost:
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task.ID = "task-" + string(time.Now().UnixNano())
		task.CreatedAt = time.Now()
		tasks[task.ID] = task
		json.NewEncoder(w).Encode(task)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func financialRecordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(financialRecords)
	case http.MethodPost:
		var record FinancialRecord
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		record.ID = "record-" + string(time.Now().UnixNano())
		record.Date = time.Now()
		financialRecords[record.ID] = record
		json.NewEncoder(w).Encode(record)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func goalsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(goals)
	case http.MethodPost:
		var goal Goal
		if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		goal.ID = "goal-" + string(time.Now().UnixNano())
		goals[goal.ID] = goal
		json.NewEncoder(w).Encode(goal)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func personasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(personas)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
