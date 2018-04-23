package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnpackUpdateRequest(t *testing.T) {
	instanceID := "i1234"
	serviceID := "s1234"
	planID := "p1234"
	acceptsIncomplete := true

	fakeUpdateReq := createFakeUpdateRequest(serviceID, planID, acceptsIncomplete)
	unpackReq, err := unpackUpdateRequest(fakeUpdateReq, map[string]string{"instance_id": instanceID})
	if err != nil {
		t.Fatalf("Unpacking update request: %v", err)
	}

	if unpackReq.InstanceID != instanceID {
		t.Fatalf("InstanceID was unpacked unsuccessfully. Expecting %s got %s", instanceID, unpackReq.InstanceID)
	}

	if unpackReq.ServiceID != serviceID {
		t.Fatalf("PlanID was unpacked unsuccessfully. Expecting %s got %s", serviceID, unpackReq.ServiceID)
	}

	if *unpackReq.PlanID != planID {
		t.Fatalf("PlanID was unpacked unsuccessfully. Expecting %s got %s", planID, *unpackReq.PlanID)
	}

	if unpackReq.AcceptsIncomplete != acceptsIncomplete {
		t.Fatalf("AcceptsIncomplete was unpacked unsuccessfully. Expecting %t got %t", acceptsIncomplete, unpackReq.AcceptsIncomplete)
	}
}

func createFakeUpdateRequest(s, p string, a bool) *http.Request {
	data := fmt.Sprintf(`{
  "context": {
    "platform": "kubernetes",
    "some_field": "some-contextual-data"
  },
  "service_id": "%s",
  "plan_id": "%s",
  "parameters": {
    "parameter1": 1,
    "parameter2": "foo"
  },
  "previous_values": {
    "plan_id": "old-plan-id-here",
    "service_id": "service-id-here",
    "organization_id": "org-guid-here",
    "space_id": "space-guid-here"
  }
}`, s, p)

	r := bytes.NewBufferString(data)
	uri := fmt.Sprintf("/v2/service_instances/i1234?accepts_incomplete=%t", a)

	return httptest.NewRequest("PATCH", uri, r)
}
