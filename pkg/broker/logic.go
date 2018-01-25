package broker

import (
	"net/http"

	"github.com/pmorie/go-open-service-broker-client/v2"
)

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
	GetCatalog(w http.ResponseWriter, r *http.Request) (*v2.CatalogResponse, error)
	// Provision encapsulates the business logic for a provision operation and
	// returns a v2.ProvisionResponse or an error. Provisioning creates a new
	// instance of a particular service.
	//
	// The parameters are:
	// - a v2.ProvisionRequest created from the original http request
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
	Provision(request v2.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*v2.ProvisionResponse, error)
	// Deprovision encapsulates the business logic for a deprovision operation
	// and returns a v2.DeprovisionResponse or an error. Deprovisioning deletes
	// an instance of a service and releases the resources associated with it.
	//
	// The parameters are:
	// - a v2.DeprovisionRequest created from the original http request
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
	Deprovision(w http.ResponseWriter, r *http.Request) (*v2.DeprovisionResponse, error)
	// LastOperation encapsulates the business logic for a last operation
	// request and returns a v2.LastOperationResponse or an error.
	// LastOperation is called when a platform checks the status of an ongoing
	// asynchronous operation on an instance of a service.
	//
	// The parameters are:
	// - a v2.LastOperationRequest created from the original http request
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
	LastOperation(w http.ResponseWriter, r *http.Request) (*v2.LastOperationResponse, error)
	// Bind encapsulates the business logic for a bind operation and returns a
	// v2.BindResponse or an error. Binding creates a new set of credentials for
	// a consumer to use an instance of a service. Not all services are
	// bindable; in order for a service to be bindable, either the service or
	// the current plan associated with the instance must declare itself to be
	// bindable.
	//
	// The parameters are:
	// - a v2.BindRequest created from the original http request
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
	Bind(w http.ResponseWriter, r *http.Request) (*v2.BindResponse, error)
	// Unbind encapsulates the business logic for an unbind operation and
	// returns a v2.UnbindResponse or an error. Unbind deletes a binding and the
	// resources associated with it.
	//
	// The parameters are:
	// - a v2.UnbindRequest created from the original http request
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
	Unbind(w http.ResponseWriter, r *http.Request) (*v2.UnbindResponse, error)
}

// Implementation provides an implementation of the BusinessLogic interface.
type Implementation struct {
	// You can add fields here or make your own type that implements BusinessLogic!
}

var _ BusinessLogic = &Implementation{}

func (b *Implementation) GetCatalog(w http.ResponseWriter, r *http.Request) (*v2.CatalogResponse, error) {
	// Your catalog business logic goes here
	return nil, nil
}

func (b *Implementation) Provision(pr v2.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*v2.ProvisionResponse, error) {
	// Your provision business logic goes here
	return nil, nil
}

func (b *Implementation) Deprovision(w http.ResponseWriter, r *http.Request) (*v2.DeprovisionResponse, error) {
	// Your deprovision business logic goes here
	return nil, nil
}

func (b *Implementation) LastOperation(w http.ResponseWriter, r *http.Request) (*v2.LastOperationResponse, error) {
	// Your last-operation business logic goes here
	return nil, nil
}

func (b *Implementation) Bind(w http.ResponseWriter, r *http.Request) (*v2.BindResponse, error) {
	// Your bind business logic goes here
	return nil, nil
}

func (b *Implementation) Unbind(w http.ResponseWriter, r *http.Request) (*v2.UnbindResponse, error) {
	// Your unbind business logic goes here
	return nil, nil
}

func (b *Implementation) ValidateBrokerAPIVersion(version string) error {
	return nil
}
