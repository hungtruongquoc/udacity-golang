package main

import (
	"fmt"
	"net/http"
)

// "Cities" is a slice of strings
var cities = []string{"Tokyo", "Delhi", "Shanghai", "Sao Paulo", "Mexico City"}

// The "index" handler function. Note the two parameters:
// 1. The ResponseWriter interface is used by the "index" handler to construct an HTTP response
// 2. The Request is an HTTP request received by the server
func index(w http.ResponseWriter, r *http.Request) {
	// "Fprintf" formats and writes to w
	fmt.Fprintf(w, "Main page\n")
}

func cityList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of most populous cities:\n")

	// Using a range to iterate through the "cities" slice and writing the name of the city
	for _, city := range cities {
		fmt.Fprintf(w, "%s\n", city)
	}
}

func main() {
	// Registering the "index" handler function to the "/" route (e.g., "https://example.com/",  "http://localhost:3000/", etc.)
	http.HandleFunc("/", index)
	// Registering the "cityList" handler function to the "/citylist" route (e.g., "https://example.com/citylist",  "http://localhost:3000/citylist", etc.)
	http.HandleFunc("/citylist", cityList)

	fmt.Println("Server is starting on port 3000...")
	// Instructing this HTTP server to listen for incoming requests on port 3000
	http.ListenAndServe(":3000", nil)
}
