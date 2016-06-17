
showfast
========
[![Go Report Card](https://goreportcard.com/badge/github.com/couchbaselabs/showfast)](https://goreportcard.com/report/github.com/couchbaselabs/showfast)

Couchbase Server performance dashboard.

Prerequisites
-------------

* go 1.6
* Couchbase Server 4.x.

Required buckets
----------------

* benchmarks
* clusters
* feed
* metrics

Usage
-----

    docker pull perflab/showfast
    docker run -t -i -e CB_HOST=... -e CB_PASS=... -p 8000:8000 perflab/showfast
