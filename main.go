package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/hoisie/web"
)

var pkgDir string

var dataSource DataSource

func home() []byte {
	content, _ := ioutil.ReadFile(pkgDir + "app/index.html")
	return content
}

func allRuns(ctx *web.Context) []byte {
	metric := ctx.Params["metric"]
	build := ctx.Params["build"]

	return dataSource.getAllRuns(metric, build)
}

func admin() []byte {
	content, _ := ioutil.ReadFile(pkgDir + "app/admin.html")
	return content
}

func release() []byte {
	content, _ := ioutil.ReadFile(pkgDir + "app/release.html")
	return content
}

func feed() []byte {
	content, _ := ioutil.ReadFile(pkgDir + "app/feed.html")
	return content
}

func getComparison(ctx *web.Context) []byte {
	baseline := ctx.Params["baseline"]
	target := ctx.Params["target"]
	return dataSource.getComparison(baseline, target)
}

func deleteBenchmark(ctx *web.Context) {
	var params struct {
		ID string `json:"id"`
	}
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(body, &params)
	dataSource.deleteBenchmark(params.ID)
}

func reverseObsolete(ctx *web.Context) {
	var params struct {
		ID string `json:"id"`
	}
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(body, &params)
	dataSource.reverseObsolete(params.ID)
}

type Config struct {
	CouchbaseAddress, BucketPassword, ListenAddress string
}

func main() {
	pkgDir = os.Getenv("GOPATH") + "/src/github.com/couchbaselabs/showfast/"
	web.Config.StaticDir = pkgDir + "app"

	configFile, err := ioutil.ReadFile(pkgDir + "config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	dataSource = DataSource{config.CouchbaseAddress, config.BucketPassword}

	web.Get("/", home)
	web.Get("/admin", admin)
	web.Get("/release", release)
	web.Get("/feed", feed)

	web.Get("/all_metrics", dataSource.getAllMetrics)
	web.Get("/all_clusters", dataSource.getAllClusters)
	web.Get("/all_timelines", dataSource.getAllTimelines)
	web.Get("/all_benchmarks", dataSource.getAllBenchmarks)
	web.Get("/all_runs", allRuns)
	web.Get("/all_releases", dataSource.getAllReleases)
	web.Get("/all_feed_records", dataSource.getAllFeedRecords)
	web.Get("/get_comparison", getComparison)
	web.Post("/delete", deleteBenchmark)
	web.Post("/reverse_obsolete", reverseObsolete)

	web.Run(config.ListenAddress)
}
