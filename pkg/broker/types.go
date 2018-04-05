package broker

import (
	"sync"
)

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Indicates if the broker should handle the requests asynchronously.
	async bool
	// Synchronize go routines.
	sync.RWMutex
	// Add fields here! These fields are provided purely as an example
	instances map[string]*dataverseInstance
	// dataverse map dataverse_id to *dataverseInstances
	dataverses map[string]*dataverseInstance
}

// dataverseInstance holds information about a dataverse service instance
type dataverseInstance struct {
	ID        string
	ServiceID string
	PlanID    string
	Description *DataverseDescription
	ServerName string
	ServerUrl string
	Params    map[string]interface{} // Maybe add the DataverseDescription to this?
}

// Dataverse JSON Structs

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
	Message string `json:"message,omitempty"`
}
