package main

import (
	"duck/avito/db"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// using gin to establish the different endpoint of the API
	router := gin.Default()
	router.GET("/user/:id", getUser)
	router.GET("/newSeg/:slug", newSeg)
	router.GET("/newSegAutoAdd/:slug/:percent", newSegAutoAdd)
	router.GET("/rmSeg/:slug", rmSeg)
	router.GET("/add/:slug/:id", addToSeg)
	router.GET("/addTTL/:slug/:id/:ttl", addToSegTTL)
	router.GET("/remove/:slug/:id", rmFromSeg)

	router.Run("localhost:8081")
}

// getUser Handler to get a user from the API
func getUser(c *gin.Context) {
	// query parameter id which is the user id
	user_id, _ := strconv.Atoi(c.Param("id"))
	//connect to the db
	db.ConnectDB()
	// get the user by the id
	user := db.GetUser(user_id)
	// return the json
	c.IndentedJSON(http.StatusOK, user)
}

// Handler to create a new segment
func newSeg(c *gin.Context) {
	slug := c.Param("slug")
	db.ConnectDB()
	status := db.CreateSegment(slug)
	if status == 1 {
		c.IndentedJSON(http.StatusOK, "{added:1}")
	} else {
		c.IndentedJSON(http.StatusOK, "{added:0}")
	}
}

// handler to create a new segment and add random users to it
func newSegAutoAdd(c *gin.Context) {
	slug := c.Param("slug")
	percent, _ := strconv.Atoi(c.Param("percent"))
	db.ConnectDB()
	status := db.CreateSegment(slug)
	segmentId := db.GetSegmentBySlug(slug)
	log.Println(segmentId)
	db.AddRandomUserPercent(percent, segmentId)
	if status == 1 {
		c.IndentedJSON(http.StatusOK, "{added:1}")
	} else {
		c.IndentedJSON(http.StatusOK, "{added:0}")
	}
}

// handler to remove a segment
func rmSeg(c *gin.Context) {
	slug := c.Param("slug")
	db.ConnectDB()
	status := db.DeleteSegment(slug)
	if status == 1 {
		c.IndentedJSON(http.StatusOK, "{deleted:1}")
	} else {
		c.IndentedJSON(http.StatusOK, "{deleted:0}")
	}
}

// handler to add a user to a segment
func addToSeg(c *gin.Context) {
	slug := c.Param("slug")
	user_id, _ := strconv.Atoi(c.Param("id"))
	slugs := strings.Split(slug, ",")
	db.ConnectDB()

	for _, value := range slugs {
		segId := db.GetSegmentBySlug(value)
		db.AddUserToSegment(user_id, segId)
	}
	c.IndentedJSON(http.StatusOK, "{done:1}")

}

// handler to add a user to a segment
func addToSegTTL(c *gin.Context) {
	slug := c.Param("slug")
	user_id, _ := strconv.Atoi(c.Param("id"))
	ttl := c.Param("ttl")
	slugs := strings.Split(slug, ",")
	db.ConnectDB()

	for _, value := range slugs {
		segId := db.GetSegmentBySlug(value)
		db.AddUserToSegmentTTL(user_id, segId, ttl)
	}
	c.IndentedJSON(http.StatusOK, "{done:1}")

}

// handler to remove a user from a segment
func rmFromSeg(c *gin.Context) {
	slug := c.Param("slug")
	user_id, _ := strconv.Atoi(c.Param("id"))
	slugs := strings.Split(slug, ",")
	db.ConnectDB()

	for _, value := range slugs {
		segId := db.GetSegmentBySlug(value)
		db.RemoveUserFromSegment(user_id, segId)
	}
	c.IndentedJSON(http.StatusOK, "{done:1}")

}
