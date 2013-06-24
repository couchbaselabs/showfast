package datasource

import (
	"log"

	"github.com/couchbaselabs/go-couchbase"
)

var Ddocs = map[string]string{
	"metrics": `{
		"views": {
			"all": {
				"map": "function (doc, meta) {emit(meta.id, doc);}"
			}
		}
	}`,
	"benchmarks": `{
		"views": {
			"metrics_by_build": {
				 "map": "function (doc, meta) {emit(doc.build, doc.metric);}"
			},
			"build_by_metric": {
				"map": "function (doc, meta) {emit(doc.metric, doc.build);}"
			},
			"build_and_value_by_metric": {
				"map": "function (doc, meta) {emit(doc.metric, [doc.build, doc.value]);}"
			},
			"values_by_build_and_metric": {
				"map": "function (doc, meta) {emit([doc.metric, doc.build], doc.value);}"
			}
		}
	}`,
}

func GetPool() (pool couchbase.Pool) {
	c, err := couchbase.Connect("http://127.0.0.1:8091/")
	if err != nil {
		log.Fatalf("Error connecting:  %v", err)
	}
	pool, _ = c.GetPool("default")
	return
}

func InstallDDoc(bucket string) {
	pool := GetPool()
	b, _ := pool.GetBucket(bucket)
	err := b.PutDDoc(bucket, Ddocs[bucket])
	if err != nil {
		log.Fatalf("%v", err)
	}
}
