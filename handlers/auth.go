package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(c *gin.Context) {
	errorParam := c.Query("error")
	c.HTML(http.StatusOK, "login/login", gin.H{
		"Error": errorParam != "",
	})
}

func Login(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Impossible de charger le fichier .env")
	}
	username := c.PostForm("username")
	password := c.PostForm("password")
	expectedUser := os.Getenv("ADMIN_USER")
	expectedHash := "$2a$10$mEYU5MCUAiPV9EG/JaM9AO1PfEQ4IPLMNsVF2y8auZFZfvm9mVbP."
	// bcrypt hash
	fmt.Println("✔ Hash détecté :", expectedHash)
	fmt.Printf("✔ Hash détecté : %v", expectedHash)
	fmt.Println("✔ Taille :", len(expectedHash))
	
	fmt.Println("username = ",username , "exceptec = ",expectedUser , "test = ",username == expectedUser)
	fmt.Println("pass = " , password , "expected = ",expectedHash, "Valid = " , bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(password)))
	if username != expectedUser || bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(password)) != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=1")
		return
	}

	session := sessions.Default(c)
	session.Set("authenticated", true)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/tasks/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		auth := session.Get("authenticated")
		if auth != true {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

