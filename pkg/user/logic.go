package user

import (
	"net/http"

	"gopkg.in/yaml.v2"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"github.com/pmorie/go-open-service-broker-skeleton/pkg/broker"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.
	return &BusinessLogic{}, nil
}

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Add fields here!
}

var _ broker.BusinessLogic = &BusinessLogic{}

func (b *BusinessLogic) GetCatalog(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error) {
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

func (b *BusinessLogic) Provision(pr *osb.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*osb.ProvisionResponse, error) {
	// Your provision business logic goes here
	return nil, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, w http.ResponseWriter, r *http.Request) (*osb.DeprovisionResponse, error) {
	// Your deprovision business logic goes here
	return nil, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error) {
	// Your last-operation business logic goes here
	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, w http.ResponseWriter, r *http.Request) (*osb.BindResponse, error) {
	// Your bind business logic goes here
	return nil, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error) {
	// Your unbind business logic goes here
	return nil, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	return nil
}
