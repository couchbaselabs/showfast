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
				"map": "function (doc, meta) {emit([doc.metric, doc.build], [doc.value, doc.snapshots, doc.master_events, doc.build_url, doc.datetime]);}"
			},
			"value_and_obsolete_by_build_and_metric": {
				"map": "function (doc, meta) {emit([doc.metric, doc.build], [doc.value, doc.obsolete == true]);}"
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
	for _, row := range rows {
		var buildURL, dateTime, masterEvents string
		if str, ok := row.Value.([]interface{})[2].(string); ok {
			masterEvents = str
		}
		if str, ok := row.Value.([]interface{})[3].(string); ok {
			buildURL = str
		}
		if str, ok := row.Value.([]interface{})[4].(string); ok {
			dateTime = str
		}
		benchmark := map[string]interface{}{
			"value":         strconv.FormatFloat(row.Value.([]interface{})[0].(float64), 'f', 1, 64),
			"snapshots":     row.Value.([]interface{})[1],
			"master_events": masterEvents,
			"build_url":     buildURL,
			"datetime":      dateTime,
		}
		benchmarks = append(benchmarks, benchmark)
	}

	return benchmarks
}

type Benchmark struct {
	ID        string   `json:"id"`
	Build     string   `json:"build"`
	BuildURL  string   `json:"build_url"`
	DateTime  string   `json:"datetime"`
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
