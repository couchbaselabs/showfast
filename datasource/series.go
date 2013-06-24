package datasource

import (
	"encoding/json"
	"strings"

	"github.com/couchbaselabs/go-couchbase"
)

func getLatestBuild(metric string) string {
	b_benchmarks := GetBucket("benchmarks")
	rows := QueryView(b_benchmarks, "benchmarks", "values_by_build_and_metric",
		map[string]interface{}{
			"startkey": []string{metric, "z"}, "endkey": []string{metric},
			"descending": true, "limit": 1})
	return rows[0].Key.([]interface{})[1].(string)
}

func rowsToTimeline(rows []couchbase.ViewRow) []byte {
	timeline := [][]interface{}{}
	for i := range rows {
		build := rows[i].Key.([]interface{})[1]
		value := rows[i].Value.(interface{})
		timeline = append(timeline, []interface{}{build, value})
	}
	t, _ := json.Marshal(timeline)
	return t
}

func GetTimeline(metric string) []byte {
	b_benchmarks := GetBucket("benchmarks")
	rows := QueryView(b_benchmarks, "benchmarks", "values_by_build_and_metric",
		map[string]interface{}{
			"startkey": []string{metric}, "endkey": []string{metric, "z"}})
	return rowsToTimeline(rows)
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
	return rowsToTimeline(rows)
}
