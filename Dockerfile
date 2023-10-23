FROM alpine:3.7

MAINTAINER Pavel Paulau <pavel@couchbase.com>

EXPOSE 8000

ENV CB_HOST ""
ENV CB_PASS ""
ENV CB_USER ""

COPY app app
COPY showfast /usr/local/bin/showfast

CMD ["showfast"]
