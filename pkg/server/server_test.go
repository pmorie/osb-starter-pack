package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pmorie/osb-starter-pack/pkg/broker"
	"github.com/pmorie/osb-starter-pack/pkg/rest"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type FakeBusinessLogic struct {
	getCatalog    func(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error)
	provision     func(pr *osb.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*osb.ProvisionResponse, error)
	deprovision   func(request *osb.DeprovisionRequest, w http.ResponseWriter, r *http.Request) (*osb.DeprovisionResponse, error)
	lastOperation func(request *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error)
	bind          func(request *osb.BindRequest, w http.ResponseWriter, r *http.Request) (*osb.BindResponse, error)
	unbind        func(request *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error)
	update        func(request *osb.UpdateInstanceRequest, w http.ResponseWriter, r *http.Request) (*osb.UpdateInstanceResponse, error)
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
	return nil
}

func (b *FakeBusinessLogic) Update(request *osb.UpdateInstanceRequest, w http.ResponseWriter, r *http.Request) (*osb.UpdateInstanceResponse, error) {
	return b.update(request, w, r)
}

func TestGetCatalog(t *testing.T) {
	expectedResponse := &osb.CatalogResponse{Services: []osb.Service{
		{
			Name: "foo",
		},
	}}

	api := &rest.APISurface{
		BusinessLogic: &FakeBusinessLogic{
			getCatalog: func(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error) {
				return expectedResponse, nil
			},
		},
	}

	s := New(api)

	fs := httptest.NewServer(s.Router)
	defer fs.Close()

	config := osb.DefaultClientConfiguration()
	config.URL = fs.URL

	client, err := osb.NewClient(config)
	if err != nil {
		t.Error(err)
	}

	actualResponse, err := client.GetCatalog()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Unexpected response\n\nExpected: %#+v\n\nGot: %#+v", expectedResponse, actualResponse)
	}
}
