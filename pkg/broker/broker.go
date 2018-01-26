package broker

import (
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// APISurface is a type that describes a OSB REST API surface.  APISurface is
// responsible for decoding HTTP requests and transforming them into the request
// object for each operation and transforming responses and errors returned from
// the broker's internal business logic into the correct places in the HTTP
// response.
type APISurface struct {
	// Router is a mux.Router that registers the handlers for the different OSB
	// API operations.
	Router *mux.Router
	// BusinessLogic contains the business logic that provides the
	// implementation for the different OSB API operations.
	BusinessLogic BusinessLogic
}

const (
	instanceIDVarKey = "instance_id"
	bindingIDVarKey  = "binding_id"
)

// NewAPISurface returns a new, ready-to-go APISurface.
func NewAPISurface() *APISurface {
	router := mux.NewRouter()

	s := &APISurface{
		Router: router,
	}

	// TODO: update

	router.HandleFunc("/v2/catalog", s.GetCatalogHandler).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}/last_operation", s.LastOperationHandler).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}", s.ProvisionHandler).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}", s.DeprovisionHandler).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.BindHandler).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.UnbindHandler).Methods("DELETE")

	return s
}

// GetCatalogHandler is the mux handler that dispatches requests to get the
// broker's catalog to the broker's BusinessLogic.
func (s *APISurface) GetCatalogHandler(w http.ResponseWriter, r *http.Request) {
	version := getBrokerAPIVersionFromRequest(r)
	if err := s.BusinessLogic.ValidateBrokerAPIVersion(version); err != nil {
		writeError(w, err)
		return
	}

	response, err := s.BusinessLogic.GetCatalog(w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, response)
}

// ProvisionHandler is the mux handler that dispatches ProvisionRequests to the
// broker's BusinessLogic.
func (s *APISurface) ProvisionHandler(w http.ResponseWriter, r *http.Request) {
	version := getBrokerAPIVersionFromRequest(r)
	if err := s.BusinessLogic.ValidateBrokerAPIVersion(version); err != nil {
		writeError(w, err)
		return
	}

	request, err := unpackProvisionRequest(r)
	if err != nil {
		writeError(w, err)
		return
	}

	glog.Infof("Received ProvisionRequest for instanceID %q", request.InstanceID)

	response, err := s.BusinessLogic.Provision(request, w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	status := http.StatusOK
	if response.Async {
		status = http.StatusAccepted
	}

	writeResponse(w, status, response)
}

func unpackProvisionRequest(r *http.Request) (*osb.ProvisionRequest, error) {
	// unpacking an osb request from an http request involves:
	// - unmarshaling the request body
	// - getting IDs out of mux vars
	// - getting query parameters from request URL
	osbRequest := &osb.ProvisionRequest{}
	if err := unmarshalRequestBody(r, osbRequest); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	osbRequest.InstanceID = vars[instanceIDVarKey]

	asyncQueryParamVal := r.URL.Query().Get(asyncQueryParamKey)
	if strings.ToLower(asyncQueryParamVal) == "true" {
		osbRequest.AcceptsIncomplete = true
	}

	return osbRequest, nil
}

// DeprovisionHandler is the mux handler that dispatches deprovision requests to
// the broker's BusinessLogic.
func (s *APISurface) DeprovisionHandler(w http.ResponseWriter, r *http.Request) {
	version := getBrokerAPIVersionFromRequest(r)
	if err := s.BusinessLogic.ValidateBrokerAPIVersion(version); err != nil {
		writeError(w, err)
		return
	}

	// TODO
}

// LastOperationHandler is the mux handler that dispatches last-operation
// requests to the broker's BusinessLogic.
func (s *APISurface) LastOperationHandler(w http.ResponseWriter, r *http.Request) {
	version := getBrokerAPIVersionFromRequest(r)
	if err := s.BusinessLogic.ValidateBrokerAPIVersion(version); err != nil {
		writeError(w, err)
		return
	}

	// TODO
}

// BindHandler is the mux handler that dispatches bind requests to the broker's
// BusinessLogic.
func (s *APISurface) BindHandler(w http.ResponseWriter, r *http.Request) {
	version := getBrokerAPIVersionFromRequest(r)
	if err := s.BusinessLogic.ValidateBrokerAPIVersion(version); err != nil {
		writeError(w, err)
		return
	}

	// TODO
}

// UnbindHandler is the mux handler that dispatches unbind requests to the
// broker's BusinessLogic.
func (s *APISurface) UnbindHandler(w http.ResponseWriter, r *http.Request) {
	version := getBrokerAPIVersionFromRequest(r)
	if err := s.BusinessLogic.ValidateBrokerAPIVersion(version); err != nil {
		writeError(w, err)
		return
	}

	// TODO
}

// writeError accepts any error and writes it to the given ResponseWriter.
func writeError(w http.ResponseWriter, err error) {
	// TODO: make a little better :)

	if httpErr, ok := osb.IsHTTPError(err); ok {
		writeResponse(w, httpErr.StatusCode, err)
		return
	}

	writeResponse(w, http.StatusInternalServerError, err)
}
