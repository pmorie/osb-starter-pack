package server

import (
	"errors"
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

func TestProvision(t *testing.T) {
	cases := []struct {
		name          string
		validateFunc  func(string) error
		provisionFunc func(req *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error)
		response      *broker.ProvisionResponse
		err           error
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
			provisionFunc: func(req *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
				return nil, errors.New("oops")
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusInternalServerError,
				Description: strPtr("oops"),
			},
		},
		{
			name: "returns osb.HTTPStatusCodeError",
			provisionFunc: func(req *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
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
			provisionFunc: func(req *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
				return &broker.ProvisionResponse{
					ProvisionResponse: osb.ProvisionResponse{
						DashboardURL: strPtr("my.service.to/12345"),
					}}, nil
			},
			response: &broker.ProvisionResponse{
				ProvisionResponse: osb.ProvisionResponse{
					DashboardURL: strPtr("my.service.to/12345"),
				}},
		},
		{
			name: "returns async",
			provisionFunc: func(req *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
				return &broker.ProvisionResponse{
					ProvisionResponse: osb.ProvisionResponse{
						Async:        true,
						DashboardURL: strPtr("my.service.to/12345"),
					}}, nil
			},
			response: &broker.ProvisionResponse{
				ProvisionResponse: osb.ProvisionResponse{
					Async:        true,
					DashboardURL: strPtr("my.service.to/12345"),
				}},
		},
		{
			name: "check originating origin identity is passed",
			provisionFunc: func(req *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
				if req.OriginatingIdentity != nil {

					return &broker.ProvisionResponse{
						ProvisionResponse: osb.ProvisionResponse{
							Async:        true,
							DashboardURL: strPtr("my.service.to/12345"),
						}}, nil
				}
				return nil, errors.New("oops")
			},
			response: &broker.ProvisionResponse{
				ProvisionResponse: osb.ProvisionResponse{
					Async:        true,
					DashboardURL: strPtr("my.service.to/12345"),
				}},
		},
		{
			name: "returns already completed",
			provisionFunc: func(req *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
				return &broker.ProvisionResponse{
					Exists: true,
					ProvisionResponse: osb.ProvisionResponse{
						DashboardURL: strPtr("my.service.to/12345"),
					}}, nil
			},
			response: &broker.ProvisionResponse{
				ProvisionResponse: osb.ProvisionResponse{
					DashboardURL: strPtr("my.service.to/12345"),
				}},
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

			request := &osb.ProvisionRequest{
				InstanceID:          "12345",
				ServiceID:           "12345",
				PlanID:              "12345",
				OrganizationGUID:    "12345",
				SpaceGUID:           "12345",
				AcceptsIncomplete:   true,
				OriginatingIdentity: originatingIdentity(),
			}

			// establish that the request we got was the request we sent
			provisionFunc := func(req *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
				if !reflect.DeepEqual(request, req) {
					t.Errorf("unexpected request; expected %v, got %v", request, req)
				}

				return tc.provisionFunc(req, c)
			}

			api := &rest.APISurface{
				Broker: &FakeBroker{
					validateAPIVersion: validateFunc,
					provision:          provisionFunc,
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

			actualResponse, err := client.ProvisionInstance(request)
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

			if e, a := &tc.response.ProvisionResponse, actualResponse; !reflect.DeepEqual(e, a) {
				t.Errorf("Unexpected response\n\nExpected: %#+v\n\nGot: %#+v", e, a)
			}
		})
	}
}
