package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/couchbase/go-couchbase"
	"github.com/mitchellh/mapstructure"
	log "gopkg.in/inconshreveable/log15.v2"
)

type dataStore struct {
	buckets map[string]*couchbase.Bucket
}

func newBucket(baseURL, bucketName, password string) (*couchbase.Bucket, error) {
	client, err := couchbase.ConnectWithAuthCreds(baseURL, bucketName, password)
	if err != nil {
		log.Error("error connecting to server", "err", err)
		return nil, err
	}

	pool, err := client.GetPool("default")
	if err != nil {
		log.Error("error getting pool", "err", err)
		return nil, err
	}

	bucket, err := pool.GetBucket(bucketName)
	if err != nil {
		log.Error("error getting bucket", "err", err)
		return nil, err
	}

	return bucket, nil
}

func newDataStore() *dataStore {
	hostname := os.Getenv("CB_HOST")
	password := os.Getenv("CB_PASS")
	if hostname == "" || password == "" {
		log.Error("missing Couchbase Server settings", "hostname", hostname, "pass", password)
		os.Exit(1)
	}

	baseURL := fmt.Sprintf("http://%s:8091/", hostname)

	buckets := map[string]*couchbase.Bucket{}
	for _, bucketName := range []string{"benchmarks", "clusters", "metrics"} {
		bucket, err := newBucket(baseURL, bucketName, password)
		if err != nil {
			os.Exit(1)
		}
		buckets[bucketName] = bucket
	}

	return &dataStore{buckets}
}

func (ds *dataStore) getBucket(bucketName string) *couchbase.Bucket {
	return ds.buckets[bucketName]
}

func (ds *dataStore) queryView(b *couchbase.Bucket, ddoc, view string, params map[string]interface{}) (*[]couchbase.ViewRow, error) {
	viewResult, err := b.View(ddoc, view, params)
	if err != nil {
		log.Error("error quering view", "err", err)
		return nil, err
	}
	return &viewResult.Rows, nil
}

type Metric struct {
	Cluster string `json:"cluster"`
	ID      string `json:"id"`
	OrderBy string `json:"orderBy"`
	Title   string `json:"title"`
}

func (ds *dataStore) getAllMetrics() (*[]Metric, error) {
	bucket := ds.getBucket("metrics")
	rows, err := ds.queryView(bucket, "v1", "all", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	metrics := []Metric{}
	for _, row := range *rows {
		value, ok := row.Value.(map[string]interface{})
		if !ok {
			log.Error("assertion failed", "value", row.Value)
			continue
		}
		var metric Metric
		if err := mapstructure.Decode(value, &metric); err != nil {
			log.Error("conversion failed", "err", err)
			continue
		}
		metrics = append(metrics, metric)
	}

	return &metrics, nil
}

type Cluster struct {
	CPU    string `json:"cpu"`
	Disk   string `json:"disk"`
	Memory string `json:"memory"`
	Name   string `json:"name"`
	OS     string `json:"os"`
}

func (ds *dataStore) getAllClusters() (*[]Cluster, error) {
	bucket := ds.getBucket("clusters")
	rows, err := ds.queryView(bucket, "v1", "all", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	clusters := []Cluster{}
	for _, row := range *rows {
		value, ok := row.Value.(map[string]interface{})
		if !ok {
			log.Error("assertion failed", "value", row.Value)
			continue
		}
		var cluster Cluster
		if err := mapstructure.Decode(value, &cluster); err != nil {
			log.Error("conversion failed", "err", err)
			continue
		}
		clusters = append(clusters, cluster)
	}

	return &clusters, nil
}

type Benchmark struct {
	Build     string   `json:"build"`
	BuildURL  string   `json:"buildURL"`
	DateTime  string   `json:"dateTime"`
	ID        string   `json:"id"`
	Metric    string   `json:"metric"`
	Obsolete  bool     `json:"obsolete"`
	Snapshots []string `json:"snapshots"`
	Value     float64  `json:"value"`
}

func (ds *dataStore) getBenchmarks(view string) (*[]Benchmark, error) {
	bucket := ds.getBucket("benchmarks")
	queryParams := map[string]interface{}{}
	rows, err := ds.queryView(bucket, "v1", view, queryParams)
	if err != nil {
		return nil, err
	}

	benchmarks := []Benchmark{}
	for _, row := range *rows {
		value, ok := row.Value.(map[string]interface{})
		if !ok {
			log.Error("assertion failed", "value", row.Value)
			continue
		}
		var benchmark Benchmark
		if err := mapstructure.Decode(value, &benchmark); err != nil {
			log.Error("conversion failed", "err", err)
			continue
		}
		benchmarks = append(benchmarks, benchmark)
	}

	return &benchmarks, nil
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

func (ds *dataStore) getAllTimelines() (*map[string][][]interface{}, error) {
	benchmarks, err := ds.getBenchmarks("all_visible")
	if err != nil {
		return nil, err
	}

	timelines := map[string][][]interface{}{}
	for _, benchmark := range *benchmarks {
		if benchmark.Obsolete {
			continue
		}
		if array, ok := timelines[benchmark.Metric]; ok {
			timelines[benchmark.Metric] = append(array, []interface{}{benchmark.Build, benchmark.Value})
		} else {
			timelines[benchmark.Metric] = [][]interface{}{{benchmark.Build, benchmark.Value}}
		}
	}
	for _, timeline := range timelines {
		sort.Sort(byBuild(timeline))
	}
	return &timelines, nil
}

func (ds *dataStore) getAllRuns(metric string, build string) (*[]Benchmark, error) {
	bucket := ds.getBucket("benchmarks")
	params := map[string]interface{}{
		"startkey": []string{metric, build},
		"endkey":   []string{metric, build},
	}
	rows, err := ds.queryView(bucket, "v1", "by_metric_and_build", params)
	if err != nil {
		return nil, err
	}

	benchmarks := []Benchmark{}
	for _, row := range *rows {
		value, ok := row.Value.(map[string]interface{})
		if !ok {
			log.Error("assertion failed", "value", row.Value)
			continue
		}
		var benchmark Benchmark
		if err := mapstructure.Decode(value, &benchmark); err != nil {
			log.Error("conversion failed", "err", err)
			continue
		}
		benchmarks = append(benchmarks, benchmark)
	}
	return &benchmarks, nil
}

func (ds *dataStore) deleteBenchmark(id string) error {
	bucket := ds.getBucket("benchmarks")
	err := bucket.Delete(id)
	if err != nil {
		log.Error("deletion failed", "err", err)
	}
	return err
}

func (ds *dataStore) reverseObsolete(id string) error {
	bucket := ds.getBucket("benchmarks")
	benchmark := Benchmark{}

	if err := bucket.Get(id, &benchmark); err != nil {
		return err
	}
	benchmark.Obsolete = !benchmark.Obsolete

	err := bucket.Set(id, 0, benchmark)
	if err != nil {
		log.Error("update failed", "err", err)
	}
	return err
}
