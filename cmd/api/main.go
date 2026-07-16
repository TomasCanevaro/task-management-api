package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"task-management-api/internal/handlers"
	"task-management-api/internal/middleware"
	"task-management-api/internal/services"
	"task-management-api/internal/storage"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API is running")
}

func healthAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {

	store := storage.NewMemoryStore()
	taskService := services.NewTaskService(store)
	taskHandler := handlers.NewTaskHandler(taskService)

	router := mux.NewRouter()

	router.HandleFunc("/", healthHandler).Methods("GET")
	router.HandleFunc("/health", healthAPIHandler).Methods("GET")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.UserMiddleware(store))

	//Protected routes

	protected.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	protected.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	protected.HandleFunc("/tasks/{id}/assign", taskHandler.AssignTask).Methods("POST")
	protected.HandleFunc("/tasks/{id}/status", taskHandler.UpdateStatus).Methods("PATCH")
	protected.HandleFunc("/notifications", taskHandler.GetNotifications).Methods("GET")

	fmt.Println("Server started on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}
