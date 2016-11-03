package main

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func getBenchmarks(c *gin.Context) {
	benchmarks, err := ds.getAllBenchmarks()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, benchmarks)

}

func getClusters(c *gin.Context) {
	clusters, err := ds.getAllClusters()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, clusters)

}

func getMetrics(c *gin.Context) {
	metrics, err := ds.getAllMetrics()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, metrics)
}

func getRuns(c *gin.Context) {
	metric := c.Param("metric")
	build := c.Param("build")
	if metric == "" || build == "" {
		c.AbortWithError(400, errors.New("bad arguments"))
		return
	}
	runs, err := ds.getAllRuns(metric, build)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, runs)
}

func getTimelines(c *gin.Context) {
	timelines, err := ds.getAllTimelines()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, timelines)
}

func deleteBenchmark(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithError(400, errors.New("bad arguments"))
		return
	}
	err := ds.deleteBenchmark(id)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, nil)
}

func changeBenchmark(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithError(400, errors.New("bad arguments"))
		return
	}
	err := ds.reverseObsolete(id)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, nil)
}

func httpEngine() *gin.Engine {
	router := gin.Default()

	router.StaticFile("/", "./app/index.html")
	router.StaticFile("/admin", "./app/admin.html")

	router.Static("/css", "./app/css")
	router.Static("/js", "./app/js")
	router.Static("/partials", "./app/partials")

	v1 := router.Group("/api/v1")
	v1.GET("benchmarks", getBenchmarks)
	v1.GET("clusters", getClusters)
	v1.GET("metrics", getMetrics)
	v1.GET("timelines", getTimelines)

	v1.GET("runs/:metric/:build", getRuns)

	v1.DELETE("benchmarks/:id", deleteBenchmark)
	v1.PATCH("benchmarks/:id", changeBenchmark)

	return router
}
