package main

import (
	"todo-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	connectDB()

	handlers.SetCollection(taskCollection)

	r := gin.Default()

	r.GET("/tasks", handlers.GetTasks)
	r.POST("/tasks", handlers.AddTask)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.GET("/tasks/:id", handlers.GetTaskByID)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.Run(":8080")
}
