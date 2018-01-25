package broker

import (
	"net/http"

	_ "github.com/golang/glog"

	"github.com/gorilla/mux"
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

// NewAPISurface returns a new, ready-to-go APISurface.
func NewAPISurface() *APISurface {
	router := mux.NewRouter()

	s := &APISurface{
		Router: router,
	}

	router.HandleFunc("/v2/catalog", s.GetCatalogHandler).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}/last_operation", s.LastOperationHandler).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}", s.ProvisionHandler).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}", s.DeprovisionHandler).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.BindHandler).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.UnbindHandler).Methods("DELETE")

	return s
}

// GetCatalogHandler is a Handler for catalog requests to the broker.
// GetCatalogHandler validates the request, delegates to the broker
// BusinessLogic to get the CatalogResponse, and handles writing the HTTP
// response.
func (s *APISurface) GetCatalogHandler(w http.ResponseWriter, r *http.Request) {
	version := getBrokerAPIVersionFromRequest(r)
	if err := s.BusinessLogic.ValidateBrokerAPIVersion(version); err != nil {
		// write the error back
	}

	response, err := s.BusinessLogic.GetCatalog(w, r)
	if err != nil {
		// check for client http error and directly serialize
		// otherwise set status to 500 and return the error message in the body.
	}

	writeResponse(w, http.StatusOK, response)
}

// ProvisionHandler is a Handler for provision requests to the broker.
func (s *APISurface) ProvisionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO

	// this layer should deserialize body/unpack request params into a
	// ProvisionRequest to pass to the business logic layer

	// this layer also handles serializing the response or error returned
	// from the business logic layer
}

// DeprovisionHandler is a Handler for deprovision requests to the broker.
func (s *APISurface) DeprovisionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// LastOperationHandler is a Handler for last operation requests to the broker.
func (s *APISurface) LastOperationHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// BindHandler is a Handler for bind requests to the broker.
func (s *APISurface) BindHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// UnbindHandler is a Handler for unbind requests to the broker.
func (s *APISurface) UnbindHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
