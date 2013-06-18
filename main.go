package main

import (
	"flag"
	"github.com/hoisie/mustache"
	"github.com/hoisie/web"
	"os"
	"strings"

	"github.com/pavel-paulau/showfast/datasource"
)

var pckg_dir string

func head() string {
	return mustache.RenderFile(pckg_dir + "templates/head.mustache")
}

func timeline(ctx *web.Context) []byte {
	metric := ctx.Params["metric"]
	return datasource.GetTimeline(metric)
}

func b2b(ctx *web.Context) []byte {
	metric := ctx.Params["metric"]
	builds := ctx.Params["builds"]
	return datasource.GetTimelineForBuilds(metric, builds)
}

func comparison(ctx *web.Context, val string) string {
	builds := strings.Split(val, "/")
	if len(val) == 0 || len(builds) > 2 {
		ctx.WriteHeader(400)
		return "Wrong number of builds to compare"
	}
	content := ""
	for _, metric := range datasource.GetMetricsForBuilds(builds) {
		metric["chart"] = mustache.RenderFile(
			pckg_dir+"templates/bars.mustache", metric)
		content += mustache.RenderFile(
			pckg_dir+"templates/metric.mustache", metric)
	}
	return mustache.RenderFile(pckg_dir+"templates/b2b.mustache",
		map[string]string{"head": head(), "content": content})
}

func home() string {
	content := ""
	for _, metric := range datasource.GetAllMetrics() {
		metric["chart"] = mustache.RenderFile(
			pckg_dir+"templates/columns.mustache", metric)
		content += mustache.RenderFile(
			pckg_dir+"templates/metric.mustache", metric)
	}
	return mustache.RenderFile(pckg_dir+"templates/home.mustache",
		map[string]string{"head": head(), "content": content})
}

func main() {
	address := flag.String("address", "127.0.0.1:8080", "Listen address")
	flag.Parse()

	pckg_dir = os.Getenv("GOPATH") + "/src/github.com/pavel-paulau/showfast/"
	web.Config.StaticDir = pckg_dir + "static"

	web.Get("/", home)
	web.Get("/timeline", timeline)
	web.Get("/comparison/(.*)", comparison)
	web.Get("/b2b", b2b)
	web.Run(*address)
}
