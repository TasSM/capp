package api

import (
	"fmt"
	"net/http"
)

func ServeRoutes(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HandleHealthcheck)
	mux.HandleFunc("/stats", HandleStats)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
	if err != nil {
		panic(err)
	}
}
