package handlers

import (
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
func getNextID(tasks []models.Task) int {
	if len(tasks) == 0{
		return 1
	} else{
		return tasks[len(tasks)-1].ID + 1
	}
}

func getUniqueImagePath(baseDir, filename string) string {
	name := filename
	fullPath := filepath.Join(baseDir, name)

	for {
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			break
		}
		name = "ouz" + name
		fullPath = filepath.Join(baseDir, name)
	}

	return name
}

func TaskAdd(c *gin.Context) {
	// Chargement des tâches actuelles
	tasks, err := utils.LoadTasks()
	if err != nil {
		c.String(http.StatusInternalServerError, "Erreur chargement JSON")
		return
	}

	// Traitement de l'image
	imagePath := ""
	file, err := c.FormFile("imageFile")
	if err == nil {
		if file.Size > 1_000_000 {
			c.String(http.StatusBadRequest, "Image trop lourde (max 1MB)")
			return
		}

		src, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "Erreur lecture image")
			return
		}
		defer src.Close()

		img, _, err := image.Decode(src)
		if err != nil {
			c.String(http.StatusBadRequest, "Format d’image non supporté")
			return
		}

		// Création du nom de fichier compressé
		originalName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
		safeName := getUniqueImagePath("static/img", originalName+".jpg")
		dstPath := filepath.Join("static/img", safeName)

		out, err := os.Create(dstPath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Erreur création fichier image")
			return
		}
		defer out.Close()

		// Compression JPEG qualité 60
		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 60}); err != nil {
			c.String(http.StatusInternalServerError, "Erreur compression image")
			return
		}

		imagePath = "/" + strings.ReplaceAll(dstPath, "\\", "/")

	}

	progress, _ := strconv.Atoi(c.PostForm("progress"))

	newTask := models.Task{
		ID:          getNextID(tasks),
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Date:        time.Now().Format("2006-01-02"),
		PlannedDate: c.PostForm("plannedDate"),
		Progress:    progress,
		Image:       imagePath,
		Tags:        parseTags(c.PostFormArray("tags[]")),
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
		"Economy":  "bg-sky-900 text-white",
		"Dev":      "bg-teal-700 text-white",
		"Study":    "bg-yellow-400 text-black",
		"Research": "bg-red-700 text-white",
		"Others":   "bg-gray-700 text-white",
	}
	for _, n := range names {
		if c, ok := colors[n]; ok {
			tags = append(tags, models.Tag{Name: n, Color: c})
		}
	}
	return tags
}
