package main

import (
	"github.com/hoisie/mustache"
	"github.com/hoisie/web"
)

func home() string {
	var num_benchmarks = 0
	var benchmarks = ""

	for i := 0; i < num_benchmarks; i++ {
		var benchmark = map[string]string{}
		benchmarks += mustache.RenderFile("templates/benchmark.mustache", benchmark)
	}

	return mustache.RenderFile("templates/home.mustache", map[string]string{"benchmarks": benchmarks})
}

func main() {
	web.Get("/", home)
	web.Run("127.0.0.1:8080")
}
