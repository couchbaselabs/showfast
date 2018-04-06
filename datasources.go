package main

import (
	"fmt"
	"os"

	"gopkg.in/couchbase/gocb.v1"
	log "gopkg.in/inconshreveable/log15.v2"
)

var couchbaseBuckets = []string{"benchmarks", "clusters", "metrics"}

type dataStore struct {
	cluster *gocb.Cluster
	buckets map[string]*gocb.Bucket
}

func newDataStore() *dataStore {
	hostname := os.Getenv("CB_HOST")
	if hostname == "" {
		log.Error("missing Couchbase Server hostname")
		os.Exit(1)
	}

	connSpecStr := fmt.Sprintf("couchbase://%s", hostname)
	cluster, err := gocb.Connect(connSpecStr)
	if err != nil {
		log.Error("failed to connect to Couchbase Server", "err", err)
		os.Exit(1)
	}

	return &dataStore{cluster, map[string]*gocb.Bucket{}}
}

func (ds *dataStore) auth() {
	password := os.Getenv("CB_PASS")
	if password == "" {
		log.Error("missing password")
		os.Exit(1)
	}

	authMap := gocb.BucketAuthenticatorMap{}
	for _, bucketName := range couchbaseBuckets {
		bucket, err := ds.cluster.OpenBucket(bucketName, password)
		if err != nil {
			log.Error("failed to open bucket", "err", err)
			os.Exit(1)
		}
		ds.buckets[bucketName] = bucket
		authMap[bucketName] = gocb.BucketAuthenticator{Password: password}
	}

	auth := gocb.ClusterAuthenticator{
		Username: "Administrator",
		Password: password,
		Buckets:  authMap,
	}
	if err := ds.cluster.Authenticate(auth); err != nil {
		log.Error("authentication failed", "err", err)
		os.Exit(1)
	}
}

func (ds *dataStore) getBucket(bucketName string) *gocb.Bucket {
	return ds.buckets[bucketName]
}

type Build struct {
	Build string `json:"build"`
}

func (ds *dataStore) getBuilds() (*[]string, error) {
	var builds []string

	query := gocb.NewN1qlQuery(
		"SELECT DISTINCT `build` " +
			"FROM benchmarks " +
			"WHERE hidden = False;")

	rows, err := ds.cluster.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		return &builds, err
	}

	var row Build
	for rows.Next(&row) {
		builds = append(builds, row.Build)
	}
	return &builds, nil
}

type Metric struct {
	Cluster     string `json:"cluster"`
	Category    string `json:"category"`
	Component   string `json:"component"`
	ID          string `json:"id"`
	OrderBy     string `json:"orderBy"`
	SubCategory string `json:"subCategory"`
	Title       string `json:"title"`
}

type MetricWithCluster struct {
	Metric
	Cluster Cluster `json:"cluster"`
}

func (ds *dataStore) addMetric(metric Metric) error {
	bucket := ds.getBucket("metrics")

	_, err := bucket.Upsert(metric.ID, metric, 0)
	if err != nil {
		log.Error("failed to insert metric", "err", err)
	}
	return err
}

func (ds *dataStore) getMetrics(component, category string, subCategory string) (*[]MetricWithCluster, error) {
	var metrics []MetricWithCluster

	query := gocb.NewN1qlQuery(
		"SELECT m.id, m.title, m.component, m.category, m.orderBy, m.subCategory, c AS `cluster` " +
			"FROM metrics m JOIN clusters c ON KEYS m.`cluster`" +
			"WHERE m.component = $1 AND m.category = $2 AND m.subCategory = $3;")
	params := []interface{}{component, category, subCategory}

	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		return &metrics, err
	}

	var row MetricWithCluster
	for rows.Next(&row) {
		metrics = append(metrics, row)
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

func (ds *dataStore) addCluster(cluster Cluster) error {
	bucket := ds.getBucket("clusters")

	_, err := bucket.Upsert(cluster.Name, cluster, 0)
	if err != nil {
		log.Error("failed to insert cluster", "err", err)
	}
	return err
}

type Benchmark struct {
	Build     string   `json:"build"`
	BuildURL  string   `json:"buildURL"`
	DateTime  string   `json:"dateTime"`
	ID        string   `json:"id"`
	Metric    string   `json:"metric"`
	Hidden    bool     `json:"hidden"`
	Snapshots []string `json:"snapshots"`
	Value     float64  `json:"value"`
}

func (ds *dataStore) findExisting(benchmark Benchmark) ([]Benchmark, error) {
	var existing []Benchmark

	query := gocb.NewN1qlQuery(
		"SELECT `build`, buildURL, dateTime, id, metric, hidden, snapshots, `value` " +
			"FROM benchmarks " +
			"WHERE metric = $1 AND `build` = $2;")
	params := []interface{}{benchmark.Metric, benchmark.Build}

	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		return existing, err
	}

	var row Benchmark
	for rows.Next(&row) {
		existing = append(existing, row)
		row = Benchmark{}
	}

	return existing, nil
}

