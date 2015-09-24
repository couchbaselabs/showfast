package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/couchbaselabs/go-couchbase"
)

var ddocs = map[string]string{
	"metrics": `{
		"views": {
			"all": {
				"map": "function (doc, meta) {emit(meta.id, doc);}"
			}
		}
	}`,
	"clusters": `{
		"views": {
			"all": {
				"map": "function (doc, meta) {emit(meta.id, doc);}"
			}
		}
	}`,
	"benchmarks": `{
		"views": {
			"metrics_by_build": {
				 "map": "function (doc, meta) {if (!doc.obsolete) {emit(doc.build, doc.metric);}}"
			},
			"values_by_build_and_metric": {
				"map": "function (doc, meta) {if (!doc.obsolete) {emit([doc.metric, doc.build], doc.value);}}"
			},
			"value_and_snapshots_by_build_and_metric": {
				"map": "function (doc, meta) {emit([doc.metric, doc.build], [doc.value, doc.snapshots, doc.master_events, doc.build_url || null]);}"
			},
			"value_and_obsolete_by_build_and_metric": {
				"map": "function (doc, meta) {emit([doc.metric, doc.build], [doc.value, doc.obsolete == true]);}"
			}
		}
	}`,
	"feed": `{
		"views": {
			"all": {
				"map": "function (doc, meta) {emit(meta.id, doc);}"
			}
		}
	}`,
}

type DataSource struct {
	CouchbaseAddress, BucketPassword string
}

func (ds *DataSource) GetBucket(bucket string) *couchbase.Bucket {
	uri := fmt.Sprintf("http://%s:%s@%s/", bucket, ds.BucketPassword, ds.CouchbaseAddress)

	client, _ := couchbase.Connect(uri)
	pool, _ := client.GetPool("default")

	b, err := pool.GetBucket(bucket)
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}
	return b
}

func (ds *DataSource) QueryView(b *couchbase.Bucket, ddoc, view string,
	params map[string]interface{}) []couchbase.ViewRow {
	params["stale"] = false
	vr, err := b.View(ddoc, view, params)
	if err != nil {
		ds.installDDoc(ddoc)
	}
	return vr.Rows
}

func (ds *DataSource) installDDoc(ddoc string) {
	b := ds.GetBucket(ddoc) // bucket name == ddoc name
	err := b.PutDDoc(ddoc, ddocs[ddoc])
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (ds *DataSource) GetAllMetrics() []byte {
	b_metrics := ds.GetBucket("metrics")
	rows := ds.QueryView(b_metrics, "metrics", "all", map[string]interface{}{})

	metrics := []map[string]interface{}{}
	for i := range rows {
		metric := rows[i].Value.(map[string]interface{})
		metric["id"] = rows[i].ID
		metrics = append(metrics, metric)
	}

	j, _ := json.Marshal(metrics)
	return j
}

func (ds *DataSource) GetAllClusters() []byte {
	b_clusters := ds.GetBucket("clusters")
	rows := ds.QueryView(b_clusters, "clusters", "all", map[string]interface{}{})

	clusters := []map[string]interface{}{}
	for i := range rows {
		cluster := rows[i].Value.(map[string]interface{})
		cluster["Name"] = rows[i].ID
		clusters = append(clusters, cluster)
	}

	j, _ := json.Marshal(clusters)
	return j
}

type ByBuild [][]interface{}

func (b ByBuild) Len() int {
	return len(b)
}

func (b ByBuild) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByBuild) Less(i, j int) bool {
	build_i := strings.Split(b[i][0].(string), "-")
	build_j := strings.Split(b[j][0].(string), "-")
	if build_i[0] == build_j[0] {
		int_i, _ := strconv.ParseInt(build_i[1], 10, 16)
		int_j, _ := strconv.ParseInt(build_j[1], 10, 16)
		return int_i < int_j
	} else {
		return build_i[0] < build_j[0]
	}
}

func (ds *DataSource) GetAllTimelines() []byte {
	b_benchmarks := ds.GetBucket("benchmarks")
	rows := ds.QueryView(b_benchmarks, "benchmarks", "values_by_build_and_metric",
		map[string]interface{}{})

	timelines := map[string][][]interface{}{}
	for i := range rows {
		metric := rows[i].Key.([]interface{})[0]
		build := rows[i].Key.([]interface{})[1]
		value := rows[i].Value.(interface{})

		if array, ok := timelines[metric.(string)]; ok {
			timelines[metric.(string)] = append(array, []interface{}{build, value})
		} else {
			timelines[metric.(string)] = [][]interface{}{[]interface{}{build, value}}
		}
	}
	for _, timeline := range timelines {
		sort.Sort(ByBuild(timeline))
	}

	j, _ := json.Marshal(timelines)
	return j
}

func (ds *DataSource) GetAllRuns(metric string, build string) []byte {
	b_benchmarks := ds.GetBucket("benchmarks")
	params := map[string]interface{}{
		"startkey": []string{metric, build},
		"endkey":   []string{metric, build},
	}
	rows := ds.QueryView(b_benchmarks, "benchmarks", "value_and_snapshots_by_build_and_metric", params)

	benchmarks := []map[string]interface{}{}
	for i, row := range rows {
		var master_events string
		if str, ok := row.Value.([]interface{})[2].(string); ok {
			master_events = str
		} else {
			master_events = ""
		}
		var build_url string
		if str, ok := row.Value.([]interface{})[3].(string); ok {
			build_url = str
		} else {
			build_url = ""
                }
		benchmark := map[string]interface{}{
			"seq":           strconv.Itoa(i + 1),
			"value":         strconv.FormatFloat(row.Value.([]interface{})[0].(float64), 'f', 1, 64),
			"snapshots":     row.Value.([]interface{})[1],
			"master_events": master_events,
			"build_url": build_url,
		}
		benchmarks = append(benchmarks, benchmark)
	}
	j, _ := json.Marshal(benchmarks)
	return j
}

