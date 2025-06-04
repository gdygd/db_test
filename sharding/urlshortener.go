package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const virtualNodes = 10

var dbMap = map[string]*sql.DB{}
var ports = []string{"5433", "5434", "5435"}
var hashRing []uint32
var nodeMap = map[uint32]string{}

func main() {
	for _, port := range ports {
		dsn := "host=10.1.0.119 port=" + port + " user=postgres password=secret dbname=postgres sslmode=disable"
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("DB 연결 실패 (port: %s): %v", port, err)
		}
		dbMap[port] = db
	}

	initHashRing(ports)
	r := gin.Default()

	r.GET("/:urlId", func(c *gin.Context) {
		fmt.Printf("/:urlId...")
		urlId := c.Param("urlId")
		server := getPortFromHashRing(urlId)
		db := dbMap[server]

		fmt.Printf("%s, %v \n", urlId, server)

		var url string
		err := db.QueryRow("SELECT url FROM url_table WHERE url_id = $1", urlId).Scan(&url)
		if err != nil {
			fmt.Printf("err.. %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		// c.Redirect(http.StatusFound, url)
		c.JSON(http.StatusOK, gin.H{
			"urlId":  urlId,
			"url":    urlId,
			"server": server,
		})
	})

	r.POST("/", func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
			return
		}

		hash := sha256.Sum256([]byte(url))
		encoded := base64.StdEncoding.EncodeToString(hash[:])
		urlId := encoded[:5]
		server := getPortFromHashRing(urlId)
		db := dbMap[server]

		_, err := db.Exec("INSERT INTO url_table (url, url_id) VALUES ($1, $2)", url, urlId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"urlId":  urlId,
			"url":    url,
			"server": server,
		})
	})

	r.Run(":8081")
}

func hashKey(key string) uint32 {
	h := sha256.Sum256([]byte(key))
	return (uint32(h[0]) << 24) | (uint32(h[1]) << 16) | (uint32(h[2]) << 8) | uint32(h[3])
}

func initHashRing(ports []string) {
	for _, port := range ports {
		for i := 0; i < virtualNodes; i++ {
			vkey := port + "-" + strconv.Itoa(i)
			hash := hashKey(vkey)
			hashRing = append(hashRing, hash)
			nodeMap[hash] = port
		}
	}
	sort.Slice(hashRing, func(i, j int) bool {
		return hashRing[i] < hashRing[j]
	})
}

func getPortFromHashRing(key string) string {
	h := hashKey(key)
	for _, ringHash := range hashRing {
		if h <= ringHash {
			return nodeMap[ringHash]
		}
	}
	return nodeMap[hashRing[0]]
}
