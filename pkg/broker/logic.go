package broker

import (
	"net/http"
	"sync"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"fmt"
	"strings"

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
		dataverse_server: "harvard",
		dataverse_url: "https://dataverse.harvard.edu",
		// call dataverse server as little as possible
		dataverses: GetDataverseServices("https://dataverse.harvard.edu", "harvard"),
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

	// name of the dataverse server
	dataverse_server string
	// url to dataverse server
	dataverse_url string
	// dataverse map dataverse_id to *DataverseDescription
	dataverses map[string]*DataverseDescription
}

var _ broker.Interface = &BusinessLogic{}

func DataverseToService(dataverses map[string]*DataverseDescription, server_name string) ([]osb.Service, error) {
	// Use DataverseDescription to populate osb.Service objects

	services := make([]osb.Service, len(dataverses))

	i := 0

	for _, dataverse := range dataverses {
		// use fields in DataverseDescription to populate osb.Service fields

		// check that each field has a value
		service_dashname := strings.ToLower(strings.Replace(dataverse.Name, " ", "-", -1))
		service_id := server_name + "-" + dataverse.Identifier
		service_description := dataverse.Description
		service_name := dataverse.Name
		service_image_url := dataverse.Image_url

		if service_description == ""{
			service_description = "A Dataverse service"
		}

		if service_image_url == ""{
			// default image for osb service
			service_image_url = "https://avatars2.githubusercontent.com/u/19862012?s=200&v=4"
		}

		services[i] = osb.Service{
				Name:          service_dashname,
				ID:            service_id,
				Description:   service_description, // comes out blank
				Bindable:      true,
				PlanUpdatable: truePtr(),
				Metadata: map[string]interface{}{
					"displayName": service_name,
					"imageUrl":    service_image_url,  // comes out blank
				},
				Plans: []osb.Plan{
					{
					Name:        "default",
					ID:          service_id + "-default",
					Description: "The default plan for " + service_name,
					Free:        truePtr(),
					Schemas: &osb.Schemas{
						ServiceInstance: &osb.ServiceInstanceSchema{
							Create: &osb.InputParametersSchema{
								Parameters: map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"credentials": map[string]interface{}{
											"type":    "string",
											"description": "API key to access restricted files and dataset on dataverse",
											"default": "",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		i += 1
	}

	return services, nil
}

// Add option to take in whitelist config
func GetDataverseServices(target_dataverse string, server_alias string) (map[string]*DataverseDescription) {

	dataverses, err := SearchForDataverses(&target_dataverse, 3)

	if err != nil{
		panic(err)
	}
	
	services := make(map[string]*DataverseDescription, len(dataverses))

	for _, dataverse := range dataverses {
		services[ server_alias + "-" +dataverse.Identifier] = dataverse
	}

	return services
}

func truePtr() *bool {
	b := true
	return &b
}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	// Your catalog business logic goes here
	response := &broker.CatalogResponse{}

	// Create Service objects from dataverses
	services, err :=  DataverseToService(b.dataverses, b.dataverse_server)

	if err != nil {
		panic(err)
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

	response := broker.ProvisionResponse{}

	exampleInstance := &exampleInstance{
		ID:        request.InstanceID,
		ServiceID: request.ServiceID,
		PlanID:    request.PlanID,
		Params:    request.Parameters,
	}

	// Check to see if this is the same instance
	if i := b.instances[request.InstanceID]; i != nil {
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

	// this should probably run asynchronously if possible
	if exampleInstance.Params["credentials"] != nil && exampleInstance.Params["credentials"].(string) != "" {
		// check that the token is valid, make a call to the Dataverse server
		// make a GET request
		
		resp, err := http.Get(b.dataverse_url + "/api/dataverses/:root?key=" + exampleInstance.Params["credentials"].(string))

		if err != nil{
			return nil, osb.HTTPStatusCodeError{
				StatusCode: http.StatusNotFound,
			}
		}

		// Must close response when finished
		defer resp.Body.Close()

		//convert resp into a DataverseResponse object
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil{
			return nil, osb.HTTPStatusCodeError{
				StatusCode: http.StatusNotFound,
			}
		}

		dataverseResp := DataverseResponseWrapper{}
		err = json.Unmarshal(body, &dataverseResp)

		// failed GET means token is invalid (what to do?)
		if err != nil || dataverseResp.Status != "OK"{
			description := "Bad api key '" + exampleInstance.Params["credentials"].(string) + "'"
			return nil, osb.HTTPStatusCodeError{
				StatusCode: http.StatusBadRequest,
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

	credentials := ""
	if instance.Params["credentials"] != nil {
			credentials = instance.Params["credentials"].(string)
	}

	response := broker.BindResponse{
		BindResponse: osb.BindResponse{
			// Get the service URL based on the serviceID (which is funny because they're the same thing right now...)
			Credentials: map[string]interface{}{
				"coordinates": b.dataverses[instance.ServiceID].Url,
				"credentials": credentials,
				},
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

// get all dataverses within a Dataverse server
// Takes a base Dataverse URL
// Returns a slice of string JSON objects, representing each dataverse
func SearchForDataverses(base *string, max_results_opt ... int) ([]*DataverseDescription, error) {
	// send a GET request to Dataverse url
	max_results := 0
	if len(max_results_opt) > 0{
		max_results = max_results_opt[0]
	}

	// Search API for dataverses
	search_uri := "/api/search"

	options := "?q=*&type=dataverse&start="

	// start with first search results, and only read back per_page number of dataverses per GET
	start := 0
	per_page := 10

	total_count := 0

	query_completed := false

	//slice to hold list of
	dataverses := make([]*DataverseDescription, 0)


	for query_completed == false {

		// make a GET request
		if max_results > 0 && max_results < start + per_page{
			// don't go over max_results
			per_page = max_results - start
		}
		resp, err := http.Get(*base + search_uri + options + strconv.Itoa(start) + "&per_page="+ strconv.Itoa(per_page))

		if err != nil{
			// exit on error
			fmt.Println("Error on http GET at address", *base + search_uri + options + strconv.Itoa(start) + "&per_page="+ strconv.Itoa(per_page))
			fmt.Println(err)
			panic("")
		}

		// Must close response when finished
		defer resp.Body.Close()

		//convert resp into a DataverseResponse object
		body, err := ioutil.ReadAll(resp.Body)

		response := DataverseResponseWrapper{}
		err = json.Unmarshal(body, &response)

		if err != nil{
			return dataverses, err
		}
		// check that Get was successful
		if response.Status != "OK"{
			fmt.Printf("Error in DataverseResponse status: %s\n", response.Status)
			panic("")
		}

		// obtain "total_count" for number of dataverses available at the server
		total_count = response.Data.Total_count

		// in case there are no results...
		if total_count == 0{
			panic("No results from GET query")
		}
		//otherwise, set condition to false if we've reached total_count
		if total_count == start + response.Data.Count_in_response{
			query_completed = true
		}
		// Reached max results
		if max_results > 0 && max_results <= start + response.Data.Count_in_response{
			query_completed = true
		}

		// iterate across each DataverseDescription
		for i := 0; i < response.Data.Count_in_response; i++ {
			// cast elements of list to DataverseDescription objects
			desc := DataverseDescription{}

			desc = response.Data.Items[i]

			// append DataverseDescription to dataverses slice
			dataverses = append(dataverses, &desc)
		}


		// update start value
		start += response.Data.Count_in_response
	}
	
	return dataverses, nil
	
}


// /Dataverse Structs

// object returned by checksum for datafiles
type DatafileChecksum struct{
	Type string `json:"type"`
	Value string `json:"value"`
}

// type for JSON portion describing a dataverse on Server
type DataverseDescription struct{
	// Wow, capitalization matters for structs in go?
	// Fields for dataverses
	Name string `json:"name"`
	Type string `json:"type"`
	Url string `json:"url"`
	Image_url string `json:"image_url,omitempty"`
	Identifier string `json:"identifier"`
	Description string `json:"description,omitempty"`
	Published_at string `json:"published_at"`

	// Fields for datasets
	Global_id string `json:"global_id,omitempty"`
	CitationHtml string `json:"citationHtml,omitempty"`
	Citation string `json:"citation,omitempty"`
	Authors []string `json:"authors,omitempty"`

	// Fields for datafiles
    File_id string `json:"file_id,omitempty"`
    File_type string `json:"file_type,omitempty"`
    File_content_type string `json:"file_content_type,omitempty"`
    Size_in_bytes int `json:"size_in_bytes,omitempty"`
    Md5 string `json:"md5,omitempty"`
    Dataset_citation string `json:"dataset_citation,omitempty"`
    Checksum DatafileChecksum `json:"checksum,omitempty"`

    // Fields for advanced search (to be added ...)
    

}

type DataverseResponse struct{
	// fields from response JSON object
	Count_in_response int `json:"count_in_response"`
	Items []DataverseDescription `json:"items"`
	Q string `json:"q"`
	Spelling_alternatives interface{} `json:"spelling_alternatives,omitempty"` // This is my lazy approach to a field I don't need
	Start int `json:"start"`
	Total_count int `json:"total_count"`

	// Only a partial list .. 

}

// type for JSON response from Dataverse API
type DataverseResponseWrapper struct{

	Data *DataverseResponse `json:"data"`
	Status string `json:"status"`
}
