package middleware

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YosaZiege/golang-react-todo/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {

	var err error
    connStr := "user=yosa dbname=todo-bd sslmode=disable password=yosa"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Could not ping the database:", err)
    }
    fmt.Println("Connected to the database")
}
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id_task, task, description, completed FROM task ORDER BY createdat DESC") // Ensure the query is correct
    if err != nil {
        log.Println("Error fetching tasks:", err) // Log the error
        http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        err := rows.Scan(&task.ID, &task.Task, &task.Description, &task.Completed) // Ensure the scan matches your model
        if err != nil {
            log.Println("Error scanning tasks:", err) // Log the error
            http.Error(w, "Error scanning tasks", http.StatusInternalServerError)
            return
        }
        tasks = append(tasks, task)
    }

    // Return tasks as JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(tasks); err != nil {
        log.Println("Error encoding tasks to JSON:", err) // Log the error
        http.Error(w, "Error encoding tasks", http.StatusInternalServerError)
    }
}


func CreateTask(w http.ResponseWriter, r *http.Request) {
    var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	// Assign the current time to created_at and updated_at fields
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Insert the task into the database
	query := `INSERT INTO task (task, description, completed, createdat, updatedat) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id_task`
	err = db.QueryRow(query, task.Task, task.Description, task.Completed, task.CreatedAt, task.UpdatedAt).Scan(&task.ID)
	if err != nil {
		log.Println("Error inserting task:", err)
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	// Return the created task as a response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Println("Error encoding task to JSON:", err)
		http.Error(w, "Error encoding task", http.StatusInternalServerError)
	}
}
func TaskComplete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskId := vars["id"] // the route being task/{id}
    query :=  `UPDATE task SET completed = true WHERE id_task = $1`

    _ , err := db.Exec(query , taskId)
    if err != nil {
        log.Println("Error updating task status:", err)
        http.Error(w , "Error updating task status", http.StatusInternalServerError)
        return
    }

    response := map[string]string{
        "message": "Task marked as complete successfully",
    }

    w.Header().Set("Content-type" , "application/json")
    json.NewEncoder(w).Encode(response)
}
func UndoTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskId := vars["id"] // the route being task/{id}

    // Update query to mark the task as incomplete
    query := `UPDATE task SET completed = false WHERE id_task = $1`

    // Execute the update statement
    _, err := db.Exec(query, taskId)
    if err != nil {
        log.Println("Error updating task status:", err)
        http.Error(w, "Error updating task status", http.StatusInternalServerError)
        return
    }

    // Prepare the response message
    response := map[string]string{
        "message": "Task marked as incomplete successfully", // Updated message
    }

    // Set response header and encode the response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // Set the response status to 200 OK
    json.NewEncoder(w).Encode(response)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskId := vars["id"] // the route being task/{id}

    // Update query to mark the task as incomplete
    query := `DELETE FROM task WHERE id_task = $1`

    // Execute the update statement
    _, err := db.Exec(query, taskId)
    if err != nil {
        log.Println("Error deleting task :", err)
        http.Error(w, "Error deleting task ", http.StatusInternalServerError)
        return
    }

    // Prepare the response message
    response := map[string]string{
        "message":  "Task deleted successfully", // Updated message
    }

    // Set response header and encode the response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // Set the response status to 200 OK
    json.NewEncoder(w).Encode(response)
}
func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {

}
