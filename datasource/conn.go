package datasource

import (
	"log"

	"github.com/couchbaselabs/go-couchbase"
)

const cbHost = "http://127.0.0.1:8091/"

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
			}
		}
	}`,
}

func GetBucket(bucket string) *couchbase.Bucket {
	b, err := couchbase.GetBucket(cbHost, "default", bucket)
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
