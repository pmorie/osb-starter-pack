package broker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pmorie/go-open-service-broker-client/v2"
)

func getBrokerAPIVersionFromRequest(r *http.Request) string {
	return r.Header.Get(v2.APIVersionHeader)
}

// writeResponse will serialize 'object' to the HTTP ResponseWriter
// using the 'code' as the HTTP status code
func writeResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

type e struct {
	Error string
}

func writeErrorResponse(w http.ResponseWriter, code int, err error) {
	type e struct {
		Error string
	}
	writeResponse(w, code, &e{
		Error: err.Error(),
	})
}

func unmarshalRequestBody(request *http.Request, obj interface{}) error {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}

	return nil
}
