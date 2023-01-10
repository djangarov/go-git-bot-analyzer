package api

import (
	"fmt"
	"net/http"
	"time"
)

type Api struct {
	Port         int
	PrivateToken string
	GitlabHost   string
}

func (api *Api) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/scan-merge-request", api.handler(api.ScanMergeRequestHandler))

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", api.Port),
		ReadTimeout: 2 * time.Minute,
		Handler:     mux,
	}
	fmt.Printf("Server started listening on port %d... \n", api.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Error: %s \n", err)
	}
}

func (api *Api) handler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Start processing request: %s %s\n", r.URL.Path, r.Method)
		beginTime := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
		defer func() {
			fmt.Printf("Time to process the request: %d mS \n", time.Since(beginTime).Milliseconds())
		}()
	}
}
