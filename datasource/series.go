package datasource

import (
	"encoding/json"
	"strconv"
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

func GetAllBenchmarks() (benchmarks []map[string]string) {
	b_benchmarks := GetBucket("benchmarks")
	rows := QueryView(b_benchmarks, "benchmarks", "values_by_build_and_metric",
		map[string]interface{}{})
	for _, row := range rows {
		benchmark := map[string]string{
			"id":     row.ID,
			"metric": row.Key.([]interface{})[0].(string),
			"build":  row.Key.([]interface{})[1].(string),
			"value":  strconv.FormatFloat(row.Value.(float64), 'f', 1, 64),
		}
		benchmarks = append(benchmarks, benchmark)
	}
	return
}

func GetObsoleteBenchmarks(metric string, build string) (benchmarks []map[string]interface{}) {
	b_benchmarks := GetBucket("benchmarks")
	params := map[string]interface{}{
		"startkey": []string{metric, build},
		"endkey":   []string{metric, build},
	}
	rows := QueryView(b_benchmarks, "benchmarks", "value_and_reports_by_build_and_metric", params)
	for _, row := range rows {
		benchmark := map[string]interface{}{
			"value":   strconv.FormatFloat(row.Value.([]interface{})[0].(float64), 'f', 1, 64),
			"report1": row.Value.([]interface{})[1],
			"report2": row.Value.([]interface{})[2],
		}
		benchmarks = append(benchmarks, benchmark)
	}
	return
}

func DeleteBenchmark(benchmark string) {
	b_benchmarks := GetBucket("benchmarks")
	b_benchmarks.Delete(benchmark)
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
