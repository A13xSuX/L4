package main

import (
	"flag"
	"log"
	"metrics/internal/httpserver"
	"net/http"
	"runtime/debug"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	gcPercent := flag.Int("gc-percent", 100, "gc target percentage")

	flag.Parse()

	oldGCPercent := debug.SetGCPercent(*gcPercent)
	log.Printf("GC percent set to %d (previous percent %d)", *gcPercent, oldGCPercent)

	mux := httpserver.NewRouter()

	log.Printf("Server start on %s", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
