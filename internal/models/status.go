package models

type TaskStatus string

const (
	StatusCreated    TaskStatus = "CREATED"
	StatusAssigned   TaskStatus = "ASSIGNED"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"
)
