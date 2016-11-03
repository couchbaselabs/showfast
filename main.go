package main

var ds *dataStore

func main() {
	ds = newDataStore()

	httpEngine().Run(":8000")
}
