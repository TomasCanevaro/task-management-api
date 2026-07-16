package models

type Role string

const (
	RoleSupervisor Role = "SUPERVISOR"
	RoleWorker     Role = "WORKER"
)
