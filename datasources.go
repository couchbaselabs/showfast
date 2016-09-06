package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/couchbase/go-couchbase"
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
	hostname, password string
}

func (ds *DataSource) getBucket(bucket string) *couchbase.Bucket {
	uri := fmt.Sprintf("http://%s:%s@%s:8091/", bucket, ds.password, ds.hostname)

	client, _ := couchbase.Connect(uri)
	pool, _ := client.GetPool("default")

	b, err := pool.GetBucket(bucket)
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}
	return b
}

func (ds *DataSource) queryView(b *couchbase.Bucket, ddoc, view string,
	params map[string]interface{}) []couchbase.ViewRow {
	params["stale"] = false
	vr, err := b.View(ddoc, view, params)
	if err != nil {
		ds.installDDoc(ddoc)
	}
	return vr.Rows
}

func (ds *DataSource) installDDoc(ddoc string) {
	b := ds.getBucket(ddoc) // bucket name == ddoc name
	err := b.PutDDoc(ddoc, ddocs[ddoc])
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (ds *DataSource) getAllMetrics(rw http.ResponseWriter, r *http.Request) {
	bMetrics := ds.getBucket("metrics")
	rows := ds.queryView(bMetrics, "metrics", "all", map[string]interface{}{})
	metrics := []map[string]interface{}{}
	for i := range rows {
		metric := rows[i].Value.(map[string]interface{})
		metric["id"] = rows[i].ID
		metrics = append(metrics, metric)
	}

	writeJSON(rw, metrics)
}

func (ds *DataSource) getAllClusters(rw http.ResponseWriter, r *http.Request) {
	bClusters := ds.getBucket("clusters")
	rows := ds.queryView(bClusters, "clusters", "all", map[string]interface{}{})

	clusters := []map[string]interface{}{}
	for i := range rows {
		cluster := rows[i].Value.(map[string]interface{})
		cluster["Name"] = rows[i].ID
		clusters = append(clusters, cluster)
	}

	writeJSON(rw, clusters)
}

type byBuild [][]interface{}

func (b byBuild) Len() int {
	return len(b)
}

func (b byBuild) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byBuild) Less(i, j int) bool {
	buildI := strings.Split(b[i][0].(string), "-")
	buildJ := strings.Split(b[j][0].(string), "-")
	if buildI[0] == buildJ[0] {
		intBuildI, _ := strconv.ParseInt(buildI[1], 10, 16)
		intBuildJ, _ := strconv.ParseInt(buildJ[1], 10, 16)
		return intBuildI < intBuildJ
	}
	return buildI[0] < buildJ[0]
}

func (ds *DataSource) getAllTimelines(rw http.ResponseWriter, r *http.Request) {
	bBenchmarks := ds.getBucket("benchmarks")
	rows := ds.queryView(bBenchmarks, "benchmarks", "values_by_build_and_metric",
		map[string]interface{}{})

	timelines := map[string][][]interface{}{}
	for i := range rows {
		metric := rows[i].Key.([]interface{})[0]
		build := rows[i].Key.([]interface{})[1]
		value := rows[i].Value.(interface{})

		if array, ok := timelines[metric.(string)]; ok {
			timelines[metric.(string)] = append(array, []interface{}{build, value})
		} else {
			timelines[metric.(string)] = [][]interface{}{{build, value}}
		}
	}
	for _, timeline := range timelines {
		sort.Sort(byBuild(timeline))
	}

	writeJSON(rw, timelines)
}

