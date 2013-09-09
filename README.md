showfast
========

Performance dashboard

Prerequisites
-------------

* go 1.1
* Couchbase Server 2.0.0 or higher

Required buckets
----------------

* benchmarks
* clusters
* metrics

Installation
------------

    go get github.com/couchbaselabs/showfast

Usage
-----

    $GOPATH/bin/showfast -address=127.0.0.1:8000
