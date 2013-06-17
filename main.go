package main

import (
	"encoding/json"
	"flag"
	"github.com/couchbaselabs/go-couchbase"
	"github.com/hoisie/mustache"
	"github.com/hoisie/web"
	"log"
	"os"
	"sort"
)

var pool couchbase.Pool
var pckg_dir string

func get_benchmarks() (benchmarks []map[string]interface{}) {
	b_metrics, err := pool.GetBucket("metrics")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}
	b_clusters, err := pool.GetBucket("clusters")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}

	res, err := b_metrics.View("metrics", "all", map[string]interface{}{
		"stale": false,
	})
	if err != nil {
		log.Fatalf("Error reading view:  %v", err)
	}

	for i := range res.Rows {
		benchmark := res.Rows[i].Value.(map[string]interface{})
		benchmark["id"] = res.Rows[i].ID

		rv := map[string]string{}
		b_clusters.Get(benchmark["cluster"].(string), &rv)
		for k, v := range rv {
			benchmark[k] = v
		}
		benchmarks = append(benchmarks, benchmark)
	}
	return benchmarks
}

type Timeline [][]interface{}

func (b Timeline) Len() int {
	return len(b)
}

func (b Timeline) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Timeline) Less(i, j int) bool {
	return b[i][0].(string) < b[j][0].(string)
}

func timeline(ctx *web.Context) []byte {
	b_benchmarks, err := pool.GetBucket("benchmarks")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}

	res, err := b_benchmarks.View("benchmarks", "by_metric", map[string]interface{}{
		"stale": false,
		"key":   ctx.Params["metric"],
	})
	if err != nil {
		log.Fatalf("Error reading view:  %v", err)
	}

	timeline := Timeline{}
	for i := range res.Rows {
		xy := res.Rows[i].Value.([]interface{})
		timeline = append(timeline, xy)
	}
	sort.Sort(timeline)
	t, _ := json.Marshal(timeline)
	return t
}

func home() string {
	content := ""
	for _, benchmark := range get_benchmarks() {
		content += mustache.RenderFile(pckg_dir+"templates/benchmark.mustache", benchmark)
	}
	return mustache.RenderFile(pckg_dir+"templates/home.mustache", map[string]string{"content": content})
}

func main() {
	address := flag.String("address", "127.0.0.1:8080", "Listen address")
	flag.Parse()

	c, err := couchbase.Connect("http://127.0.0.1:8091/")
	if err != nil {
		log.Fatalf("Error connecting:  %v", err)
	}
	pool, _ = c.GetPool("default")

	pckg_dir = os.Getenv("GOPATH") + "/src/github.com/pavel-paulau/showfast/"
	web.Config.StaticDir = pckg_dir + "static"
	web.Get("/", home)
	web.Get("/timeline", timeline)
	web.Run(*address)
}
