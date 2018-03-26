package broker

import (
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/pmorie/osb-broker-lib/pkg/broker"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"reflect"
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
	// Indicates if the broker should handle the requests asynchronously.
	async bool
	// Synchronize go routines.
	sync.RWMutex
	// Add fields here! These fields are provided purely as an example
	instances map[string]*exampleInstance
}

var _ broker.Interface = &BusinessLogic{}

func truePtr() *bool {
	b := true
	return &b
}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	// Your catalog business logic goes here
	response := &broker.CatalogResponse{}
	osbResponse := &osb.CatalogResponse{
		Services: []osb.Service{
			{
				Name:          "example-starter-pack-service",
				ID:            "4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a",
				Description:   "The example service from the osb starter pack!",
				Bindable:      true,
				PlanUpdatable: truePtr(),
				Metadata: map[string]interface{}{
					"displayName": "Example starter pack service",
					"imageUrl":    "https://avatars2.githubusercontent.com/u/19862012?s=200&v=4",
				},
				Plans: []osb.Plan{
					{
						Name:        "default",
						ID:          "86064792-7ea2-467b-af93-ac9694d96d5b",
						Description: "The default plan for the starter pack example service",
						Free:        truePtr(),
						Schemas: &osb.Schemas{
							ServiceInstance: &osb.ServiceInstanceSchema{
								Create: &osb.InputParametersSchema{
									Parameters: map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"color": map[string]interface{}{
												"type":    "string",
												"default": "Clear",
												"enum": []string{
													"Clear",
													"Beige",
													"Grey",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	glog.Infof("catalog response: %#+v", osbResponse)

	response.CatalogResponse = *osbResponse

	return response, nil
}

func (b *BusinessLogic) Provision(request *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
	// Your provision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.ProvisionResponse{}

	exampleInstance := &exampleInstance{
		ID:        request.InstanceID,
		ServiceID: request.ServiceID,
		PlanID:    request.PlanID,
		Params:    request.Parameters,
	}

	// Check to see if this is the same instance
	if b.instances[request.InstanceID] != nil {
		i := b.instances[request.InstanceID]
		if i.Match(exampleInstance) {
			response.Exists = true
			return &response, nil
		} else {
			// Instance ID in use, this is a conflict.
			description := "InstanceID in use"
			return nil, osb.HTTPStatusCodeError{
				StatusCode: http.StatusConflict,
				Description: &description,
			}
		}
	}
	b.instances[request.InstanceID] = exampleInstance

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	// Your deprovision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.DeprovisionResponse{}

	delete(b.instances, request.InstanceID)

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
	// Your last-operation business logic goes here

	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
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

	response := broker.BindResponse{
		BindResponse: osb.BindResponse{
			Credentials: instance.Params,
		},
	}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
	// Your unbind business logic goes here
	return &broker.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
	// Your logic for updating a service goes here.
	response := broker.UpdateInstanceResponse{}
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
	ID        string
	ServiceID string
	PlanID    string
	Params    map[string]interface{}
}

func (i *exampleInstance) Match(other *exampleInstance) bool {
	return reflect.DeepEqual(i, other)
}
