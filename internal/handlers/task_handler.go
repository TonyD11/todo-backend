package handlers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"todo-backend/internal/model"
)

var taskCollection *mongo.Collection

func SetCollection(col *mongo.Collection) {
	taskCollection = col
}

func AddTask(c *gin.Context) {

	var newTask model.Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validation
	if newTask.Title == "" {
		c.JSON(400, gin.H{"error": "title is required"})
		return
	}

	if newTask.Priority != "low" && newTask.Priority != "medium" && newTask.Priority != "high" {
		c.JSON(400, gin.H{"error": "invalid priority"})
		return
	}

	// Default values
	newTask.ID = primitive.NewObjectID()
	newTask.Done = false
	newTask.CreatedAt = time.Now()
	newTask.UpdatedAt = time.Now()

	result, err := taskCollection.InsertOne(context.Background(), newTask)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save"})
		return
	}

	c.JSON(201, gin.H{
		"id":          result.InsertedID,
		"title":       newTask.Title,
		"description": newTask.Description,
		"done":        newTask.Done,
		"priority":    newTask.Priority,
		"created_at":  newTask.CreatedAt,
		"updated_at":  newTask.UpdatedAt,
	})
}

func GetTasks(c *gin.Context) {

	var tasks []model.Task

	filter := bson.M{}

	// Query filters
	if done := c.Query("done"); done != "" {
		filter["done"] = done == "true"
	}

	if priority := c.Query("priority"); priority != "" {
		filter["priority"] = priority
	}

	cursor, err := taskCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch tasks"})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var task model.Task
		if err := cursor.Decode(&task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	c.JSON(200, tasks)
}

func UpdateTask(c *gin.Context) {

	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var updated model.Task

	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updateFields := bson.M{}

	if updated.Title != "" {
		updateFields["title"] = updated.Title
	}
	if updated.Description != "" {
		updateFields["description"] = updated.Description
	}
	if updated.Priority != "" {
		updateFields["priority"] = updated.Priority
	}
	updateFields["done"] = updated.Done
	updateFields["updated_at"] = time.Now()

	result, err := taskCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "task not found"})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}
func DeleteTask(c *gin.Context) {

	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	result, err := taskCollection.DeleteOne(
		context.Background(),
		bson.M{"_id": objID},
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "delete failed"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "task not found"})
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}
func GetTaskByID(c *gin.Context) {

	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var task model.Task

	err = taskCollection.FindOne(
		context.Background(),
		bson.M{"_id": objID},
	).Decode(&task)

	if err != nil {
		c.JSON(404, gin.H{"error": "task not found"})
		return
	}

	c.JSON(200, task)
}
