package broker

import (
	"net/http"
	"io/ioutil"
	"encoding/json"

	"strconv"
	"fmt"
	"strings"
	"regexp"

	"os"
	"path/filepath"

	"github.com/golang/glog"
	"reflect"

	osb "github.com/pmorie/go-open-service-broker-client/v2"

)

//Displays the individual metadata of the metadatablocks from a dataverse through injection of their ids
// if they are of type dataverse
func DataverseMeta(base string, id float64) {

	//Injecting string version of id into the search uri
	search_uri := "/api/dataverses/" + fmt.Sprint(id)

	//Variable that will store the json from the GET request
	var status map[string]interface{}

	//executing GET request
	resp, err := http.Get(base + search_uri)

	if err != nil {
		//exit on error
		fmt.Println("Error on http GET at address", base+search_uri)
		fmt.Println(err)
		panic("")
	}

	//Must close response when finished
	defer resp.Body.Close()

	//convert resp into a DataverseResponse object
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &status)
	if err != nil {
		fmt.Println(err)
	}

	//Checking GET response json
	if status["status"] == "ERROR" {

		//Skipping..
		fmt.Println("passing, not a dataverse.")

	} else {

		//Printing metadata
		fmt.Println(string(body))

	}

}

//Gathers the ids of the metadatablocks from a dataverse, calls on other function which displays the metadata
// of the individual metadatablocks (DataverseMeta), and returns an array of those ids if needed
func DataverseMetadataIds(base string) []float64{

	//This search uri finds the metadata of the dataverse and displays info objects with according id's
	search_uri := "/api/metadatablocks"

	//Variable that will store the json from the GET request
	var metadata map[string]interface{}

	//make a GET request
	resp, err := http.Get(base + search_uri)

	if err != nil {
		//exit on error
		fmt.Println("Error on http GET at address", base+search_uri)
		fmt.Println(err)
		panic("")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	//fmt.Println(string(body))

	err = json.Unmarshal(body, &metadata)
	if err != nil {
		fmt.Println(err)
	}

	//Creating return array object
	ids := make([]float64, 0)

	i := 0

	//Iterating through metadatablock
	for {

		//This recover function is to catch the panic when the index is out of range
		// indictating that the function has reached the end of the list of metadatablocks
		// and thus returns the array
		defer func() []float64 {
			if r := recover(); r != nil {
				//fmt.Println("Recovered", r)
				//fmt.Println(ids)
				return ids
			}
			return make([]float64, 0)
		}()

		//Retrieving ids from metadatablock
		test := metadata["data"].([]interface{})[i].(map[string]interface{})["id"]
		
		//Asserting test's type before injecting to separate function and return array
		string_test := test.(float64)
		
		//Calling function which displays the individual metadata of the metadatablocks
		DataverseMeta(base, string_test)
		
		//Appending to the return array
		ids = append(ids, string_test)
		
		//Incrementing iterator
		i++

	}


}

func DataverseToService(dataverses map[string]*dataverseInstance) ([]osb.Service, error) {
	// Use DataverseDescription to populate osb.Service objects

	services := make([]osb.Service, len(dataverses))

	i := 0
	reg, err := regexp.Compile("[^a-zA-Z0-9-.]+")
	if err != nil {
		return nil, err
	}

	for _, dataverse := range dataverses {
		// use fields in DataverseDescription to populate osb.Service fields

		// check that each field has a value
		// This name MUST be alphanumeric, dashes, and periods ONLY (no spaces)
		service_dashname := strings.ToLower(reg.ReplaceAllString(strings.Replace(dataverse.Description.Name, " ", "-", -1), ""))
		
		service_id := dataverse.ServiceID
		plan_id := dataverse.PlanID
		service_description := dataverse.Description.Description
		service_name := dataverse.Description.Name
		service_image_url := dataverse.Description.Image_url

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
					ID:          plan_id,
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
											"description": "API key to access restricted files and datasets on Dataverse",
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
func GetDataverseInstances(target_dataverse string, server_alias string) (map[string]*dataverseInstance) {

	dataverses, err := SearchForDataverses(&target_dataverse, 10)

	if err != nil{
		panic(err)
	}
	
	services := make(map[string]*dataverseInstance, len(dataverses))

	for _, dataverse := range dataverses {
		services[ server_alias + "-" +dataverse.Identifier] = &dataverseInstance{
			ID: server_alias + "-" +dataverse.Identifier,
			ServiceID: server_alias + "-" +dataverse.Identifier,
			PlanID: server_alias + "-" +dataverse.Identifier + "-default",
			ServerName: server_alias,
			ServerUrl: target_dataverse,
			Description: dataverse,
		}
	}

	return services
}

func FileToService(path string) ([]*dataverseInstance, error) {
	// take a file and turn it into dataverseInstances
	// each file stores a JSON object for a whitelisted dataverse service

	files, err := ioutil.ReadDir(path)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	instances := make([]*dataverseInstance, len(files))

	for i, f := range files {
		// read each file
		text, err := ioutil.ReadFile(filepath.Join(path,f.Name()))

		if err != nil{
			return nil, err
		}

		//Unmarshal string into dataverseInstance object
		dataverse := &dataverseInstance{}
		err = json.Unmarshal(text, dataverse)

		if err != nil {
			return nil, err
		}

		instances[i] = dataverse

	}

	return instances, nil

}

func ServiceToFile(instance *dataverseInstance, path string) (bool, error) {
	// take a service and store as JSON object in file
	// save as a file in path

	err := os.MkdirAll(path, os.ModePerm)

	if err != nil{
		return false, err
	}

	// get JSON from instance
	jsonInstance, err := json.Marshal(instance)

	if err != nil{
		return false, err
	}


	// write to file
	err = ioutil.WriteFile(filepath.Join(path,instance.ServiceID+".json"), jsonInstance, 0777)

	if err != nil {
		return false, err
	}

	return true, nil
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

func TestDataverseToken(serverUrl string, token string) (bool, error) {
	// ping the url, return bool for success or failure, and error code on fail
	resp, err := http.Get(serverUrl + "/api/dataverses/:root?key=" + token)

	if err != nil{
		return false, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	// Must close response when finished
	defer resp.Body.Close()

	//convert resp into a DataverseResponse object
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil{
		return false, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	dataverseResp := DataverseResponseWrapper{}
	err = json.Unmarshal(body, &dataverseResp)

	if err != nil || dataverseResp.Status != "OK"{
		return false, osb.HTTPStatusCodeError{
			StatusCode: http.StatusBadRequest,
			Description: &dataverseResp.Message,
		}
	}

	// reaching here means successful ping
	return true, nil
}

func truePtr() *bool {
	b := true
	return &b
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	return nil
}

func (i *dataverseInstance) Match(other *dataverseInstance) bool {
	return reflect.DeepEqual(i, other)
}
