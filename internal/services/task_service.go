package services

import (
	"errors"
	"fmt"
	"time"

	"task-management-api/internal/models"
	"task-management-api/internal/repository"
)

type TaskService struct {
	store repository.TaskRepository
}

func NewTaskService(store repository.TaskRepository) *TaskService {
	return &TaskService{
		store: store,
	}
}

var (
	ErrForbidden       = errors.New("forbidden")
	ErrTaskNotFound    = errors.New("task not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidStatus   = errors.New("invalid status transition")
	ErrAlreadyAssigned = errors.New("task already assigned")
)

func (s *TaskService) CreateTask(title, description string, supervisor *models.User) (*models.Task, error) {

	if supervisor.Role != models.RoleSupervisor {
		return nil, ErrForbidden
	}

	task := &models.Task{
		ID:          s.store.NextTaskID(),
		Title:       title,
		Description: description,
		Status:      models.StatusCreated,
		CreatedBy:   supervisor.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.store.SaveTask(task)

	return task, nil
}

func (s *TaskService) AssignTask(taskID int, workerID int, supervisor *models.User) error {

	if supervisor.Role != models.RoleSupervisor {
		return ErrForbidden
	}

	task, ok := s.store.GetTask(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	if task.Status != models.StatusCreated {
		return ErrAlreadyAssigned
	}

	worker, ok := s.store.GetUser(workerID)
	if !ok {
		return ErrUserNotFound
	}

	if worker.Role != models.RoleWorker {
		return ErrUserNotFound
	}

	task.AssignedTo = &workerID
	task.Status = models.StatusAssigned
	task.UpdatedAt = time.Now()

	s.store.SaveTask(task)

	message := fmt.Sprintf("Task %d assigned to worker %d", task.ID, workerID)

	fmt.Println(message)

	notification := &models.Notification{
		ID:        s.store.NextNotificationID(),
		UserID:    workerID,
		TaskID:    task.ID,
		Message:   message,
		CreatedAt: time.Now(),
	}

	s.store.SaveNotification(notification)

	return nil
}

func (s *TaskService) UpdateStatus(taskID int, status models.TaskStatus, worker *models.User) error {

	task, ok := s.store.GetTask(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	if worker.Role != models.RoleWorker {
		return ErrForbidden
	}

	if task.AssignedTo == nil {
		return ErrForbidden
	}

	if *task.AssignedTo != worker.ID {
		return ErrForbidden
	}

	switch task.Status {

	case models.StatusAssigned:

		if status != models.StatusInProgress {
			return ErrInvalidStatus
		}

	case models.StatusInProgress:

		if status != models.StatusCompleted {
			return ErrInvalidStatus
		}

	default:
		return ErrInvalidStatus
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	s.store.SaveTask(task)

	return nil
}

func (s *TaskService) ListTasks(user *models.User) []*models.Task {

	tasks := s.store.GetAllTasks()

	if user.Role == models.RoleSupervisor {
		return tasks
	}

	var assigned []*models.Task

	for _, task := range tasks {

		if task.AssignedTo != nil && *task.AssignedTo == user.ID {
			assigned = append(assigned, task)
		}
	}

	return assigned
}

func (s *TaskService) GetTaskByID(taskID int, user *models.User) (*models.Task, error) {
	task, ok := s.store.GetTask(taskID)
	if !ok {
		return nil, ErrTaskNotFound
	}

	// Supervisors can view any task
	if user.Role == models.RoleSupervisor {
		return task, nil
	}

	// Workers can only view their assigned tasks
	if task.AssignedTo != nil && *task.AssignedTo == user.ID {
		return task, nil
	}

	return nil, ErrForbidden
}

func (s *TaskService) GetNotifications(user *models.User) []*models.Notification {
	return s.store.GetNotificationsByUser(user.ID)
}
