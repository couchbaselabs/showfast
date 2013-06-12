package main

import (
	"github.com/couchbaselabs/go-couchbase"
	"github.com/hoisie/mustache"
	"github.com/hoisie/web"
)

var c, _ = couchbase.Connect("http://127.0.0.1:8091/")
var pool, _ = c.GetPool("default")
var b_metrics, _ = pool.GetBucket("metrics")
var b_clusters, _ = pool.GetBucket("clusters")

func get_benchmarks() []map[string]interface{} {
	res, _ := b_metrics.View("metrics", "all", map[string]interface{}{
		"stale": false,
	})

	benchmarks := []map[string]interface{}{}
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

func home() string {
	var benchmarks = get_benchmarks()
	var content = ""

	for _, benchmark := range benchmarks {
		content += mustache.RenderFile("templates/benchmark.mustache", benchmark)
	}

	return mustache.RenderFile("templates/home.mustache", map[string]string{"content": content})
}

func main() {
	web.Get("/", home)
	web.Run("127.0.0.1:8080")
}
