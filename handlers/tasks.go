package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ouzrourextra/official-ouzrour/models"
	"github.com/ouzrourextra/official-ouzrour/utils"
)

func TasksIndex(c *gin.Context) {
	tasks, err := utils.LoadTasks()
	if err != nil {
		c.String(http.StatusInternalServerError, "Erreur chargement des tâches")
		return
	}
	c.HTML(http.StatusOK, "tasks/tasks", tasks)
}

func TaskAdd(c *gin.Context) {
	progress, _ := strconv.Atoi(c.PostForm("progress"))

	newTask := models.Task{
		ID:          int(time.Now().Unix()), // ID unique basé sur timestamp
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Date:        time.Now().Format("2006-01-02"),
		PlannedDate: c.PostForm("plannedDate"),
		Progress:    progress,
		Image:       c.PostForm("image"),
		Tags:        parseTags(c.PostFormArray("tags[]")),
	}

	tasks, err := utils.LoadTasks()
	if err != nil {
		c.String(http.StatusInternalServerError, "Erreur chargement JSON")
		return
	}

	tasks = append(tasks, newTask)
	if err := utils.SaveTasks(tasks); err != nil {
		c.String(http.StatusInternalServerError, "Erreur sauvegarde JSON")
		return
	}

	c.Redirect(http.StatusSeeOther, "/tasks/")
}

func TaskUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "ID invalide")
		return
	}

	tasks, err := utils.LoadTasks()
	if err != nil {
		c.String(http.StatusInternalServerError, "Erreur chargement JSON")
		return
	}

	task, i, err := utils.FindTaskByID(tasks, id)
	if err != nil {
		c.String(http.StatusNotFound, "Tâche introuvable")
		return
	}

	progress, _ := strconv.Atoi(c.PostForm("progress"))
	task.Title = c.PostForm("title")
	task.Description = c.PostForm("description")
	task.PlannedDate = c.PostForm("plannedDate")
	task.Image = c.PostForm("image")
	task.Progress = progress
	task.Tags = parseTags(c.PostFormArray("tags[]"))
	tasks[i] = *task

	if err := utils.SaveTasks(tasks); err != nil {
		c.String(http.StatusInternalServerError, "Erreur mise à jour")
		return
	}

	c.Redirect(http.StatusSeeOther, "/tasks/")
}

func TaskDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "ID invalide")
		return
	}

	tasks, err := utils.LoadTasks()
	if err != nil {
		c.String(http.StatusInternalServerError, "Erreur lecture JSON")
		return
	}

	var newTasks []models.Task
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}

	if err := utils.SaveTasks(newTasks); err != nil {
		c.String(http.StatusInternalServerError, "Erreur suppression")
		return
	}

	c.Redirect(http.StatusSeeOther, "/tasks/")
}

func parseTags(names []string) []models.Tag {
	var tags []models.Tag
	colors := map[string]string{
		"Economy": "bg-sky-900",
		"Dev":      "bg-teal-700 text-white",
		"Study":    "bg-yellow-400 text-black",
		"Research":   "bg-red-700 text-white",
		"Others"  : "bg-gray-700 text-white",
	}
	for _, n := range names {
		if c, ok := colors[n]; ok {
			tags = append(tags, models.Tag{Name: n, Color: c})
		}
	}
	return tags
}
