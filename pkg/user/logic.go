package user

import (
	"net/http"
	"sync"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"fmt"
	"strings"
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

// Conrad: what does this do
func DataverseToYAML() string {

	harvard := "https://dataverse.harvard.edu"
	target_dataverse := harvard //demo_dataverse

	dataverses, err := GetDataverses(&target_dataverse, 1)

	if err != nil{
		panic(err)
	}

	output := `
---
services:
` + DataverseToService(dataverses)

	return output

}

func (b *BusinessLogic) GetCatalog(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error) {

	response := &osb.CatalogResponse{}

	data := DataverseToYAML()

	err := yaml.Unmarshal([]byte(data), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// func (b *BusinessLogic) GetCatalog(w http.ResponseWriter, r *http.Request) (*osb.CatalogResponse, error) {
// 	// Your catalog business logic goes here
// 	response := &osb.CatalogResponse{}

// 	data := `
// ---
// services:
// - name: example-starter-pack-service
//   id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a
//   description: The example service from the osb starter pack!
//   bindable: true
//   plan_updateable: true
//   metadata:
//     displayName: "Example starter-pack service"
//     imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
//   plans:
//   - name: default
//     id: 86064792-7ea2-467b-af93-ac9694d96d5b
//     description: The default plan for the starter pack example service
//     free: true
//     schemas:
//       service_instance:
//         create:
//           "$schema": "http://json-schema.org/draft-04/schema"
//           "type": "object"
//           "title": "Parameters"
//           "properties":
//           - "name":
//               "title": "Some Name"
//               "type": "string"
//               "maxLength": 63
//               "default": "My Name"
//           - "color":
//               "title": "Color"
//               "type": "string"
//               "default": "Clear"
//               "enum":
//               - "Clear"
//               - "Beige"
//               - "Grey"
// - name: example-starter-pack-service2
//   id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246b
//   description: Second service
//   bindable: true
//   plan_updateable: true
//   metadata:
//     displayName: "Example second service"
//     imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
//   plans:
//   - name: default
//     id: 86064792-7ea2-467b-af93-ac9694d96d5c
//     description: The default plan for the second starter pack example service
//     free: true
//     schemas:
//       service_instance:
//         create:
//           "$schema": "http://json-schema.org/draft-04/schema"
//           "type": "object"
//           "title": "Parameters"
//           "properties":
//           - "name":
//               "title": "Some Name"
//               "type": "string"
//               "maxLength": 63
//               "default": "My Name"
//           - "color":
//               "title": "Color"
//               "type": "string"
//               "default": "Clear"
//               "enum":
//               - "Clear"
//               - "Beige"
//               - "Grey"
// `

// 	err := yaml.Unmarshal([]byte(data), &response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

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

// get all dataverses within a Dataverse server
// Takes a base Dataverse URL
// Returns a slice of string JSON objects, representing each dataverse
func GetDataverses(base *string, max_results_opt ... int) ([]*DataverseDescription, error) {
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

	/*
	for i := 0; i < len(dataverses); i++{
		// print for good measure
		byt, err := json.MarshalIndent(dataverses[i], "", "\t")
		if err != nil{
			panic("Error in  json Marshaling")
			
		}
		fmt.Printf("%s\n", string(byt))
	}
	*/
	
	
	return dataverses, nil
	
}

func DataverseToService(dataverses []*DataverseDescription) string {

	var services string

	for i := 0; i < len(dataverses); i++ {

		// // Descriptions cannot be empty, but the Dataverse object
		// // does not require their own descriptions to be non-empty
		// // So, need to check
		// var desc string
		// if dataverses[i].Description == "" {
		// 	desc = "none"
		// } else {
		// 	desc = dataverses[i].Description
		// }

		services = services + fmt.Sprintf(
`- name: %s
  id: %s
  description: none
  bindable: true
  plan_updateable: true
  metadata:
    displayName: "%s"
    imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
  plans:
  - name: default
    id: %s
    description: The default plan for the second starter pack example service
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
`, strings.Replace(dataverses[i].Name, " ", "", -1), ReturnGUID(), dataverses[i].Identifier, ReturnGUID())

	}

	return services
}

// MAY NOT BE COMPLIANT WITH GUID GEN
func ReturnGUID() string {

	// u := make([]byte, 16)
	// _, err := rand.Read(u)

	// if err != nil {
 //    	return ""
	// }

	// u[8] = (u[8] | 0x80) & 0xBF // what does this do?
	// u[6] = (u[6] | 0x40) & 0x4F // what does this do?

	return "4f6e6cf6-ffdd-425f-a2c7-3c9258ad246b"
}

// /Dataverse Structs

// type for JSON portion describing a dataverse on Server
type DataverseDescription struct{
	// Wow, capitalization matters for structs in go?
	Description string `json:"description"`
	Identifier string `json:"identifier"`
	Name string `json:"name"`
	Published_at string `json:"published_at"`
	Type string `json:"type"`
	Url string `json:"url"`

	// the ones below might only populate for data files, maybe omit empty"
    Image_url string `json:"image_url,omitempty"`
    File_id string `json:"file_id,omitempty"`
    File_type string `json:"file_type,omitempty"`
    File_content_type string `json:"file_content_type,omitempty"`
    Size_in_bytes int `json:"size_in_bytes,omitempty"`
    Md5 string `json:"md5,omitempty"`
    Dataset_citation string `json:"dataset_citation,omitempty"`

}

type DataverseResponse struct{
	// fields from response JSON object
	Count_in_response int `json:"count_in_response"`
	Items []DataverseDescription `json:"items"`
	Q string `json:"q"`
	Spelling_alternatives interface{} `json:"spelling_alternatives,omitempty"`
	Start int `json:"start"`
	Total_count int `json:"total_count"`
}

// type for JSON response from Dataverse API
type DataverseResponseWrapper struct{

	Data *DataverseResponse `json:"data"`
	Status string `json:"status"`
}