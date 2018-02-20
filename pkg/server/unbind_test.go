package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pmorie/osb-starter-pack/pkg/rest"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

func TestUnbind(t *testing.T) {
	cases := []struct {
		name         string
		validateFunc func(string) error
		unbindFunc   func(req *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error)
		response     *osb.UnbindResponse
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
			name: "unbind returns errors.New",
			unbindFunc: func(req *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error) {
				return nil, errors.New("oops")
			},
			err: osb.HTTPStatusCodeError{
				StatusCode:  http.StatusInternalServerError,
				Description: strPtr("oops"),
			},
		},
		{
			name: "unbind returns osb.HTTPStatusCodeError",
			unbindFunc: func(req *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error) {
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
			unbindFunc: func(req *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error) {
				return &osb.UnbindResponse{}, nil
			},
			response: &osb.UnbindResponse{},
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
					unbind:             tc.unbindFunc,
				},
			}

			s := New(api)
			fs := httptest.NewServer(s.Router)
			defer fs.Close()

			config := osb.DefaultClientConfiguration()
			config.URL = fs.URL

			client, err := osb.NewClient(config)
			if err != nil {
				t.Error(err)
			}

			actualResponse, err := client.Unbind(&osb.UnbindRequest{
				BindingID:  "12345",
				InstanceID: "12345",
				ServiceID:  "12345",
				PlanID:     "12345",
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
