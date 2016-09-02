
showfast
========
[![Go Report Card](https://goreportcard.com/badge/github.com/couchbaselabs/showfast)](https://goreportcard.com/report/github.com/couchbaselabs/showfast)
[![codebeat badge](https://codebeat.co/badges/1fc2a490-a39e-49be-b218-2b9289da5ae7)](https://codebeat.co/projects/github-com-couchbaselabs-showfast)
[![Docker image](https://images.microbadger.com/badges/image/perflab/showfast.svg)](http://microbadger.com/images/perflab/showfast)

Couchbase Server performance dashboard.

Prerequisites
-------------

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
