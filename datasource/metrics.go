package datasource

import (
	"log"
	"sort"
)

func GetAllMetrics() (metrics []map[string]interface{}) {
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
		metric := res.Rows[i].Value.(map[string]interface{})
		metric["id"] = res.Rows[i].ID

		rv := map[string]string{}
		b_clusters.Get(metric["cluster"].(string), &rv)
		for k, v := range rv {
			metric[k] = v
		}
		metrics = append(metrics, metric)
	}
	return
}

func appendIfUnique(slice []string, s string) []string {
	for i := range slice {
		if slice[i] == s {
			return slice
		}
	}
	return append(slice, s)
}

func getMetricIDsForBuild(build string) (ids []string) {
	pool := GetPool()
	b_benchmarks, err := pool.GetBucket("benchmarks")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}

	res, err := b_benchmarks.View("benchmarks", "by_build", map[string]interface{}{
		"stale": false,
		"key":   build,
	})
	if err != nil {
		log.Fatalf("Error reading view:  %v", err)
	}

	for i := range res.Rows {
		id := res.Rows[i].Value.(string)
		ids = appendIfUnique(ids, id)
	}
	return
}

func getMetricIDs(builds []string) (ids []string) {
	switch len(builds) {
	case 2:
		m1 := getMetricIDsForBuild(builds[0])
		m2 := getMetricIDsForBuild(builds[1])
		ids = getIntersection(m1, m2)
	case 1:
		ids = getMetricIDsForBuild(builds[0])
	}
	return
}

func getIntersection(m1, m2 []string) (intersection []string) {
	sort.Strings(m1)
	sort.Strings(m2)
	for i, j := 0, 0; i < len(m1) && j < len(m2); i++ {
		if m1[i] == m2[j] {
			intersection = append(intersection, m1[i])
			j++
		}
	}
	return
}

func GetMetricsForBuilds(builds []string) (metrics []map[string]interface{}) {
	pool := GetPool()
	b_metrics, err := pool.GetBucket("metrics")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}
	b_clusters, err := pool.GetBucket("clusters")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}

	ids := getMetricIDs(builds)
	for i := range ids {
		metric := map[string]interface{}{
			"id":     ids[i],
			"builds": builds,
		}
		b_metrics.Get(ids[i], &metric)

		rv := map[string]string{}
		b_clusters.Get(metric["cluster"].(string), &rv)
		for k, v := range rv {
			metric[k] = v
		}
		metrics = append(metrics, metric)
	}
	return
}
