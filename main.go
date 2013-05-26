package main

import (
	"github.com/hoisie/web"
)

func index() string {
	return "Performance dashboard"
}

func main() {
	web.Get("/", index)
	web.Run("127.0.0.1:8080")
}
