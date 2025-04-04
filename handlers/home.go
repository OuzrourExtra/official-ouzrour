package handlers

import (
	"net/http"
	"os"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ouzrourextra/official-ouzrour/models"
)


func loadTasks() []models.Task {
	data, err := os.ReadFile("data/tasks.json")
	if err != nil {
		return []models.Task{}
	}

	var tasks []models.Task
	_ = json.Unmarshal(data, &tasks)
	return tasks
}

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home/index", gin.H{
		"Title": "Tâches Essentielles",
		"Tasks": loadTasks(),
	})
}
