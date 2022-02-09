package main

// A small project to mimic a simple service with cache and RDMS connection
// extending official GO guides.
// ┌────────────────────────────────────────────────────────────────────┐
// │                                                                    │
// │                                                                    │
// │                                                  ┌───────────┐     │
// │                                                  │           │     │
// │                         ┌─────────┬┬─────┐       │           │     │
// │                         │         ││     │       │           │     │
// │                         │         ││     │       │           │     │
// │    ┌──────────┐         │         ││     │       │           │     │
// │    │          │ Request │ HTTP    ││Redis│Cold   │   RDMS    │     │
// │    │ Client   ├────────►│ Service ││     ├─────► │           │     │
// │    │          │         │         ││     │Request│           │     │
// │    └──────────┘         │         ││     │       │           │     │
// │                         │         ││     │       │           │     │
// │                         └─────────┴┴─────┘       │           │     │
// │                                                  └───────────┘     │
// │                                                                    │
// │                                                                    │
// │                                                                    │
// │                                                                    │
// └────────────────────────────────────────────────────────────────────┘

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func getAlbums(c *gin.Context) {
	albums, err := allAlbums()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)

}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	id, err := addAlbum(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newAlbum.ID = id
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Fatalf("getAlbumById: %v", err)

	}

	// check cache
	alb, err := getFromCache(strconv.FormatInt(id, 10))
	if err == nil {
		c.Header("X-Cache", "HIT from memory")
		c.IndentedJSON(http.StatusOK, alb)
		return
	}
	c.Header("X-Cache", "MISS")

	alb, err = albumById(id)
	if err == nil {
		storeInCache(strconv.FormatInt(id, 10), alb)
		c.IndentedJSON(http.StatusOK, alb)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func main() {
	conDB()
	conRedis()
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.Run("localhost:8080")

	// albums, err := albumsByArtist("John Coltrane")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Albums found: %v,\n", albums)

	// alb, err := albumById(2)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Album found: ", alb)

	// albID, err := addAlbum(album{
	// 	Title:  "The Modern Sound of Betty Carter",
	// 	Artist: "Betty Carter",
	// 	Price:  49.99,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("ID of added album: %v\n", albID)

	// storeInCache("kk", album{
	// 	Title:  "The Modern Sound of Betty Carter",
	// 	Artist: "Betty Carter",
	// 	Price:  49.99,
	// })

}
