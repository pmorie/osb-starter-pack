package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pmorie/osb-broker-lib/pkg/broker"
	"github.com/pmorie/osb-broker-lib/pkg/metrics"
	"github.com/pmorie/osb-broker-lib/pkg/rest"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
	prom "github.com/prometheus/client_golang/prometheus"
)

func TestDeprovision(t *testing.T) {
	cases := []struct {
		name            string
		validateFunc    func(string) error
		deprovisionFunc func(req *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error)
		response        *broker.DeprovisionResponse
		err             error
	}{
		{
			name: "version validation error",
			validateFunc: func(string) error {
				return errors.New("oops")
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusPreconditionFailed,
				Description: strPtr("oops"),
			},
		},
		{
			name: "returns errors.New",
			deprovisionFunc: func(req *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
				return nil, errors.New("oops")
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusInternalServerError,
				Description: strPtr("oops"),
			},
		},
		{
			name: "validate incoming parameters",
			deprovisionFunc: func(req *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
				if req.PlanID == "" {
					return nil, errors.New("deprovision request missing plan_id query parameter")
				}
				return &broker.DeprovisionResponse{}, nil
			},
			response: &broker.DeprovisionResponse{},
		},
		{
			name: "returns osb.HTTPStatusCodeError",
			deprovisionFunc: func(req *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
				return nil, osb.HTTPStatusCodeError{
					StatusCode:  http.StatusBadGateway,
					Description: strPtr("custom error"),
				}
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusBadGateway,
				Description: strPtr("custom error"),
			},
		},
		{
			name: "returns sync",
			deprovisionFunc: func(req *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
				return &broker.DeprovisionResponse{}, nil
			},
			response: &broker.DeprovisionResponse{},
		},
		{
			name: "returns async",
			deprovisionFunc: func(req *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
				return &broker.DeprovisionResponse{
					DeprovisionResponse: osb.DeprovisionResponse{
						Async: true,
					},
				}, nil
			},
			response: &broker.DeprovisionResponse{
				DeprovisionResponse: osb.DeprovisionResponse{
					Async: true,
				},
			},
		},
		{
			name: "check originating origin identity is passed",
			deprovisionFunc: func(req *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
				if req.OriginatingIdentity != nil {
					return &broker.DeprovisionResponse{
						DeprovisionResponse: osb.DeprovisionResponse{
							Async: true,
						},
					}, nil
				}
				return nil, fmt.Errorf("ops")
			},
			response: &broker.DeprovisionResponse{
				DeprovisionResponse: osb.DeprovisionResponse{
					Async: true,
				},
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			validateFunc := defaultValidateFunc
			if tc.validateFunc != nil {
				validateFunc = tc.validateFunc
			}

			// Prom. metrics
			reg := prom.NewRegistry()
			osbMetrics := metrics.New()
			reg.MustRegister(osbMetrics)

			request := &osb.DeprovisionRequest{
				InstanceID:          "12345",
				ServiceID:           "12345",
				PlanID:              "12345",
				AcceptsIncomplete:   true,
				OriginatingIdentity: originatingIdentity(),
			}

			// establish that the request we got was the request we sent
			deprovisionFunc := func(req *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
				if !reflect.DeepEqual(request, req) {
					t.Errorf("unexpected request; expected %v, got %v", request, req)
				}

				return tc.deprovisionFunc(req, c)
			}

			api := &rest.APISurface{
				Broker: &FakeBroker{
					validateAPIVersion: validateFunc,
					deprovision:        deprovisionFunc,
				},
				Metrics: osbMetrics,
			}

			s := New(api, reg)
			fs := httptest.NewServer(s.Router)
			defer fs.Close()

			config := defaultClientConfiguration()
			config.URL = fs.URL

			client, err := osb.NewClient(config)
			if err != nil {
				t.Error(err)
			}

			actualResponse, err := client.DeprovisionInstance(request)
			if err != nil {
				if tc.err != nil {
					if e, a := tc.err, err; !reflect.DeepEqual(e, a) {
						t.Errorf("Unexpected error; expected %v, got %v", e, a)
						return
					}
					return
				}
				t.Error(err)
				return
			}

			if e, a := &tc.response.DeprovisionResponse, actualResponse; !reflect.DeepEqual(e, a) {
				t.Errorf("Unexpected response\n\nExpected: %#+v\n\nGot: %#+v", e, a)
			}
		})
	}
}
