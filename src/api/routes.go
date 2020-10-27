package api

import (
	"encoding/json"
	"net/http"

	"github.com/TasSM/appCache/defs"
)

type cacheClientRouter struct {
	service *defs.CacheClientService
}

func NewCacheClientRouter(svc *defs.CacheClientService) defs.CacheClientRouter {
	return &cacheClientRouter{service: svc}
}

func (cr *cacheClientRouter) HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	out := make(map[string]string)
	out["status"] = "OK"
	jsonout, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonout))
}

func (cr *cacheClientRouter) HandleStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// retrieve response from the stored service

	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte(jsonout))
}
