package main

import (
	"encoding/json"
	"flag"
	"github.com/couchbaselabs/go-couchbase"
	"github.com/hoisie/mustache"
	"github.com/hoisie/web"
)

var c, _ = couchbase.Connect("http://127.0.0.1:8091/")
var pool, _ = c.GetPool("default")
var b_metrics, _ = pool.GetBucket("metrics")
var b_clusters, _ = pool.GetBucket("clusters")
var b_benchmarks, _ = pool.GetBucket("benchmarks")

func get_benchmarks() (benchmarks []map[string]interface{}) {
	res, _ := b_metrics.View("metrics", "all", map[string]interface{}{
		"stale": false,
	})

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
	res, _ := b_benchmarks.View("test", "by_metric", map[string]interface{}{
		"stale": false,
		"key":   ctx.Params["metric"],
	})

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
		content += mustache.RenderFile("templates/benchmark.mustache", benchmark)
	}
	return mustache.RenderFile("templates/home.mustache", map[string]string{"content": content})
}

func main() {
	address := flag.String("address", "127.0.0.1:8080", "Listen address")
	flag.Parse()

	web.Get("/", home)
	web.Get("/timeline", timeline)
	web.Run(*address)
}