func (ds *dataStore) hideExisting(benchmark Benchmark) error {
	existing, err := ds.findExisting(benchmark)
	if err != nil {
		log.Error("failed to find existing benchmarks", "err", err)
		return err
	}

	bucket := ds.getBucket("benchmarks")

	for _, b := range existing {
		b.Hidden = true
		if _, err := bucket.Upsert(b.ID, b, 0); err != nil {
			log.Error("update failed", "err", err)
			return err
		}
	}
	return nil
}

func (ds *dataStore) addBenchmark(benchmark Benchmark) error {
	if err := ds.hideExisting(benchmark); err != nil {
		return err
	}

	bucket := ds.getBucket("benchmarks")

	_, err := bucket.Upsert(benchmark.ID, benchmark, 0)
	if err != nil {
		log.Error("failed to insert benchmark", "err", err)
	}
	return err
}

func (ds *dataStore) getBenchmarks(component, category string, subCategory string) (*[]Benchmark, error) {
	var benchmarks []Benchmark

	var query *gocb.N1qlQuery
	var params []interface{}

	if subCategory == "" {
		query = gocb.NewN1qlQuery(
			"SELECT b.`build`, b.id, b.hidden, b.metric, b.`value` " +
				"FROM metrics m " +
				"JOIN benchmarks b " +
				"ON KEY b.metric FOR m " +
				"WHERE m.component = $1 AND m.category = $2" +
				"ORDER BY b.metric, b.`build` DESC, b.hidden;")
		params = []interface{}{component, category}
	} else {
		query = gocb.NewN1qlQuery(
			"SELECT b.`build`, b.id, b.hidden, b.metric, b.`value` " +
				"FROM metrics m " +
				"JOIN benchmarks b " +
				"ON KEY b.metric FOR m " +
				"WHERE m.component = $1 AND m.category = $2 AND m.subCategory = $3" +
				"ORDER BY b.metric, b.`build` DESC, b.hidden;")
		params = []interface{}{component, category, subCategory}
	}

	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		return &benchmarks, err
	}

	var row Benchmark
	for rows.Next(&row) {
		benchmarks = append(benchmarks, row)
		row = Benchmark{}
	}

	return &benchmarks, nil
}

type Run struct {
	Build string  `json:"build"`
	Value float64 `json:"value"`
}

func (ds *dataStore) getTimeline(metric string) (*[][]interface{}, error) {
	var runs [][]interface{}

	query := gocb.NewN1qlQuery(
		"SELECT `build`, `value` " +
			"FROM benchmarks " +
			"WHERE metric = $1 AND hidden = false " +
			"ORDER BY `build`;")
	params := []interface{}{metric}

	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		return &runs, err
	}

	var row Run
	for rows.Next(&row) {
		run := []interface{}{row.Build, row.Value}
		runs = append(runs, run)
	}

	return &runs, nil
}

func (ds *dataStore) getAllRuns(metric string, build string) (*[]Benchmark, error) {
	var benchmarks []Benchmark

	query := gocb.NewN1qlQuery(
		"SELECT `build`, buildURL, dateTime, snapshots, `value` " +
			"FROM benchmarks " +
			"WHERE metric = $1 AND `build` = $2;")
	params := []interface{}{metric, build}

	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		return &benchmarks, err
	}

	var row Benchmark
	for rows.Next(&row) {
		benchmarks = append(benchmarks, row)
		row = Benchmark{}
	}
	return &benchmarks, nil
}

type Comparison struct {
	Category  string `json:"category"`
	Component string `json:"component"`
	OS        string `json:"os"`
	Title     string `json:"title"`
	Results   []Run  `json:"results"`
}

func (ds *dataStore) compare(build1, build2 string) (*[]Comparison, error) {
	var comparison []Comparison

	query := gocb.NewN1qlQuery(
		"SELECT m.component, m.category, m.title, c.os, ARRAY_AGG({\"build\": b.`build`, \"value\": b.`value`}) AS results " +
			"FROM benchmarks b USE INDEX (benchmarks_comparison) " +
			"JOIN metrics m ON KEYS b.metric " +
			"JOIN clusters c ON KEYS m.`cluster` " +
			"WHERE b.hidden = False " +
			"AND (b.`build` = $1 OR b.`build` = $2) " +
			"GROUP BY m.component, m.category, m.title, m.`cluster`, c.os " +
			"HAVING COUNT(*) > 1 " +
			"ORDER BY c.os, m.title;")
	params := []interface{}{build1, build2}

	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		return &comparison, err
	}
	var row Comparison
	for rows.Next(&row) {
		comparison = append(comparison, row)
		row = Comparison{}
	}
	return &comparison, nil
}

func (ds *dataStore) deleteBenchmark(key string) error {
	bucket := ds.getBucket("benchmarks")

	_, err := bucket.Remove(key, 0)
	if err != nil {
		log.Error("benchmark deletion failed", "err", err)
	}
	return err
}

func (ds *dataStore) reverseHidden(key string) error {
	bucket := ds.getBucket("benchmarks")

	benchmark := Benchmark{}
	if _, err := bucket.Get(key, &benchmark); err != nil {
		log.Error("failed to get benchmark", "err", err)
		return err
	}
	benchmark.Hidden = !benchmark.Hidden

	if _, err := bucket.Upsert(key, benchmark, 0); err != nil {
		log.Error("benchmark update failed", "err", err)
		return err
	}

	return nil
}
