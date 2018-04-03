// Test Provision and Bind
package broker

import(
	"testing"

	"github.com/pmorie/osb-broker-lib/pkg/broker"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	logic "github.com/SamiSousa/dataverse-broker/pkg/broker"
)

func TestBrokerLogic(t *testing.T){
	// create a BusinessLogic struct instance (tests dataverse functions)
	businessLogic, errCreate := logic.NewBusinessLogic(logic.Options{CatalogPath: "", Async: false})

	if errCreate != nil{
		t.Errorf("Error on BusinessLogic creation: %#+v", errCreate)
	}

	// Run Provision on a couple of test cases:
	// credentials blank
	_, errProvisionBlank := businessLogic.Provision(
		&osb.ProvisionRequest{
			InstanceID:	"harvard-ephelps",
			AcceptsIncomplete:	false,
			ServiceID:	"harvard-ephelps",
			PlanID:	"harvard-ephelps-default",
			OrganizationGUID:	"bdc",
			SpaceGUID:	"bdc",
			Parameters:	map[string]interface{}{},
		}, 
		// empty because we don't use it
		&broker.RequestContext{})

	if errProvisionBlank != nil {
		t.Errorf("Error on Provision with blank token: %#+v", errProvisionBlank)
	}

	// improper credentials
	_, errProvisionImproper := businessLogic.Provision(
		&osb.ProvisionRequest{
			InstanceID:	"harvard-ephelps",
			AcceptsIncomplete:	false,
			ServiceID:	"harvard-ephelps",
			PlanID:	"harvard-ephelps-default",
			OrganizationGUID:	"bdc",
			SpaceGUID:	"bdc",
			Parameters:	map[string]interface{}{
					"credentials":"not-real-token",
				},
		}, 
		// empty because we don't use it
		&broker.RequestContext{})

	// we want an error here
	if errProvisionImproper == nil {
		t.Errorf("Error on Provision with invalid token: no error returned")
	}

	/*
	// proper credentials
	_, err = businessLogic.Provision(
		&osb.ProvisionRequest{
			InstanceID:	"harvard-ephelps",
			AcceptsIncomplete:	false,
			ServiceID:	"harvard-ephelps",
			PlanID:	"harvard-ephelps-default",
			OrganizationGUID:	"bdc",
			SpaceGUID:	"bdc",
			Parameters:	map[string]interface{}{
					"credentials":"totally-real-token", // replace this with real token in secure way
				},
		}, 
		// empty because we don't use it
		&broker.RequestContext{})

	// this should succeed, using token from config
	if err != nil {
		t.Errorf("Error on Provision with valid token: %#+v", err)
	}
	*/

	// Run Bind on a couple of test cases
	// credentials blank
	_, errBindBlank := businessLogic.Bind(
		&osb.BindRequest{
			BindingID:	"harvard-ephelps",
			InstanceID:	"harvard-ephelps",
			AcceptsIncomplete:	false,
			ServiceID:	"harvard-ephelps",
			PlanID:	"harvard-ephelps-default",
			Parameters:	map[string]interface{}{},
		}, 
		&broker.RequestContext{})

	if errBindBlank != nil{
		t.Errorf("Error on Bind with no token: %#+v", errBindBlank)
	}

	/*
	// credentials nonblank
	bindResultProper, err := businessLogic.Bind(
		&osb.BindRequest{
			BindingID:	"harvard-ephelps",
			InstanceID:	"harvard-ephelps",
			AcceptsIncomplete:	false,
			ServiceID:	"harvard-ephelps",
			PlanID:	"harvard-ephelps-default",
			Parameters:	map[string]interface{}{
				"credentials": "totally-real-token",
			},
		}, 
		&broker.RequestContext{})

	if err != nil{
		t.Errorf("Error on Bind with valid token: %#+v", err)
	}

	if bindResultProper.BindResponse.Credentials["credentials"] != "totally-real-token" || bindResultProper.BindResponse.Credentials["coordinates"] == nil{
		t.Errorf("Error on Bind: credentials and coordinates not passed properly")
	}
	*/

}