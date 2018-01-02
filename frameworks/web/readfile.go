package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleReadConfig - REST API endpoint for reading configuration file
func (rest *RestAPI) HandleReadFile(w http.ResponseWriter, r *http.Request) {
	application := mux.Vars(r)["application"]
	filename := mux.Vars(r)["filename"]

	reader, err := rest.Registry.Storage.Locate(application, filename)
	// Do we have an error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		// Its OK
		bytes, err := ioutil.ReadAll(reader)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("ERROR", err.Error())
			w.Write([]byte(err.Error()))
		}
	}
}
