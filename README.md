[![Go Report Card](https://goreportcard.com/badge/github.com/couchbaselabs/showfast)](https://goreportcard.com/report/github.com/couchbaselabs/showfast)

showfast
========

Performance dashboard. Powered by web.go, AngularJS, nvd3 and Couchbase Server.

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

    $GOPATH/bin/showfast

 All parameters are specified in config.json
 