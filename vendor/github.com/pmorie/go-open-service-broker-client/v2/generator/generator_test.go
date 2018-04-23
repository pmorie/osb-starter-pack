package generator

import (
	"encoding/json"
	"testing"

	"github.com/golang/glog"
)

func TestCreateGenerator(t *testing.T) {
	g := CreateGenerator(3, Parameters{
		Services: ServiceRanges{
			Plans:               5,
			Tags:                6,
			Metadata:            4,
			Requires:            2,
			Bindable:            10,
			BindingsRetrievable: 1,
		},
		Plans: PlanRanges{
			Metadata: 4,
			Bindable: 10,
			Free:     4,
		},
	})
	AssignPoolGoT(g)

	catalog, err := g.GetCatalog()
	if err != nil {
		t.Errorf("Got error, %v", err)
	}

	catalogBytes, err := json.MarshalIndent(catalog, "", "  ")

	catalogJson := string(catalogBytes)

	glog.Info(catalogJson)
}
