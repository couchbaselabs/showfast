package datasource

import (
	"fmt"
	"log"

	"github.com/couchbaselabs/go-couchbase"
)

const cbHost = "http://%s:password@127.0.0.1:8091/"

var ddocs = map[string]string{
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
				 "map": "function (doc, meta) {if (!doc.obsolete) {emit(doc.build, doc.metric);}}"
			},
			"values_by_build_and_metric": {
				"map": "function (doc, meta) {if (!doc.obsolete) {emit([doc.metric, doc.build], doc.value);}}"
			},
			"value_and_href_by_build_and_metric": {
				"map": "function (doc, meta) {emit([doc.metric, doc.build], [doc.value, doc.report]);}"
			}
		}
	}`,
}

func GetBucket(bucket string) *couchbase.Bucket {
	uri := fmt.Sprintf(cbHost, bucket)
	client, _ := couchbase.Connect(uri)
	pool, _ := client.GetPool("default")
	b, err := pool.GetBucket(bucket)
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}
	return b
}

func QueryView(b *couchbase.Bucket, ddoc, view string,
	params map[string]interface{}) []couchbase.ViewRow {
	params["stale"] = false
	vr, err := b.View(ddoc, view, params)
	if err != nil {
		InstallDDoc(ddoc)
	}
	return vr.Rows
}

func InstallDDoc(ddoc string) {
	b := GetBucket(ddoc) // bucket name == ddoc name
	err := b.PutDDoc(ddoc, ddocs[ddoc])
	if err != nil {
		log.Fatalf("%v", err)
	}
}
