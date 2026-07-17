package services

import (
	"testing"

	"task-management-api/internal/repository"
	"task-management-api/internal/storage"
)

func setupService() (*TaskService, repository.TaskRepository) {
	store := storage.NewMemoryStore()
	return NewTaskService(store), store
}

func TestSupervisorCanCreateTask(t *testing.T) {

	service, repo := setupService()

	supervisor, _ := repo.GetUser(1)

	task, err := service.CreateTask(
		"Write README",
		"Finish documentation",
		supervisor,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.Title != "Write README" {
		t.Fatalf("expected title 'Write README'")
	}

	if task.Status != "CREATED" {
		t.Fatalf("expected CREATED status")
	}
}

func TestWorkerCannotCreateTask(t *testing.T) {

	service, repo := setupService()

	worker, _ := repo.GetUser(2)

	_, err := service.CreateTask(
		"Illegal Task",
		"",
		worker,
	)

	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden")
	}
}

func TestAssignTask(t *testing.T) {

	service, repo := setupService()

	supervisor, _ := repo.GetUser(1)

	task, _ := service.CreateTask(
		"Task",
		"",
		supervisor,
	)

	err := service.AssignTask(
		task.ID,
		2,
		supervisor,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.AssignedTo == nil {
		t.Fatal("task should be assigned")
	}

	if task.Status != "ASSIGNED" {
		t.Fatal("expected ASSIGNED")
	}
}

func TestInvalidTransition(t *testing.T) {

	service, repo := setupService()

	supervisor, _ := repo.GetUser(1)
	worker, _ := repo.GetUser(2)

	task, _ := service.CreateTask(
		"Task",
		"",
		supervisor,
	)

	service.AssignTask(
		task.ID,
		2,
		supervisor,
	)

	err := service.UpdateStatus(
		task.ID,
		"COMPLETED",
		worker,
	)

	if err != ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus")
	}
}

func TestCompleteWorkflow(t *testing.T) {

	service, repo := setupService()

	supervisor, _ := repo.GetUser(1)
	worker, _ := repo.GetUser(2)

	task, _ := service.CreateTask(
		"Task",
		"",
		supervisor,
	)

	service.AssignTask(
		task.ID,
		2,
		supervisor,
	)

	err := service.UpdateStatus(
		task.ID,
		"IN_PROGRESS",
		worker,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = service.UpdateStatus(
		task.ID,
		"COMPLETED",
		worker,
	)

	if err != nil {
		t.Fatal(err)
	}

	if task.Status != "COMPLETED" {
		t.Fatalf("expected COMPLETED")
	}
}
