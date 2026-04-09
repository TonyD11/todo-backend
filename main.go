package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	tasks := []string{
		"Leatn Go",
		"Build Backend",
		"Complete Internship Task",
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Todo Backend Running",
		})
	})

	r.GET("/tasks", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"tasks": tasks,
		})
	})

	r.POST("/tasks", func(c *gin.Context) {
		var newTask struct {
			Task string `json:"task"`
		}

		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		tasks = append(tasks, newTask.Task)
		c.JSON(201, gin.H{
			"message": "Task added",
			"tasks":   tasks,
		})
	})

	r.PUT("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedTask struct {
			Task string `json:"task"`
		}

		if err := c.ShouldBindJSON(&updatedTask); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Convert id to int
		index, err := strconv.Atoi(id)
		if err != nil || index < 0 || index >= len(tasks) {
			c.JSON(400, gin.H{"error": "Invalid task ID"})
			return
		}

		tasks[index] = updatedTask.Task
		c.JSON(200, gin.H{
			"message": "Task updated",
			"tasks":   tasks,
		})
	})
	r.Run(":8080")
}
