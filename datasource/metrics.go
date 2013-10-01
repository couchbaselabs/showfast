package datasource

import (
	"sort"
)

func (ds *DataSource) GetAllMetrics() (metrics []map[string]interface{}) {
	b_metrics := ds.GetBucket("metrics")
	b_clusters := ds.GetBucket("clusters")
	rows := ds.QueryView(b_metrics, "metrics", "all", map[string]interface{}{})

	for i := range rows {
		metric := rows[i].Value.(map[string]interface{})
		metric["id"] = rows[i].ID

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

func (ds *DataSource) getMetricIDsForBuild(build string) (ids []string) {
	b_benchmarks := ds.GetBucket("benchmarks")
	rows := ds.QueryView(b_benchmarks, "benchmarks", "metrics_by_build",
		map[string]interface{}{"key": build})

	for i := range rows {
		id := rows[i].Value.(string)
		ids = appendIfUnique(ids, id)
	}
	return
}

func (ds *DataSource) getMetricIDs(builds []string) (ids []string) {
	switch len(builds) {
	case 2:
		m1 := ds.getMetricIDsForBuild(builds[0])
		m2 := ds.getMetricIDsForBuild(builds[1])
		ids = getIntersection(m1, m2)
	case 1:
		ids = ds.getMetricIDsForBuild(builds[0])
	}
	return
}

func getIntersection(m1, m2 []string) (intersection []string) {
	sort.Strings(m1)
	sort.Strings(m2)
	for i, j := 0, 0; i < len(m1) && j < len(m2); {
		if m1[i] == m2[j] {
			intersection = append(intersection, m1[i])
			i++
		}
		j++
	}
	return
}

func (ds *DataSource) GetMetricsForBuilds(builds []string) (metrics []map[string]interface{}) {
	b_metrics := ds.GetBucket("metrics")
	b_clusters := ds.GetBucket("clusters")

	ids := ds.getMetricIDs(builds)
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