func (ds *DataSource) getAllRuns(metric string, build string) interface{} {
	bBenchmarks := ds.getBucket("benchmarks")
	params := map[string]interface{}{
		"startkey": []string{metric, build},
		"endkey":   []string{metric, build},
	}
	rows := ds.queryView(bBenchmarks, "benchmarks", "value_and_snapshots_by_build_and_metric", params)

	benchmarks := []map[string]interface{}{}
	for i, row := range rows {
		var masterEvents string
		if str, ok := row.Value.([]interface{})[2].(string); ok {
			masterEvents = str
		} else {
			masterEvents = ""
		}
		var buildURL string
		if val, ok := row.Value.([]interface{}); ok && len(val) > 3 {
			if str, ok := val[3].(string); ok {
				buildURL = str
			} else {
				buildURL = ""
			}
		} else {
			buildURL = ""
		}
		benchmark := map[string]interface{}{
			"seq":           strconv.Itoa(i + 1),
			"value":         strconv.FormatFloat(row.Value.([]interface{})[0].(float64), 'f', 1, 64),
			"snapshots":     row.Value.([]interface{})[1],
			"master_events": masterEvents,
			"build_url":     buildURL,
		}
		benchmarks = append(benchmarks, benchmark)
	}

	return benchmarks
}

type Benchmark struct {
	ID        string   `json:"id"`
	Build     string   `json:"build"`
	BuildURL  string   `json:"build_url"`
	Metric    string   `json:"metric"`
	Obsolete  bool     `json:"obsolete"`
	Snapshots []string `json:"snapshots"`
	Value     float64  `json:"value"`
}

func (ds *DataSource) getAllBenchmarks(rw http.ResponseWriter, r *http.Request) {
	bBenchmarks := ds.getBucket("benchmarks")
	rows := ds.queryView(bBenchmarks, "benchmarks", "value_and_obsolete_by_build_and_metric",
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

	writeJSON(rw, benchmarks)
}

func (ds *DataSource) deleteBenchmark(id string) {
	bBenchmarks := ds.getBucket("benchmarks")
	bBenchmarks.Delete(id)
}

func (ds *DataSource) reverseObsolete(id string) {
	bBenchmarks := ds.getBucket("benchmarks")
	benchmark := Benchmark{}
	bBenchmarks.Get(id, &benchmark)
	benchmark.Obsolete = !benchmark.Obsolete
	err := bBenchmarks.Set(id, 0, benchmark)
	if err != nil {
		log.Printf("Error updating benchmark:  %v\n", err)
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

func (ds *DataSource) getAllReleases(rw http.ResponseWriter, r *http.Request) {
	bBenchmarks := ds.getBucket("benchmarks")
	rows := ds.queryView(bBenchmarks, "benchmarks", "metrics_by_build",
		map[string]interface{}{})

	releases := []string{}
	for _, row := range rows {
		release := row.Key.(string)[:5]
		releases = appendIfUnique(releases, release)
	}

	writeJSON(rw, releases)
}

func (ds *DataSource) getComparison(baseline, target string) interface{} {
	bMetrics := ds.getBucket("metrics")
	bBenchmarks := ds.getBucket("benchmarks")
	bClusters := ds.getBucket("clusters")

	rows := ds.queryView(bBenchmarks, "benchmarks", "values_by_build_and_metric",
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

	reducedMetrics := map[string]map[string]interface{}{}
	for metricName, builds := range metrics {
		if strings.HasPrefix(builds["baseline"].(string), baseline) &&
			strings.HasPrefix(builds["target"].(string), target) {
			metric := map[string]string{}
			bMetrics.Get(metricName, &metric)
			cluster := map[string]string{}
			bClusters.Get(metric["cluster"], &cluster)

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

			reducedMetrics[metricName] = map[string]interface{}{
				"title":      metric["title"],
				"cluster":    cluster,
				"baseline":   builds["baseline"],
				"target":     builds["target"],
				"comparison": comparison,
				"class":      class,
			}
		}
	}

	return reducedMetrics
}

func (ds *DataSource) getAllFeedRecords(rw http.ResponseWriter, r *http.Request) {
	bFeed := ds.getBucket("feed")
	params := map[string]interface{}{
		"descending": true,
		"limit":      100,
	}
	rows := ds.queryView(bFeed, "feed", "all", params)

	records := []map[string]interface{}{}
	for _, row := range rows {
		record := row.Value.(map[string]interface{})
		records = append(records, record)
	}

	writeJSON(rw, records)
}
