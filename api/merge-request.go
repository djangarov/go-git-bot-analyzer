package api

import (
	"fmt"
	"net/http"
)

func (api *Api) ScanMergeRequestHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("ScanMergeRequestHandler")
}
