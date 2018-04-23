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

func TestUpdateInstance(t *testing.T) {
	cases := []struct {
		name         string
		validateFunc func(string) error
		updateFunc   func(req *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error)
		response     *osb.UpdateInstanceResponse
		err          error
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
			name: "update returns errors.New",
			updateFunc: func(req *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
				return nil, errors.New("oops")
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusInternalServerError,
				Description: strPtr("oops"),
			},
		},
		{
			name: "update returns osb.HTTPStatusCodeError",
			updateFunc: func(req *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
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
			name: "update returns sync",
			updateFunc: func(req *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
				return &broker.UpdateInstanceResponse{}, nil
			},
			response: &osb.UpdateInstanceResponse{},
		},
		{
			name: "update returns async",
			updateFunc: func(req *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
				return &broker.UpdateInstanceResponse{
					UpdateInstanceResponse: osb.UpdateInstanceResponse{
						Async: true,
					}}, nil
			},
			response: &osb.UpdateInstanceResponse{
				Async: true,
			},
		},
		{
			name: "update check originating origin identity is passed",
			updateFunc: func(req *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
				if req.OriginatingIdentity != nil {

					return &broker.UpdateInstanceResponse{
						UpdateInstanceResponse: osb.UpdateInstanceResponse{
							Async: true,
						}}, nil
				}
				return nil, fmt.Errorf("oops")
			},
			response: &osb.UpdateInstanceResponse{
				Async: true,
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

			api := &rest.APISurface{
				Broker: &FakeBroker{
					validateAPIVersion: validateFunc,
					update:             tc.updateFunc,
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
			o := osb.OriginatingIdentity{
				Platform: "kubernetes",
				Value:    `{"username":"test", "groups": [], "extra": {}}`,
			}

			actualResponse, err := client.UpdateInstance(&osb.UpdateInstanceRequest{
				InstanceID:          "12345",
				ServiceID:           "12345",
				AcceptsIncomplete:   true,
				OriginatingIdentity: &o,
			})
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

			if e, a := tc.response, actualResponse; !reflect.DeepEqual(e, a) {
				t.Errorf("Unexpected response\n\nExpected: %#+v\n\nGot: %#+v", e, a)
			}
		})
	}
}
