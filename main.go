package main

import (
	"bytes"
	"flag"
	"log"

	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"
	"net/http"

	"github.com/rs/cors"

	"github.com/saroopmathur/rest-api/router"
)

// setupGlobalMiddleware will setup CORS
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://xpresstrust.alcyone.in", "http://18.222.133.118:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept", "Authorization"},
	})

	return c.Handler(handler)
}

// our main function
func main() {
	// Seed random generator
	randSeed()

	log.Printf("Listening on :8000\n")

	port := flag.String("p", "8000", "port to listen at")
	directory := flag.String("d", "./images", "folder containing images")
	flag.Parse()

	// Create router and start listen on port 8000
	router := router.NewRouter()

	// For seving staic files
	staticDir := "/images/"
	router.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir(*directory))))

	log.Fatal(http.ListenAndServe(":"+*port, setupGlobalMiddleware(router)))
}

// Use crypto rand to seed math rand
func randSeed() {
	var seed int64

	r := make([]byte, 8)
	_, _ = crand.Read(r)
	buf := bytes.NewBuffer(r)
	binary.Read(buf, binary.LittleEndian, &seed)
	mrand.Seed(seed)
}
