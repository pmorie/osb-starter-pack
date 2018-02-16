package user

import (
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"github.com/pmorie/osb-starter-pack/pkg/broker"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.
	return &BusinessLogic{
		async:     o.Async,
		instances: make(map[string]*exampleInstance, 10),
	}, nil
}

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Indiciates if the broker should handle the requests asynchronously.
	async bool
	// Synchronize go routines.
	sync.RWMutex
	// Add fields here! These fields are provided purely as an example
	instances map[string]*exampleInstance
}

var _ broker.BusinessLogic = &BusinessLogic{}

func (b *BusinessLogic) GetCatalog(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error) {
	// Your catalog business logic goes here
	response := &osb.CatalogResponse{}

	data := `
---
services:
- name: example-starter-pack-service
  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
  description: The example service from the osb starter pack!
  bindable: true
  plan_updateable: true
  metadata:
    displayName: "Example starter-pack service"
    imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
  plans:
  - name: default
    id: 86064792-7ea2-467b-af93-ac9694d96d5b
    description: The default plan for the starter pack example service
    free: true
    schemas:
      service_instance:
        create:
          "$schema": "http://json-schema.org/draft-04/schema"
          "type": "object"
          "title": "Parameters"
          "properties":
          - "name":
              "title": "Some Name"
              "type": "string"
              "maxLength": 63
              "default": "My Name"
          - "color":
              "title": "Color"
              "type": "string"
              "default": "Clear"
              "enum":
              - "Clear"
              - "Beige"
              - "Grey"
`

	err := yaml.Unmarshal([]byte(data), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (b *BusinessLogic) Provision(pr *osb.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*osb.ProvisionResponse, error) {
	// Your provision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := osb.ProvisionResponse{}

	exampleInstance := &exampleInstance{ID: pr.InstanceID, Params: pr.Parameters}
	b.instances[pr.InstanceID] = exampleInstance

	if pr.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, w http.ResponseWriter, r *http.Request) (*osb.DeprovisionResponse, error) {
	// Your deprovision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := osb.DeprovisionResponse{}

	delete(b.instances, request.InstanceID)

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error) {
	// Your last-operation business logic goes here

	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, w http.ResponseWriter, r *http.Request) (*osb.BindResponse, error) {
	// Your bind business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	instance, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	response := osb.BindResponse{
		Credentials: instance.Params,
	}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error) {
	// Your unbind business logic goes here
	return &osb.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, w http.ResponseWriter, r *http.Request) (*osb.UpdateInstanceResponse, error) {
	// Your logic for updating a service goes here.
	response := osb.UpdateInstanceResponse{}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	return nil
}

// example types

// exampleInstance is intended as an example of a type that holds information about a service instance
type exampleInstance struct {
	ID     string
	Params map[string]interface{}
}
