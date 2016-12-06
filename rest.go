package main

import (
	"errors"
	"github.com/gin-gonic/gin"

	log "gopkg.in/inconshreveable/log15.v2"
)

func addBenchmark(c *gin.Context) {
	var b Benchmark
	if err := c.BindJSON(&b); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		log.Error("error adding benchmark", "err", err)
		return
	}
	err := ds.addBenchmark(b)
	if err != nil {
		c.AbortWithError(500, err)
	}
}

func getBenchmarks(c *gin.Context) {
	component := c.Param("component")
	category := c.Param("category")

	benchmarks, err := ds.getBenchmarks(component, category)
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
	err := ds.reverseHidden(id)
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
		log.Error("error adding cluster", "err", err)
		return
	}
	err := ds.addCluster(cluster)
	if err != nil {
		c.AbortWithError(500, err)
	}
}

func addMetric(c *gin.Context) {
	var m Metric
	if err := c.BindJSON(&m); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		log.Error("error adding metric", "err", err)
		return
	}
	err := ds.addMetric(m)
	if err != nil {
		c.AbortWithError(500, err)
	}
}

func getMetrics(c *gin.Context) {
	component := c.Param("component")
	category := c.Param("category")

	metrics, err := ds.getMetrics(component, category)
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

func getTimeline(c *gin.Context) {
	metric := c.Param("metric")

	timeline, err := ds.getTimeline(metric)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, timeline)
}

func httpEngine() *gin.Engine {
	router := gin.Default()

	router.StaticFile("/", "./app/index.html")
	router.StaticFile("/admin", "./app/admin.html")
	router.Static("/static", "./app")

	rg := router.Group("/api/v1")

	rg.POST("benchmarks", addBenchmark)
	rg.GET("benchmarks/:component/:category", getBenchmarks)
	rg.PATCH("benchmarks/:id", changeBenchmark)
	rg.DELETE("benchmarks/:id", deleteBenchmark)

	rg.POST("clusters", addCluster)

	rg.POST("metrics", addMetric)
	rg.GET("metrics/:component/:category", getMetrics)

	rg.GET("timeline/:metric", getTimeline)

	rg.GET("runs/:metric/:build", getRuns)

	return router
}
