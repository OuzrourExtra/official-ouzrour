package utils

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ouzrourextra/official-ouzrour/models"
)

const filePath = "data/tasks.json"

func LoadTasks() ([]models.Task, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func SaveTasks(tasks []models.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func FindTaskByID(tasks []models.Task, id int) (*models.Task, int, error) {
	for i, task := range tasks {
		if task.ID == id {
			return &tasks[i], i, nil
		}
	}
	return nil, -1, errors.New("task not found")
}
