package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func createUser() {

}

func selectUserData(c *gin.Context) {
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

func updateUser() {

}

func deleteUser() {

}
