package main

import (
	"bitbucket.org/kardianos/osext"
	"encoding/json"
	"flag"
	"github.com/couchbaselabs/go-couchbase"
	"github.com/hoisie/mustache"
	"github.com/hoisie/web"
	"log"
	"path/filepath"
	"strings"
)

var pool couchbase.Pool
var exec_dir string

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

	for _, row := range res.Rows {
		benchmark := row.Value.(map[string]interface{})
		benchmark["id"] = row.ID

		rv := map[string]string{}
		b_clusters.Get(benchmark["cluster"].(string), &rv)
		for k, v := range rv {
			benchmark[k] = v
		}
		benchmarks = append(benchmarks, benchmark)
	}
	return benchmarks
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

	values := [][]interface{}{}
	for _, row := range res.Rows {
		xy := row.Value.([]interface{})
		values = append(values, xy)
	}
	b, _ := json.Marshal(values)
	return b
}

func home() string {
	var benchmarks = get_benchmarks()
	var content = ""
	for _, benchmark := range benchmarks {
		content += mustache.RenderFile(exec_dir+"templates/benchmark.mustache", benchmark)
	}
	return mustache.RenderFile(exec_dir+"templates/home.mustache", map[string]string{"content": content})
}

func main() {
	address := flag.String("address", "127.0.0.1:8080", "Listen address")
	flag.Parse()

	c, err := couchbase.Connect("http://127.0.0.1:8091/")
	if err != nil {
		log.Fatalf("Error connecting:  %v", err)
	}
	pool, _ = c.GetPool("default")

	filename, _ := osext.Executable()
	if strings.Contains(filename, "go-build") {
		exec_dir = ""
	} else {
		exec_dir = filepath.Dir(filename)
	}

	web.Get("/", home)
	web.Get("/timeline", timeline)
	web.Run(*address)
}
