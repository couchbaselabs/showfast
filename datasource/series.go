package datasource

import (
	"encoding/json"
	"log"
	"sort"
	"strings"
)

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

func GetTimeline(metric string) (t []byte) {
	pool := GetPool()
	b_benchmarks, err := pool.GetBucket("benchmarks")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}

	res, err := b_benchmarks.View("benchmarks", "build_and_value_by_metric",
		map[string]interface{}{"stale": false, "key": metric})
	if err != nil {
		InstallDDoc("benchmarks")
	}

	timeline := Timeline{}
	for i := range res.Rows {
		xy := res.Rows[i].Value.([]interface{})
		timeline = append(timeline, xy)
	}
	sort.Sort(timeline)
	t, _ = json.Marshal(timeline)
	return
}

func getLatestBuild(metric string) string {
	pool := GetPool()
	b_benchmarks, err := pool.GetBucket("benchmarks")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}
	res, err := b_benchmarks.View("benchmarks", "build_by_metric",
		map[string]interface{}{"stale": false, "key": metric})
	if err != nil {
		InstallDDoc("benchmarks")
	}
	builds := []string{}
	for i := range res.Rows {
		builds = append(builds, res.Rows[i].Value.(string))
	}
	sort.Strings(builds)
	return builds[len(builds)-1]
}

func GetTimelineForBuilds(metric string, builds_s string) (t []byte) {
	pool := GetPool()
	b_benchmarks, err := pool.GetBucket("benchmarks")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}

	var builds = strings.Split(builds_s, "/")
	if len(builds) == 1 {
		builds = append(builds, getLatestBuild(metric))
	}
	keys := [][]string{}
	for i := range builds {
		key := []string{metric, builds[i]}
		keys = append(keys, key)
	}

	res, err := b_benchmarks.View("benchmarks", "values_by_build_and_metric",
		map[string]interface{}{"stale": false, "keys": keys})
	if err != nil {
		InstallDDoc("benchmarks")
	}

	timeline := Timeline{}
	for i := range res.Rows {
		xy := []interface{}{res.Rows[i].Key.([]interface{})[1],
			res.Rows[i].Value.(interface{})}
		timeline = append(timeline, xy)
	}
	sort.Sort(timeline)
	t, _ = json.Marshal(timeline)
	return
}
