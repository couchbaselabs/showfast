package datasource

import (
	"encoding/json"
	"log"
	"sort"
)

type Timeline [][]interface{}

func (b Timeline) Len() int {
	return len(b)
}

func (b Timeline) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Timeline) Less(i, j int) bool {
	return b[i][0].(string) < b[j][0].(string)
}

func GetTimeline(metric string) (t []byte) {
	pool := GetPool()
	b_benchmarks, err := pool.GetBucket("benchmarks")
	if err != nil {
		log.Fatalf("Error reading bucket:  %v", err)
	}

	res, err := b_benchmarks.View("benchmarks", "by_metric",
		map[string]interface{}{"stale": false, "key": metric})
	if err != nil {
		log.Fatalf("Error reading view:  %v", err)
	}

	timeline := Timeline{}
	for i := range res.Rows {
		xy := res.Rows[i].Value.([]interface{})
		timeline = append(timeline, xy)
	}
	sort.Sort(timeline)
	t, _ = json.Marshal(timeline)
	return
}
