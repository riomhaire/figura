package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/riomhaire/figura/entities"
)

func (rest *RestAPI) HandleReadConfig(w http.ResponseWriter, r *http.Request) {
	application := mux.Vars(r)["application"]

	response := rest.Registry.ConfigurationReader.Lookup(r.Header.Get("Authorization"), application)
	rest.Registry.Logger.Log("TRACE", response.Message)
	code := HttpApplicationFailure
	data := []byte("Internal Error")

	switch response.ResultType {
	case entities.NoError:
		code = HttpSuccess
		data = response.Data
	case entities.UnknownApplication:
		code = HttpUnknown
		data = []byte(response.Message)
	case entities.AuthenticationError:
		code = HttpApplicationFailure
		data = []byte(response.Message)
	case entities.NotImplimentedError:
		code = HttpApplicationFailure
		data = []byte(response.Message)
	}

	// Write actual response to output stream
	w.WriteHeader(code)
	w.Write(data)
	if code != HttpSuccess {
		rest.Registry.Logger.Log("ERROR", fmt.Sprintf("App Error %v -> Http %v : %v", response.ResultType, code, string(data)))
	}
}
