package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/SamiSousa/dataverse-broker/pkg/rest"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

func TestLastOperation(t *testing.T) {
	cases := []struct {
		name         string
		validateFunc func(string) error
		lastOpFunc   func(req *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error)
		response     *osb.LastOperationResponse
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
			name: "lastOperation returns errors.New",
			lastOpFunc: func(req *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error) {
				return nil, errors.New("oops")
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusInternalServerError,
				Description: strPtr("oops"),
			},
		},
		{
			name: "lastOperation returns osb.HTTPStatusCodeError",
			lastOpFunc: func(req *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error) {
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
			name: "OK",
			lastOpFunc: func(req *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error) {
				return &osb.LastOperationResponse{
					State: osb.StateSucceeded,
				}, nil
			},
			response: &osb.LastOperationResponse{
				State: osb.StateSucceeded,
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

			api := &rest.APISurface{
				BusinessLogic: &FakeBusinessLogic{
					validateAPIVersion: validateFunc,
					lastOperation:      tc.lastOpFunc,
				},
			}

			s := New(api)
			fs := httptest.NewServer(s.Router)
			defer fs.Close()

			config := defaultClientConfiguration()
			config.URL = fs.URL

			client, err := osb.NewClient(config)
			if err != nil {
				t.Error(err)
			}

			actualResponse, err := client.PollLastOperation(&osb.LastOperationRequest{
				InstanceID: "12345",
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
