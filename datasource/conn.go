package datasource

import (
	"log"

	"github.com/couchbaselabs/go-couchbase"
)

func GetPool() (pool couchbase.Pool) {
	c, err := couchbase.Connect("http://127.0.0.1:8091/")
	if err != nil {
		log.Fatalf("Error connecting:  %v", err)
	}
	pool, _ = c.GetPool("default")
	return
}
