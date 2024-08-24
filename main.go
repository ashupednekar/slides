
package main

import (
	"log"
	"net/http"

	"golang.org/x/tools/present"
)

func main() {
	// Define the content path (where your slide files are located)
	contentPath := "."

	// Set the default template to render the slides
	present.Init()

	// Start the HTTP server to serve the presentation
	http.Handle("/", http.FileServer(http.Dir(contentPath)))
	err := http.ListenAndServe(":3999", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
