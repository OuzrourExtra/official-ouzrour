package main

import (
    "encoding/json"
    "net/http"
    "os"
	"fmt"

    "github.com/gin-gonic/gin"
    "github.com/ouzrourextra/official-ouzrour/config"
)

type Tag struct {
	Name  string `json:"Name"`
	Color string `json:"Color"`
}

type Task struct {
	ID          int     `json:"id"`
	Title       string  `json:"Title"`
	Description string  `json:"Description"`
	Date        string  `json:"Date"`
	PlannedDate string  `json:"PlannedDate"` // ✅ ce champ doit exister et commencer par majuscule
	Progress    int     `json:"Progress"`
	Image       string  `json:"Image"`
	Tags        []Tag   `json:"Tags"`
}


func loadTasks() []Task {
    file, err := os.ReadFile("data/tasks.json")
    if err != nil {
        return []Task{}
    }

    var tasks []Task
    _ = json.Unmarshal(file, &tasks)
    return tasks
}

func main() {
    r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

    r.Static("/static", "./static")
    r.HTMLRender = config.LoadTemplates()
	
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "home/index", gin.H{
            "Title": "Tâches Essentielles",
            "Tasks": loadTasks(),
        })
    })

	
	fmt.Printf("🧪 Nb tâches chargées : %d\n", len(loadTasks()))


    r.Run("localhost:8080")
}
