package generator

import (
	"encoding/json"
	"testing"

	"github.com/golang/glog"
)

func TestGetCatalog(t *testing.T) {
	g := &Generator{
		Services: []Service{
			{
				Plans: []Plan{
					{
						FromPool: Pull{
							Tags:     3,
							Metadata: 4,
							Free:     1,
						},
					},
					{
						FromPool: Pull{
							Tags:     3,
							Metadata: 4,
						},
					},
				},
				FromPool: Pull{
					Tags:                3,
					Metadata:            4,
					BindingsRetrievable: 1,
					Bindable:            1,
					Requires:            2,
				},
			},
		},
	}
	AssignPoolGoT(g)

	catalog, err := g.GetCatalog()
	if err != nil {
		t.Errorf("Got error, %v", err)
	}

	catalogBytes, err := json.MarshalIndent(catalog, "", "  ")

	catalogJson := string(catalogBytes)

	glog.Info(catalogJson)
}

func TestGetPlans(t *testing.T) {

	g := Generator{
		PlanPool: []string{"AAA", "BBB", "CCC", "DDD", "EEE"},
	}
	glog.Info(g.planNames(1, 5))
	glog.Info(g.planNames(2, 5))
}
