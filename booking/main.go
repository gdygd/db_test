package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Seat struct {
	ID       int    `json:"id"`
	IsBooked int    `json:"isbooked"`
	Name     string `json:"name"`
}

var db *sql.DB

func main() {
	dsn := "host=10.1.0.119 port=5432 user=root password=secret dbname=simple_bank sslmode=disable"
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()
	r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"http://10.1.0.119:8080"},
			AllowMethods: []string{
				http.MethodHead,
				http.MethodOptions,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete},
			AllowHeaders: []string{
				"Content-Type",
				"Authorization"},
		}))

	r.StaticFile("/", "index.html")

	r.GET("/seats", getSeats)
	r.PUT("/:id/:name", bookSeat)

	r.Run(":8080")
}

func getSeats(c *gin.Context) {
	rows, err := db.Query("SELECT id, isbooked, name FROM seats")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var seats []Seat
	for rows.Next() {
		var s Seat
		if err := rows.Scan(&s.ID, &s.IsBooked, &s.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		seats = append(seats, s)
	}
	c.JSON(http.StatusOK, seats)
}

func bookSeat(c *gin.Context) {
	fmt.Println("#0 bookSeat")
	idStr := c.Param("id")
	name := c.Param("name")
	fmt.Println("#0 bookSeat")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	fmt.Println("#1 bookSeat")
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	fmt.Println("#2 bookSeat")
	var booked int
	err = tx.QueryRow("SELECT isbooked FROM seats WHERE id = $1 FOR UPDATE", id).Scan(&booked)
	fmt.Println("#3 bookSeat")
	if err == sql.ErrNoRows || booked == 1 {
		c.JSON(http.StatusOK, gin.H{"error": "Seat already booked"})
		fmt.Println("#4 bookSeat")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("#5 bookSeat")
		return
	}

	fmt.Println("#6 bookSeat")
	_, err = tx.Exec("UPDATE seats SET isbooked = 1, name = $1 WHERE id = $2", name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "Booked successfully"})
}
