package datasource

import (
	"encoding/json"
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

func GetTimeline(metric string) []byte {
	b_benchmarks := GetBucket("benchmarks")
	rows := QueryView(b_benchmarks, "benchmarks", "build_and_value_by_metric",
		map[string]interface{}{"key": metric})

	timeline := Timeline{}
	for i := range rows {
		xy := rows[i].Value.([]interface{})
		timeline = append(timeline, xy)
	}
	sort.Sort(timeline)
	t, _ := json.Marshal(timeline)
	return t
}

func getLatestBuild(metric string) string {
	b_benchmarks := GetBucket("benchmarks")
	rows := QueryView(b_benchmarks, "benchmarks", "build_by_metric",
		map[string]interface{}{"key": metric})

	builds := []string{}
	for i := range rows {
		builds = append(builds, rows[i].Value.(string))
	}
	sort.Strings(builds)
	return builds[len(builds)-1]
}

func GetTimelineForBuilds(metric string, builds_s string) []byte {
	builds := strings.Split(builds_s, "/")
	if len(builds) == 1 {
		builds = append(builds, getLatestBuild(metric))
	}
	keys := [][]string{}
	for i := range builds {
		key := []string{metric, builds[i]}
		keys = append(keys, key)
	}

	b_benchmarks := GetBucket("benchmarks")
	rows := QueryView(b_benchmarks, "benchmarks", "values_by_build_and_metric",
		map[string]interface{}{"keys": keys})

	timeline := Timeline{}
	for i := range rows {
		xy := []interface{}{rows[i].Key.([]interface{})[1],
			rows[i].Value.(interface{})}
		timeline = append(timeline, xy)
	}
	sort.Sort(timeline)
	t, _ := json.Marshal(timeline)
	return t
}
