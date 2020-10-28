package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TasSM/appCache/defs"
)

type CacheClientRouter struct {
	service defs.CacheClientService
}

func NewCacheClientRouter(svc defs.CacheClientService) defs.CacheClientRouter {
	return &CacheClientRouter{service: svc}
}

func (cr *CacheClientRouter) HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"Status":"OK"}`))
}

func (cr *CacheClientRouter) HandleStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	stats, e1 := cr.service.GetStatistics()
	if e1 != nil {
		log.Printf("ERROR - Retrieving statistics from redis failed: %v", e1)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Message":"Internal Server Error"}`))
		return
	}
	res, e2 := json.Marshal(stats)
	if e2 != nil {
		log.Printf("ERROR - Statistics retrieved are invalid: %v", e2)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Message":"Internal Server Error"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
