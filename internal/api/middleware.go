package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TasSM/capp/internal/defs"
)

func logEnabledRequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logInfo := &defs.HTTPReqInfo{
			Method:    r.Method,
			URI:       r.URL.String(),
			Referer:   r.Header.Get("Referrer"),
			UserAgent: r.Header.Get("User-Agent"),
			IpAddr:    defs.RequestGetRemoteAddress(r),
		}
		out, err := json.Marshal(logInfo)
		if err != nil {
			panic(err)
		}
		log.Println(string(out))
		next.ServeHTTP(w, r)
	})
}
