package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	dbUser     = "root"
	dbPassword = "secret"
	dbName     = "simple_bank"
	dbHost     = "10.1.0.119"
	dbPort     = "5432"
)

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var db *sql.DB

func initDB() {
	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("DB 핑 실패: %v", err)
	}
	log.Println("DB 연결 성공")
}

func getUsername(c *gin.Context) {
	rows, err := db.Query("select username from users;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}
	c.JSON(http.StatusOK, items)
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/items", getUsername)

	r.Run(":8089")
}