type Benchmark struct {
	ID        string   `json:"id"`
	Metric    string   `json:"metric"`
	Build     string   `json:"build"`
	Value     float64  `json:"value"`
	Obsolete  bool     `json:"obsolete"`
	Snapshots []string `json:"snapshots"`
}

func (ds *DataSource) GetAllBenchmarks() []byte {
	b_benchmarks := ds.GetBucket("benchmarks")
	rows := ds.QueryView(b_benchmarks, "benchmarks", "value_and_obsolete_by_build_and_metric",
		map[string]interface{}{})

	benchmarks := []Benchmark{}
	for _, row := range rows {
		benchmark := Benchmark{
			ID:       row.ID,
			Metric:   row.Key.([]interface{})[0].(string),
			Build:    row.Key.([]interface{})[1].(string),
			Value:    row.Value.([]interface{})[0].(float64),
			Obsolete: row.Value.([]interface{})[1].(bool),
		}
		benchmarks = append(benchmarks, benchmark)
	}
	j, _ := json.Marshal(benchmarks)
	return j
}

func (ds *DataSource) DeleteBenchmark(id string) {
	b_benchmarks := ds.GetBucket("benchmarks")
	b_benchmarks.Delete(id)
}

func (ds *DataSource) ReverseObsolete(id string) {
	b_benchmarks := ds.GetBucket("benchmarks")
	benchmark := Benchmark{}
	b_benchmarks.Get(id, &benchmark)
	benchmark.Obsolete = !benchmark.Obsolete
	err := b_benchmarks.Set(id, 0, benchmark)
	if err != nil {
		log.Println("Error updating benchmark:  %v", err)
	}
}

func appendIfUnique(slice []string, s string) []string {
	for i := range slice {
		if slice[i] == s {
			return slice
		}
	}
	return append(slice, s)
}

func (ds *DataSource) GetAllReleases() []byte {
	b_benchmarks := ds.GetBucket("benchmarks")
	rows := ds.QueryView(b_benchmarks, "benchmarks", "metrics_by_build",
		map[string]interface{}{})

	releases := []string{}
	for _, row := range rows {
		release := row.Key.(string)[:5]
		releases = appendIfUnique(releases, release)
	}

	j, _ := json.Marshal(releases)
	return j
}

func (ds *DataSource) GetComparison(baseline, target string) []byte {
	b_metrics := ds.GetBucket("metrics")
	b_benchmarks := ds.GetBucket("benchmarks")
	b_clusters := ds.GetBucket("clusters")
	rows := ds.QueryView(b_benchmarks, "benchmarks", "values_by_build_and_metric",
		map[string]interface{}{})

	metrics := map[string]map[string]interface{}{}
	for _, row := range rows {
		metric := row.Key.([]interface{})[0].(string)
		build := row.Key.([]interface{})[1].(string)
		value := row.Value.(float64)
		if _, ok := metrics[metric]; ok {
			if strings.HasPrefix(build, baseline) &&
				build > metrics[metric]["baseline"].(string) {
				metrics[metric]["baseline"] = build
				metrics[metric]["baseline_value"] = value
			}
			if strings.HasPrefix(build, target) &&
				build > metrics[metric]["target"].(string) {
				metrics[metric]["target"] = build
				metrics[metric]["target_value"] = value
			}
		} else {
			metrics[metric] = map[string]interface{}{
				"baseline":       build,
				"target":         build,
				"baseline_value": value,
				"target_value":   value,
			}
		}
	}
	reduced_metrics := map[string]map[string]interface{}{}
	for metric_name, builds := range metrics {
		if strings.HasPrefix(builds["baseline"].(string), baseline) &&
			strings.HasPrefix(builds["target"].(string), target) {
			metric := map[string]string{}
			b_metrics.Get(metric_name, &metric)
			cluster := map[string]string{}
			b_clusters.Get(metric["cluster"], &cluster)

			diff := 100 * (builds["target_value"].(float64) - builds["baseline_value"].(float64)) /
				builds["baseline_value"].(float64)

			var coeff float64
			if metric["larger_is_better"] == "false" {
				coeff = -1
			} else {
				coeff = 1
			}

			comparison := "The same"
			class := "same"
			if coeff*diff > 10 {
				diff := strconv.FormatFloat(diff*coeff, 'f', 1, 64)
				comparison = fmt.Sprintf("%s%% better", diff)
				class = "better"
			} else if coeff*diff < -10 {
				diff := strconv.FormatFloat(-diff*coeff, 'f', 1, 64)
				comparison = fmt.Sprintf("%s%% worse", diff)
				class = "worse"
			}

			reduced_metrics[metric_name] = map[string]interface{}{
				"title":      metric["title"],
				"cluster":    cluster,
				"baseline":   builds["baseline"],
				"target":     builds["target"],
				"comparison": comparison,
				"class":      class,
			}
		}
	}
	j, _ := json.Marshal(reduced_metrics)
	return j
}

func (ds *DataSource) GetAllFeedRecords() []byte {
	b_feed := ds.GetBucket("feed")
	params := map[string]interface{}{
		"descending": true,
		"limit":      100,
	}
	rows := ds.QueryView(b_feed, "feed", "all", params)

	records := []map[string]interface{}{}
	for _, row := range rows {
		record := row.Value.(map[string]interface{})
		records = append(records, record)
	}
	j, _ := json.Marshal(records)
	return j
}
