package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

// album represents data about a record album.
type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Hello World!"})
}

func getAlbums(c *gin.Context) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		albums = append(albums, album)
		if err := rows.Err(); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum Album

	err := c.BindJSON(&newAlbum)

	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err})
		return
	}

	err = db.QueryRow(
		"INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		newAlbum.Title,
		newAlbum.Artist,
		newAlbum.Price,
	).Scan(&newAlbum.ID)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err})
		return
	}

	var createdAlbum Album
	createdAlbum.ID = newAlbum.ID
	createdAlbum.Artist = newAlbum.Artist
	createdAlbum.Price = newAlbum.Price
	createdAlbum.Title = newAlbum.Title
	c.IndentedJSON(http.StatusCreated, &createdAlbum)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var album Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
			return
		}
	}
	c.IndentedJSON(http.StatusOK, &album)
}

func updateAlbumByID(c *gin.Context) {
	var updateAlbum Album

	id := c.Param("id")

	if err := c.BindJSON(&updateAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bind error"})
		return
	}

	_, err := db.Exec("UPDATE album SET title = $1, artist = $2, price = $3 WHERE id = $4", updateAlbum.Title, updateAlbum.Artist, updateAlbum.Price, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "update error"})
		return
	}

	c.IndentedJSON(http.StatusNoContent, &updateAlbum)
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var album Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
			return
		}
	}

	if _, err := db.Exec("DELETE FROM album WHERE id = $1", id); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err})
		return
	}

	c.IndentedJSON(http.StatusNoContent, &album)
}

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		getEnv("DB_SSLMODE", "disable"),
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/", home)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.PATCH("/albums/:id", updateAlbumByID)
	router.DELETE("/albums/:id", deleteAlbumByID)

	router.Run()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
