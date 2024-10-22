package models

import (
	"time"
)

type Task struct {
	ID          uint      // Primary key
	Task       string    // Title of the task
	Description string    // Detailed description of the task
	Completed   bool      // Indicates if the task is completed
	CreatedAt   time.Time // Timestamp for creation
	UpdatedAt   time.Time // Timestamp for updates
}
