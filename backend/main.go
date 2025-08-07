package main

import (
	"fmt"
	"database/sql"
	"log"
	"os"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

var db *sql.DB

func main() {
	// Read DB config from env
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Could not open DB connection:  ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("❌ Could not ping DB: ", err)
	}

	fmt.Println("Successfully connected to the database!")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Allow all origins
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Go Backend!"})
	})

	r.GET("/api/messages", getMessages)

	r.Run(":8080") // Listen on port 8080
}

func getMessages(c *gin.Context) {
	if db == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection is nil"})
        return
    }

    rows, err := db.Query("SELECT id, message FROM test_table")
    if err != nil {
        log.Printf("❌ Error querying DB: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query DB"})
        return
    }
    defer rows.Close()

    var messages []map[string]interface{}

    for rows.Next() {
        var id int
        var message string
        if err := rows.Scan(&id, &message); err != nil {
            continue
        }
        messages = append(messages, gin.H{"id": id, "message": message})
    }

    c.JSON(http.StatusOK, messages)
}