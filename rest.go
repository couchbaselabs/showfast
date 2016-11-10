package main

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func addBenchmark(c *gin.Context) {
	var b Benchmark
	if err := c.BindJSON(&b); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}
	err := ds.addBenchmark(b)
	if err != nil {
		c.AbortWithError(500, err)
	}
}

func getBenchmarks(c *gin.Context) {
	benchmarks, err := ds.getBenchmarks("by_metric_and_build")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, benchmarks)

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

func addCluster(c *gin.Context) {
	var cluster Cluster
	if err := c.BindJSON(&cluster); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}
	err := ds.addCluster(cluster)
	if err != nil {
		c.AbortWithError(500, err)
	}
}

func getClusters(c *gin.Context) {
	clusters, err := ds.getAllClusters()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, clusters)

}

func addMetric(c *gin.Context) {
	var m Metric
	if err := c.BindJSON(&m); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}
	err := ds.addMetric(m)
	if err != nil {
		c.AbortWithError(500, err)
	}
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

func httpEngine() *gin.Engine {
	router := gin.Default()

	router.StaticFile("/", "./app/index.html")
	router.StaticFile("/admin", "./app/admin.html")

	router.Static("/css", "./app/css")
	router.Static("/js", "./app/js")
	router.Static("/partials", "./app/partials")

	v1 := router.Group("/api/v1")

	v1.POST("benchmarks", addBenchmark)
	v1.GET("benchmarks", getBenchmarks)
	v1.PATCH("benchmarks/:id", changeBenchmark)
	v1.DELETE("benchmarks/:id", deleteBenchmark)

	v1.POST("clusters", addCluster)
	v1.GET("clusters", getClusters)

	v1.POST("metrics", addMetric)
	v1.GET("metrics", getMetrics)

	v1.GET("timelines", getTimelines)
	v1.GET("runs/:metric/:build", getRuns)

	return router
}
