package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

var tasks = []string{
	"Learn Go",
	"Build Backend",
	"Complete Internship Task",
}

func GetTasks(c *gin.Context) {
	c.JSON(200, gin.H{
		"tasks": tasks,
	})
}
func AddTask(c *gin.Context) {
	var newTask struct {
		Task string `json:"task"`
	}

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	tasks = append(tasks, newTask.Task)

	c.JSON(201, gin.H{
		"message": "Task Added",
		"tasks":   tasks,
	})
}
func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var updatedTask struct {
		Task string `json:"task"`
	}

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	index, err := strconv.Atoi(id)
	if err != nil || index < 0 || index >= len(tasks) {
		c.JSON(400, gin.H{
			"error": "Invalid Task ID",
		})
		return
	}

	tasks[index] = updatedTask.Task

	c.JSON(200, gin.H{
		"message": "Task Updated",
		"tasks":   tasks,
	})
}
