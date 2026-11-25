
showfast
========

Couchbase Server performance dashboard.

Prerequisites
-------------

* Couchbase Server 4.x.

Building from the source
------------------------

First, install Go vendor tool:
```
> go get -u github.com/kardianos/govendor
```

Pull third-party packages into vendor folder:
```
> govendor sync
```

Build the project:
```
> go build
```

Usage
-----

```
> docker pull perflab/showfast
> docker run -t -i -e CB_HOST=... -e CB_PASS=... -p 8000:8000 perflab/showfast
```
