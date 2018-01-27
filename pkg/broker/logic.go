package broker

import (
	"net/http"

	"gopkg.in/yaml.v2"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// TODO: move BusinessLogic into its own package

// BusinessLogic contains the business logic for the broker's operations.
// BusinessLogic is the interface broker authors should implement and is
// embedded in an APISurface.
type BusinessLogic interface {
	// ValidateBrokerAPIVersion encapsulates the business logic of validating
	// the OSB API version sent to the broker with every request and returns
	// an error.
	//
	// For more information, see:
	//
	// https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#api-version-header
	ValidateBrokerAPIVersion(version string) error
	// GetCatalog encapsulates the business logic for returning the broker's
	// catalog of services. Brokers must tell platforms they're integrating with
	// which services they provide. GetCatalog is called when a platform makes
	// initial contact with the broker to find out about that broker's services.
	//
	// For more information, see:
	//
	// https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#catalog-management
	GetCatalog(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error)
	// Provision encapsulates the business logic for a provision operation and
	// returns a osb.ProvisionResponse or an error. Provisioning creates a new
	// instance of a particular service.
	//
	// The parameters are:
	// - a osb.ProvisionRequest created from the original http request
	// - a response writer, in case fine-grained control over the response is
	//   required
	// - the original http request, in case access is required (to get special
	//   request headers, for example)
	//
	// Implementers should return a ProvisionResponse for a successful operation
	// or an error. The APISurface handles translating ProvisionResponses or
	// errors into the correct form in the http response.
	//
	// For more information, see:
	//
	// https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#provisioning
	Provision(request *osb.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*osb.ProvisionResponse, error)
	// Deprovision encapsulates the business logic for a deprovision operation
	// and returns a osb.DeprovisionResponse or an error. Deprovisioning deletes
	// an instance of a service and releases the resources associated with it.
	//
	// The parameters are:
	// - a osb.DeprovisionRequest created from the original http request
	// - a response writer, in case fine-grained control over the response is
	//   required
	// - the original http request, in case access is required (to get special
	//   request headers, for example)
	//
	// Implementers should return a DeprovisionResponse for a successful
	// operation or an error. The APISurface handles translating
	// DeprovisionResponses or errors into the correct form in the http
	// response.
	//
	// For more information, see:
	//
	// https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#deprovisioning
	Deprovision(request *osb.DeprovisionRequest, w http.ResponseWriter, r *http.Request) (*osb.DeprovisionResponse, error)
	// LastOperation encapsulates the business logic for a last operation
	// request and returns a osb.LastOperationResponse or an error.
	// LastOperation is called when a platform checks the status of an ongoing
	// asynchronous operation on an instance of a service.
	//
	// The parameters are:
	// - a osb.LastOperationRequest created from the original http request
	// - a response writer, in case fine-grained control over the response is
	//   required
	// - the original http request, in case access is required (to get special
	//   request headers, for example)
	//
	// Implementers should return a LastOperationResponse for a successful
	// operation or an error. The APISurface handles translating
	// LastOperationResponses or errors into the correct form in the http
	// response.
	//
	// For more information, see:
	//
	// https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#polling-last-operation
	LastOperation(request *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error)
	// Bind encapsulates the business logic for a bind operation and returns a
	// osb.BindResponse or an error. Binding creates a new set of credentials for
	// a consumer to use an instance of a service. Not all services are
	// bindable; in order for a service to be bindable, either the service or
	// the current plan associated with the instance must declare itself to be
	// bindable.
	//
	// The parameters are:
	// - a osb.BindRequest created from the original http request
	// - a response writer, in case fine-grained control over the response is
	//   required
	// - the original http request, in case access is required (to get special
	//   request headers, for example)
	//
	// Implementers should return a BindResponse for a successful operation or
	// an error. The APISurface handles translating BindResponses or errors into
	// the correct form in the http response.
	//
	// For more information, see:
	//
	// https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#binding
	Bind(request *osb.BindRequest, w http.ResponseWriter, r *http.Request) (*osb.BindResponse, error)
	// Unbind encapsulates the business logic for an unbind operation and
	// returns a osb.UnbindResponse or an error. Unbind deletes a binding and the
	// resources associated with it.
	//
	// The parameters are:
	// - a osb.UnbindRequest created from the original http request
	// - a response writer, in case fine-grained control over the response is
	//   required
	// - the original http request, in case access is required (to get special
	//   request headers, for example)
	//
	// Implementers should return a UnbindResponse for a successful operation or
	// an error. The APISurface handles translating UnbindResponses or errors
	// into the correct form in the http response.
	//
	// For more information, see:
	//
	// https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#unbinding
	Unbind(request *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error)
}

// Implementation provides an implementation of the BusinessLogic interface.
type Implementation struct {
	// You can add fields here or make your own type that implements BusinessLogic!
}

var _ BusinessLogic = &Implementation{}

func (b *Implementation) GetCatalog(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error) {
	// Your catalog business logic goes here
	response := &osb.CatalogResponse{}

	data := `
---
services:
- name: skeleton-example-service
  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
  description: The example service from the broker skeleton!
  bindable: true
  plan_updateable: true
  plans:
  - name: default
    id: 86064792-7ea2-467b-af93-ac9694d96d5b
    description: The default plan for the skeleton example service
    free: true
`

	err := yaml.Unmarshal([]byte(data), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (b *Implementation) Provision(pr *osb.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*osb.ProvisionResponse, error) {
	// Your provision business logic goes here
	return nil, nil
}

func (b *Implementation) Deprovision(request *osb.DeprovisionRequest, w http.ResponseWriter, r *http.Request) (*osb.DeprovisionResponse, error) {
	// Your deprovision business logic goes here
	return nil, nil
}

func (b *Implementation) LastOperation(request *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error) {
	// Your last-operation business logic goes here
	return nil, nil
}

func (b *Implementation) Bind(request *osb.BindRequest, w http.ResponseWriter, r *http.Request) (*osb.BindResponse, error) {
	// Your bind business logic goes here
	return nil, nil
}

func (b *Implementation) Unbind(request *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error) {
	// Your unbind business logic goes here
	return nil, nil
}

func (b *Implementation) ValidateBrokerAPIVersion(version string) error {
	return nil
}
