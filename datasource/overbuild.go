package datasource

import (
	"log"
)

func GetBenchmarks() (benchmarks []map[string]interface{}) {
	pool := GetPool()
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
	return
}
