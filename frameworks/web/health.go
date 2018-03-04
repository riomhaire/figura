package web

import (
	"net/http"
)

// HandleHealth - export basic health status for DevOps
func (r *RestAPI) HandleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("{ \"status\":\"up\"}"))

}
