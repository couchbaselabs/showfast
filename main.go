package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/hoisie/web"
)

var pckg_dir string

var data_source DataSource

func home() []byte {
	content, _ := ioutil.ReadFile(pckg_dir + "app/index.html")
	return content
}

func all_runs(ctx *web.Context) []byte {
	metric := ctx.Params["metric"]
	build := ctx.Params["build"]

	return data_source.GetAllRuns(metric, build)
}

func admin() []byte {
	content, _ := ioutil.ReadFile(pckg_dir + "app/admin.html")
	return content
}

func delete(ctx *web.Context) {
	id := ctx.Params["id"]
	data_source.DeleteBenchmark(id)
}

type Config struct {
	CouchbaseAddress, BucketPassword, ListenAddress string
}

func main() {
	pckg_dir = os.Getenv("GOPATH") + "/src/github.com/couchbaselabs/showfast/"
	web.Config.StaticDir = pckg_dir + "app"

	config_file, err := ioutil.ReadFile(pckg_dir + "config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal(config_file, &config)
	if err != nil {
		log.Fatal(err)
	}

	data_source = DataSource{config.CouchbaseAddress, config.BucketPassword}

	web.Get("/", home)
	web.Get("/admin", admin)

	web.Get("/all_metrics", data_source.GetAllMetrics)
	web.Get("/all_clusters", data_source.GetAllClusters)
	web.Get("/all_timelines", data_source.GetAllTimelines)
	web.Get("/all_benchmarks", data_source.GetAllBenchmarks)
	web.Get("/all_runs", all_runs)
	web.Post("/delete", delete)

	web.Run(config.ListenAddress)
}
