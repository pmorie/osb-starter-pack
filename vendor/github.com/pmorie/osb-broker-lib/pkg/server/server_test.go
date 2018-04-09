package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"github.com/pmorie/osb-broker-lib/pkg/broker"
	"github.com/pmorie/osb-broker-lib/pkg/metrics"
	"github.com/pmorie/osb-broker-lib/pkg/rest"
)

// TODO: is this more of an integration test?

// FakeBroker provides an implementation of the broker.Interface.
type FakeBroker struct {
	validateAPIVersion func(string) error
	getCatalog         func(c *broker.RequestContext) (*broker.CatalogResponse, error)
	provision          func(pr *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error)
	deprovision        func(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error)
	lastOperation      func(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error)
	bind               func(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error)
	unbind             func(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error)
	update             func(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error)
}

var _ broker.Interface = &FakeBroker{}

func (b *FakeBroker) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	return b.getCatalog(c)
}

func (b *FakeBroker) Provision(pr *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
	return b.provision(pr, c)
}

func (b *FakeBroker) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	return b.deprovision(request, c)
}

func (b *FakeBroker) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
	return b.lastOperation(request, c)
}

func (b *FakeBroker) Bind(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
	return b.bind(request, c)
}

func (b *FakeBroker) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
	return b.unbind(request, c)
}

func (b *FakeBroker) ValidateBrokerAPIVersion(version string) error {
	return b.validateAPIVersion(version)
}

func (b *FakeBroker) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
	return b.update(request, c)
}

func defaultValidateFunc(_ string) error {
	return nil
}

func strPtr(s string) *string {
	return &s
}

func defaultClientConfiguration() *osb.ClientConfiguration {
	conf := osb.DefaultClientConfiguration()
	conf.Verbose = true

	return conf
}

func originatingIdentity() *osb.OriginatingIdentity {
	return &osb.OriginatingIdentity{
		Platform: "kubernetes",
		Value:    `{"username":"test", "groups": [], "extra": {}}`,
	}
}

func TestNewHTTPHandler(t *testing.T) {
	type args struct {
		broker        broker.Interface
		servicePath   string
		serviceMethod string
		request       []byte
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "test APISurface GetCatalog(...)",
			args: args{
				broker: &fakeBroker{
					validateBrokerAPIVersion: func(version string) error { return nil },
					getCatalog: func(c *broker.RequestContext) (*broker.CatalogResponse, error) {
						return &broker.CatalogResponse{}, nil
					},
				},
				servicePath:   "/v2/catalog",
				serviceMethod: http.MethodGet,
			},
			wantStatusCode: 200,
		},
		{
			name: "test APISurface LastOperation(...)",
			args: args{
				broker: &fakeBroker{
					validateBrokerAPIVersion: func(version string) error { return nil },
					lastOperation: func(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
						return &broker.LastOperationResponse{}, nil
					},
				},
				servicePath:   "/v2/service_instances/foo/last_operation",
				serviceMethod: http.MethodGet,
			},
			wantStatusCode: 200,
		},
		{
			name: "test APISurface Provision(...)",
			args: args{
				broker: &fakeBroker{
					validateBrokerAPIVersion: func(version string) error { return nil },
					provision: func(request *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
						return &broker.ProvisionResponse{}, nil
					},
				},
				servicePath:   "/v2/service_instances/foo",
				serviceMethod: http.MethodPut,
				request:       []byte("{}"),
			},
			wantStatusCode: 201,
		},
		{
			name: "test APISurface Deprovision(...)",
			args: args{
				broker: &fakeBroker{
					validateBrokerAPIVersion: func(version string) error { return nil },
					deprovision: func(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
						return &broker.DeprovisionResponse{}, nil
					},
				},
				servicePath:   "/v2/service_instances/foo",
				serviceMethod: http.MethodDelete,
				request:       []byte("{}"),
			},
			wantStatusCode: 200,
		},
		{
			name: "test APISurface Update(...)",
			args: args{
				broker: &fakeBroker{
					validateBrokerAPIVersion: func(version string) error { return nil },
					update: func(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
						return &broker.UpdateInstanceResponse{}, nil
					},
				},
				servicePath:   "/v2/service_instances/foo",
				serviceMethod: http.MethodPatch,
				request:       []byte("{}"),
			},
			wantStatusCode: 200,
		},
		{
			name: "test APISurface Deprovision(...)",
			args: args{
				broker: &fakeBroker{
					validateBrokerAPIVersion: func(version string) error { return nil },
					deprovision: func(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
						return &broker.DeprovisionResponse{}, nil
					},
				},
				servicePath:   "/v2/service_instances/foo",
				serviceMethod: http.MethodDelete,
				request:       []byte("{}"),
			},
			wantStatusCode: 200,
		},
		{
			name: "test APISurface Bind(...)",
			args: args{
				broker: &fakeBroker{
					validateBrokerAPIVersion: func(version string) error { return nil },
					bind: func(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
						return &broker.BindResponse{}, nil
					},
				},
				servicePath:   "/v2/service_instances/foo/service_bindings/bar",
				serviceMethod: http.MethodPut,
				request:       []byte("{}"),
			},
			wantStatusCode: 201,
		},
		{
			name: "test APISurface Unbind(...)",
			args: args{
				broker: &fakeBroker{
					validateBrokerAPIVersion: func(version string) error { return nil },
					unbind: func(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
						return &broker.UnbindResponse{}, nil
					},
				},
				servicePath:   "/v2/service_instances/foo/service_bindings/bar",
				serviceMethod: http.MethodDelete,
				request:       []byte("{}"),
			},
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, err := rest.NewAPISurface(tt.args.broker, metrics.New())
			handler := NewHTTPHandler(api)
			server := httptest.NewServer(handler)
			defer server.Close()
			u, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			u.Path = path.Join(u.Path, tt.args.servicePath)
			client := http.DefaultClient
			req, err := http.NewRequest(tt.args.serviceMethod, u.String(), bytes.NewReader(tt.args.request))
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != tt.wantStatusCode {
				t.Fatalf("Received status code: %d, want: %d\n", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

// BusinessLogic provides an implementation of the broker.Interface interface.
type fakeBroker struct {
	validateBrokerAPIVersion func(version string) error
	getCatalog               func(c *broker.RequestContext) (*broker.CatalogResponse, error)
	provision                func(request *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error)
	deprovision              func(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error)
	lastOperation            func(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error)
	bind                     func(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error)
	unbind                   func(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error)
	update                   func(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error)
}

var _ broker.Interface = &fakeBroker{}

func (b *fakeBroker) ValidateBrokerAPIVersion(version string) error {
	return b.validateBrokerAPIVersion(version)
}
func (b *fakeBroker) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	return b.getCatalog(c)
}
func (b *fakeBroker) Provision(request *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
	return b.provision(request, c)
}
func (b *fakeBroker) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	return b.deprovision(request, c)
}
func (b *fakeBroker) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
	return b.lastOperation(request, c)
}
func (b *fakeBroker) Bind(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
	return b.bind(request, c)
}
func (b *fakeBroker) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
	return b.unbind(request, c)
}
func (b *fakeBroker) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
	return b.update(request, c)
}
