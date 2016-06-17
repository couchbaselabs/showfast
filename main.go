package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/hoisie/web"
)

const (
	address   = ":8000"
	staticDir = "app"
)

var (
	dataSource     DataSource
	cbHost, cbPass string
)

func home() []byte {
	content, _ := ioutil.ReadFile("app/index.html")
	return content
}

func admin() []byte {
	content, _ := ioutil.ReadFile("app/admin.html")
	return content
}

func release() []byte {
	content, _ := ioutil.ReadFile("app/release.html")
	return content
}

func feed() []byte {
	content, _ := ioutil.ReadFile("app/feed.html")
	return content
}

func allRuns(ctx *web.Context) []byte {
	metric := ctx.Params["metric"]
	build := ctx.Params["build"]

	return dataSource.getAllRuns(metric, build)
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

func init() {
	cbHost = os.Getenv("CB_HOST")
	cbPass = os.Getenv("CB_PASS")
	if cbHost == "" || cbPass == "" {
		log.Fatalln("Missing Couchbase Server settings.")
	}
}

func main() {
	dataSource = DataSource{cbHost, cbPass}

	web.Get("/", home)
	web.Get("/admin", admin)
	web.Get("/release", release)
	web.Get("/feed", feed)

	web.Config.StaticDir = staticDir

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

	web.Run(address)
}
