package broker

import (
	"encoding/json"
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
