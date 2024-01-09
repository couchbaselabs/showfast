package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/couchbase/gocb.v1"
	log "gopkg.in/inconshreveable/log15.v2"
)

var couchbaseBuckets = []string{"benchmarks", "clusters", "metrics", "impressive"}

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

	connSpecStr := fmt.Sprintf("couchbases://%s", hostname)
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

    username := os.Getenv("CB_USER")
	if username == "" {
		log.Error("missing username")
		os.Exit(1)
	}

	ds.cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: username,
		Password: password,
	})

	authMap := gocb.BucketAuthenticatorMap{}
	for _, bucketName := range couchbaseBuckets {
		bucket, err := ds.cluster.OpenBucket(bucketName, "")
		if err != nil {
			log.Error("failed to open bucket", "err", err)
			os.Exit(1)
		}
		ds.buckets[bucketName] = bucket
		authMap[bucketName] = gocb.BucketAuthenticator{Password: password}
	}

	// auth := gocb.ClusterAuthenticator{
	// Username: Username,
	// Password: password,
	// Buckets:  authMap,
	// }
	// if err := ds.cluster.Authenticate(auth); err != nil {
	// log.Error("authentication failed", "err", err)
	// os.Exit(1)
	// }
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
	Chirality   int    `json:"chirality"`
	MemQuota    int64  `json:"memquota"`
	Provider    string `json:"provider"`
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

func (ds *dataStore) getMetrics(component, category string, subCategory string) (*[]interface{}, error) {
	var metrics []interface{}

	query := gocb.NewN1qlQuery(
		"SELECT m.id, m.title, m.component, m.category, m.orderBy, m.subCategory,  m.memquota, m.provider, c AS `cluster` " +
			"FROM metrics m JOIN clusters c ON KEYS m.`cluster`" +
			"WHERE m.component = $1 AND m.category = $2 AND m.subCategory = $3 " +
			"ORDER BY m.category;")
	params := []interface{}{component, category, subCategory}

	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		return &metrics, err
	}

	var row interface{}
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

type Test struct {
	Category string  `json:"category"`
	Title    string  `json:"title"`
	Metric   string  `json:"metric"`
	Active   bool    `json:"active"`
	Windows  bool    `json:"windows"`
	B1       string  `json:"b1"`
	B2       string  `json:"b2"`
	V1       float64 `json:"v1"`
	V2       float64 `json:"v2"`
	JobUrl1  string  `json:"jobUrl1"`
	JobUrl2  string  `json:"jobUrl2"`
}

func (ds *dataStore) findExisting(benchmark Benchmark) ([]Benchmark, error) {
	existing := []Benchmark{}

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

func getJobLink(metric string, build string) (string, error) {
	var url string
	var query *gocb.N1qlQuery
	var params []interface{}

	query = gocb.NewN1qlQuery(
		"SELECT RAW url FROM jenkins WHERE " +
			"SPLIT(SPLIT(test_config,'/')[-1], '.')[0] = $1 and version=$2")

	params = []interface{}{metric, build}
	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		log.Error(err.Error())
		return url, err
	}
	var row string
	for rows.Next(&row) {
		url = row
	}

	return url, err
}

func getClusterNames() *[]string {
	clusterNames := []string{}
	query := gocb.NewN1qlQuery("SELECT RAW meta().id from clusters")
	rows, err := ds.cluster.ExecuteN1qlQuery(query, nil)
	if err != nil {
		log.Error(err.Error())
		return &clusterNames
	}

	var row string
	for rows.Next(&row) {
		clusterNames = append(clusterNames, row)
	}
	return &clusterNames
}

func cutoffClusterName(fullName string, clusterNames *[]string) string {
	for _, v := range *clusterNames {
		tokens := strings.Split(fullName, "_"+v)
		if len(tokens) > 1 {
			return tokens[0]
		}
	}
	return fullName
}

func (ds *dataStore) getImpressiveTests(component string, build1 string, build2 string, active bool) (*[]Test, error) {
	tests := []Test{}
	var query *gocb.N1qlQuery
	var params []interface{}
	isActivePredicate := ""
	if active {
		isActivePredicate = "AND i.active = True "
	}

	query = gocb.NewN1qlQuery(
		"SELECT i.category, i.title, i.active, i.windows, i.metric, " +
			"ARRAY_AGG({b.`value`, b.`build`})[0].`build` as `b1`, " +
			"ARRAY_AGG({b.`value`, b.`build`})[0].`value` as `v1`, " +
			"ARRAY_AGG({b.`value`, b.`build`})[1].`build` as `b2`, " +
			"ARRAY_AGG({b.`value`, b.`build`})[1].`value` as `v2` " +
			"FROM impressive i LEFT JOIN benchmarks b ON i.metric = b.metric AND b.`build` IN [$3, $4] " +
			"WHERE i.type = $1 AND i.component = $2 " + isActivePredicate +
			"GROUP BY i.category, i.title, i.active, i.windows, i.metric " +
			"ORDER BY i.category")

	params = []interface{}{"test", component, build1, build2}
	rows, err := ds.cluster.ExecuteN1qlQuery(query, params)
	if err != nil {
		log.Error(err.Error())
		return &tests, err
	}
	var row Test
	for rows.Next(&row) {
		if row.B1 == build2 || row.B2 == build1 {
			tmpB := row.B1
			tmpV := row.V1
			row.B1 = row.B2
			row.V1 = row.V2
			row.B2 = tmpB
			row.V2 = tmpV
		}

		tests = append(tests, row)
		row = Test{}
	}

	clusterNames := getClusterNames()

	for k, v := range tests {
		if v.B2 == "" {
			url, _ := getJobLink(cutoffClusterName(v.Metric, clusterNames), build2)
			tests[k].JobUrl2 = url
		}
		if v.B1 == "" {
			url, _ := getJobLink(cutoffClusterName(v.Metric, clusterNames), build1)
			tests[k].JobUrl1 = url
		}
	}

	return &tests, nil
}

func (ds *dataStore) getBenchmarks(component, category string, subCategory string) (*[]Benchmark, error) {
	benchmarks := []Benchmark{}

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
	runs := [][]interface{}{}

	query := gocb.NewN1qlQuery(
		"SELECT `build`, `value` " +
			"FROM benchmarks " +
			"WHERE metric = $1 AND hidden = false " +
			"ORDER BY SPLIT(`build`, '-')[0], TONUMBER(SPLIT(`build`, '-')[1]);")

	//this is a heck to show average of gsi rebalance tests for smart batching
	// if strings.Contains(metric, "aether_5indexes") {
	// 	query = gocb.NewN1qlQuery(
	// 		"SELECT `build`, ROUND(AVG(`value`),3) as `value` " +
	// 			"FROM benchmarks " +
	// 			"WHERE metric = $1 " +
	// 			"Group by `build`" +
	// 			"ORDER BY SPLIT(`build`, '-')[0], TONUMBER(SPLIT(`build`, '-')[1]);")
	// }
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
	benchmarks := []Benchmark{}

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
	comparison := []Comparison{}

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
