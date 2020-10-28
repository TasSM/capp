package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TasSM/appCache/defs"
)

func ServeRoutes(port string, service defs.CacheClientService) {
	router := NewCacheClientRouter(service)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", router.HandleHealthcheck)
	mux.HandleFunc("/stats", router.HandleStatistics)
	log.Printf("Starting HTTP server on port: %v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
	if err != nil {
		log.Printf("ERROR - Failed to serve http on port: %v", port)
		panic(err)
	}
}
