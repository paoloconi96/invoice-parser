package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func __main() {
	router := gin.Default()
	router.GET("/api/v1/albums", getAlbums)
	router.GET("/api/v1/albums/:id", getAlbumById)
	router.POST("/api/v1/albums", addAlbum)

	router.Run("localhost:8000")
}

func getAlbums(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(ctx *gin.Context) {
	var newAlbum album

	if err := ctx.BindJSON(&newAlbum); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
	}

	albums = append(albums, newAlbum)
}

func getAlbumById(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, album := range albums {
		if album.ID == id {
			ctx.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	ctx.IndentedJSON(http.StatusNotFound, "Album not found.")
}
