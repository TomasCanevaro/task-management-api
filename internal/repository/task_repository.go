package repository

import "task-management-api/internal/models"

type TaskRepository interface {

	// User operations
	GetUser(id int) (*models.User, bool)

	// Task operations
	SaveTask(task *models.Task)
	GetTask(id int) (*models.Task, bool)
	GetAllTasks() []*models.Task

	// Notification operations
	SaveNotification(notification *models.Notification)
	GetNotificationsByUser(userID int) []*models.Notification

	// ID generators
	NextTaskID() int
	NextNotificationID() int
}
