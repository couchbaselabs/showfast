package main

var ds *dataStore

func main() {
	ds = newDataStore()
	ds.auth()

	httpEngine().Run(":8000")
}
