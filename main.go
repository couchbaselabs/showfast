package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "gopkg.in/inconshreveable/log15.v2"
)

const (
	address   = "0.0.0.0:8000"
)

var (
	dataSource     DataSource
	cbHost, cbPass string
)

func home(rw http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("app/index.html")
	rw.Write(content)
}

func admin(rw http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("app/admin.html")
	rw.Write(content)
}

func release(rw http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("app/release.html")
	rw.Write(content)
}

func feed(rw http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("app/feed.html")
	rw.Write(content)
}

func allRuns(rw http.ResponseWriter, r *http.Request) {
	metric := r.URL.Query()["metric"][0]
	build := r.URL.Query()["build"][0]

	writeJSON(rw, dataSource.getAllRuns(metric, build))
}

func getComparison(rw http.ResponseWriter, r *http.Request) {
	baseline := r.URL.Query()["baseline"][0]
	target := r.URL.Query()["target"][0]

	writeJSON(rw, dataSource.getComparison(baseline, target))
}

func deleteBenchmark(rw http.ResponseWriter, r *http.Request) {
	var params struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &params); err == nil {
		dataSource.deleteBenchmark(params.ID)
	}
}

func reverseObsolete(rw http.ResponseWriter, r *http.Request) {
	var params struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &params); err == nil {
		dataSource.reverseObsolete(params.ID)
	}
}

func init() {
	cbHost = os.Getenv("CB_HOST")
	cbPass = os.Getenv("CB_PASS")
	if cbHost == "" || cbPass == "" {
		log.Error("Missing Couchbase Server settings.")
		os.Exit(1)
	}
}

func main() {
	dataSource = DataSource{cbHost, cbPass}

	router := mux.NewRouter()

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/admin", admin).Methods("GET")
	router.HandleFunc("/feed", feed).Methods("GET")
	router.HandleFunc("/release", release).Methods("GET")

	router.HandleFunc("/all_metrics", dataSource.getAllMetrics).Methods("GET")
	router.HandleFunc("/all_clusters", dataSource.getAllClusters).Methods("GET")
	router.HandleFunc("/all_timelines", dataSource.getAllTimelines).Methods("GET")
	router.HandleFunc("/all_benchmarks", dataSource.getAllBenchmarks).Methods("GET")
	router.HandleFunc("/all_releases", dataSource.getAllReleases).Methods("GET")
	router.HandleFunc("/all_feed_records", dataSource.getAllFeedRecords).Methods("GET")

	router.HandleFunc("/all_runs", allRuns).Methods("GET")
	router.HandleFunc("/get_comparison", getComparison).Methods("GET")
	router.HandleFunc("/delete", deleteBenchmark).Methods("POST")
	router.HandleFunc("/reverse_obsolete", reverseObsolete).Methods("POST")

	router.PathPrefix("/css").Handler(http.FileServer(http.Dir("./app")))
	router.PathPrefix("/js").Handler(http.FileServer(http.Dir("./app")))
	router.PathPrefix("/partials").Handler(http.FileServer(http.Dir("./app")))

	http.Handle("/", router)

	banner := fmt.Sprintf("\n\t.:: Serving http://%s/ ::.\n", address)
        fmt.Println(banner)

	http.ListenAndServe(address, accessLog(http.DefaultServeMux))
}
