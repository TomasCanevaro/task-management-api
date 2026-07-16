package storage

import (
	"sync"

	"task-management-api/internal/models"
)

type MemoryStore struct {
	mu sync.RWMutex

	Users         map[int]*models.User
	Tasks         map[int]*models.Task
	Notifications map[int]*models.Notification

	nextTaskID         int
	nextNotificationID int
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		Users:         make(map[int]*models.User),
		Tasks:         make(map[int]*models.Task),
		Notifications: make(map[int]*models.Notification),

		nextTaskID:         1,
		nextNotificationID: 1,
	}

	store.seedUsers()

	return store
}

func (m *MemoryStore) seedUsers() {
	m.Users[1] = &models.User{
		ID:    1,
		Name:  "Supervisor",
		Role:  models.RoleSupervisor,
		Email: "supervisor@example.com",
	}

	m.Users[2] = &models.User{
		ID:    2,
		Name:  "Worker One",
		Role:  models.RoleWorker,
		Email: "worker1@example.com",
	}

	m.Users[3] = &models.User{
		ID:    3,
		Name:  "Worker Two",
		Role:  models.RoleWorker,
		Email: "worker2@example.com",
	}
}

func (m *MemoryStore) NextTaskID() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.nextTaskID
	m.nextTaskID++

	return id
}

func (m *MemoryStore) NextNotificationID() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.nextNotificationID
	m.nextNotificationID++

	return id
}

func (m *MemoryStore) GetUser(id int) (*models.User, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.Users[id]
	return user, ok
}

func (m *MemoryStore) SaveTask(task *models.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Tasks[task.ID] = task
}

func (m *MemoryStore) GetTask(id int) (*models.Task, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.Tasks[id]
	return task, ok
}

func (m *MemoryStore) GetAllTasks() []*models.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(m.Tasks))

	for _, task := range m.Tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (m *MemoryStore) SaveNotification(notification *models.Notification) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Notifications[notification.ID] = notification
}

func (m *MemoryStore) GetNotificationsByUser(userID int) []*models.Notification {
	m.mu.RLock()
	defer m.mu.RUnlock()

	notifications := make([]*models.Notification, 0)

	for _, notification := range m.Notifications {
		if notification.UserID == userID {
			notifications = append(notifications, notification)
		}
	}

	return notifications
}
