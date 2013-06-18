package main

import (
	"flag"
	"github.com/hoisie/mustache"
	"github.com/hoisie/web"
	"os"

	"github.com/pavel-paulau/showfast/datasource"
)

var pckg_dir string

func timeline(ctx *web.Context) []byte {
	metric := ctx.Params["metric"]
	return datasource.GetTimeline(metric)
}

func home() string {
	content := ""
	for _, benchmark := range datasource.GetBenchmarks() {
		content += mustache.RenderFile(pckg_dir+"templates/benchmark.mustache",
			benchmark)
	}
	return mustache.RenderFile(pckg_dir+"templates/home.mustache",
		map[string]string{"content": content})
}

func main() {
	address := flag.String("address", "127.0.0.1:8080", "Listen address")
	flag.Parse()

	pckg_dir = os.Getenv("GOPATH") + "/src/github.com/pavel-paulau/showfast/"
	web.Config.StaticDir = pckg_dir + "static"

	web.Get("/", home)
	web.Get("/timeline", timeline)
	web.Run(*address)
}
