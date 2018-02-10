package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// TODO: export in go-open-service-broker-client
const (
	asyncQueryParamKey = "accepts_incomplete"
)

func getBrokerAPIVersionFromRequest(r *http.Request) string {
	return r.Header.Get(osb.APIVersionHeader)
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

// writeError accepts any error and writes it to the given ResponseWriter along
// with a status code.
//
// If the error is an osb.HTTPStatusCodeError, the error's StatusCode field will
// be used and the response body will contain the error's Description and
// ErrorMessage fields.
//
// Otherwise, the given defaultStatusCode will be used, and the response body
// will have the result of calling the error's Error method set in the
// 'description' field.
//
// For more information about OSB errors, see:
//
// https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#service-broker-errors
func writeError(w http.ResponseWriter, err error, defaultStatusCode int) {
	if httpErr, ok := osb.IsHTTPError(err); ok {
		writeResponse(w, httpErr.StatusCode, err)
		return
	}

	writeErrorResponse(w, defaultStatusCode, err)
}

func writeErrorResponse(w http.ResponseWriter, code int, err error) {
	type e struct {
		Description string `json:"description"`
	}
	writeResponse(w, code, &e{
		Description: err.Error(),
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
