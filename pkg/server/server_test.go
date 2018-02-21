package server

import (
	"net/http"

	"github.com/pmorie/osb-starter-pack/pkg/broker"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// TODO: is this more of an integration test?

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type FakeBusinessLogic struct {
	validateAPIVersion func(string) error
	getCatalog         func(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error)
	provision          func(pr *osb.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*osb.ProvisionResponse, error)
	deprovision        func(request *osb.DeprovisionRequest, w http.ResponseWriter, r *http.Request) (*osb.DeprovisionResponse, error)
	lastOperation      func(request *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error)
	bind               func(request *osb.BindRequest, w http.ResponseWriter, r *http.Request) (*osb.BindResponse, error)
	unbind             func(request *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error)
	update             func(request *osb.UpdateInstanceRequest, w http.ResponseWriter, r *http.Request) (*osb.UpdateInstanceResponse, error)
}

var _ broker.BusinessLogic = &FakeBusinessLogic{}

func (b *FakeBusinessLogic) GetCatalog(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error) {
	return b.getCatalog(w, r)
}

func (b *FakeBusinessLogic) Provision(pr *osb.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*osb.ProvisionResponse, error) {
	return b.provision(pr, w, r)
}

func (b *FakeBusinessLogic) Deprovision(request *osb.DeprovisionRequest, w http.ResponseWriter, r *http.Request) (*osb.DeprovisionResponse, error) {
	return b.deprovision(request, w, r)
}

func (b *FakeBusinessLogic) LastOperation(request *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error) {
	return b.lastOperation(request, w, r)
}

func (b *FakeBusinessLogic) Bind(request *osb.BindRequest, w http.ResponseWriter, r *http.Request) (*osb.BindResponse, error) {
	return b.bind(request, w, r)
}

func (b *FakeBusinessLogic) Unbind(request *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error) {
	return b.unbind(request, w, r)
}

func (b *FakeBusinessLogic) ValidateBrokerAPIVersion(version string) error {
	return b.validateAPIVersion(version)
}

func (b *FakeBusinessLogic) Update(request *osb.UpdateInstanceRequest, w http.ResponseWriter, r *http.Request) (*osb.UpdateInstanceResponse, error) {
	return b.update(request, w, r)
}

func defaultValidateFunc(_ string) error {
	return nil
}

func strPtr(s string) *string {
	return &s
}

func defaultClientConfiguration() *osb.ClientConfiguration {
	conf := osb.DefaultClientConfiguration()
	conf.Verbose = true

	return conf
}
