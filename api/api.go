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
	mux.HandleFunc("/api/scan-merge-request", api.ScanMergeRequestHandler)

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
