package broker

import (

	"net/http"

	"github.com/golang/glog"
	"github.com/pmorie/osb-broker-lib/pkg/broker"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.

	// This is not ideal, create an environment variable for this path?
	dataverseInstances, err := FileToService(o.CatalogPath)

	if err != nil {
		return nil, err
	}

	dataverseMap := make(map[string]*dataverseInstance, len(dataverseInstances))

	for _, dataverse := range dataverseInstances {
		dataverseMap[dataverse.ServiceID] = dataverse
	}

	return &BusinessLogic{
		async:     o.Async,
		instances: make(map[string]*dataverseInstance, 10),
		// call dataverse server as little as possible
		dataverses: dataverseMap,
	}, nil
}

var _ broker.Interface = &BusinessLogic{}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	// Your catalog business logic goes here
	response := &broker.CatalogResponse{}

	// Create Service objects from dataverses
	services, err := DataverseToService(b.dataverses)

	if err != nil {
		return nil, err
	}

	osbResponse := &osb.CatalogResponse{
		Services : services,
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

	glog.Infof("provision request: %#+v", request)

	response := broker.ProvisionResponse{}

	dataverseInstance := &dataverseInstance{
		ID:        request.InstanceID,
		ServiceID: request.ServiceID,
		PlanID:    request.PlanID,
		ServerName: b.dataverses[request.ServiceID].ServerName,
		ServerUrl: b.dataverses[request.ServiceID].ServerUrl,
		Description: b.dataverses[request.ServiceID].Description,
		Params:    request.Parameters,
	}

	// Check to see if this is the same instance
	if i := b.instances[request.InstanceID]; i != nil {
		if i.Match(dataverseInstance) {
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

	// this should probably run asynchronously if possible
	if dataverseInstance.Params["credentials"] != nil && dataverseInstance.Params["credentials"].(string) != "" {
		// check that the token is valid, make a call to the Dataverse server
		succ, err := TestDataverseToken(dataverseInstance.ServerUrl, dataverseInstance.Params["credentials"].(string))

		if err != nil {
			return nil, err
		} else if succ != true {
			description := "Could not reach server"
			return nil, osb.HTTPStatusCodeError{
				StatusCode: http.StatusBadRequest,
				Description: &description,
			}
		}

	}
	
  
	b.instances[request.InstanceID] = dataverseInstance

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	glog.Infof("provision response: %#+v", response)

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

	glog.Infof("bind request: %#+v", request)

	instance, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	credentials := ""
	if instance.Params["credentials"] != nil {
			credentials = instance.Params["credentials"].(string)
	}

	response := broker.BindResponse{
		BindResponse: osb.BindResponse{
			Credentials: map[string]interface{}{
				"coordinates": instance.Description.Url,
				"credentials": credentials,
				},
		},

	}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	glog.Infof("bind response: %#+v", response)

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
