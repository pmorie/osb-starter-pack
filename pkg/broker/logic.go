package broker

import (
	"net/http"

	"github.com/pmorie/go-open-service-broker-client/v2"
)

// BusinessLogic contains the business logic for the broker's operations.
type BusinessLogic interface {
	GetCatalog(w http.ResponseWriter, r *http.Request) (*v2.CatalogResponse, error)
	Provision(w http.ResponseWriter, r *http.Request) (*v2.ProvisionResponse, error)
	Deprovision(w http.ResponseWriter, r *http.Request) (*v2.DeprovisionResponse, error)
	LastOperation(w http.ResponseWriter, r *http.Request) (*v2.LastOperationResponse, error)
	Bind(w http.ResponseWriter, r *http.Request) (*v2.BindResponse, error)
	Unbind(w http.ResponseWriter, r *http.Request) (*v2.UnbindResponse, error)
}

// Implementation provides an implementation of the BusinessLogic interface.
type Implementation struct {
}

var _ BusinessLogic = &Implementation{}

func (*Implementation) GetCatalog(w http.ResponseWriter, r *http.Request) (*v2.CatalogResponse, error) {
	return nil, nil
}

func (b *Implementation) Provision(w http.ResponseWriter, r *http.Request) (*v2.ProvisionResponse, error) {
	return nil, nil
}

func (b *Implementation) Deprovision(w http.ResponseWriter, r *http.Request) (*v2.DeprovisionResponse, error) {
	return nil, nil
}

func (b *Implementation) LastOperation(w http.ResponseWriter, r *http.Request) (*v2.LastOperationResponse, error) {
	return nil, nil
}

func (b *Implementation) Bind(w http.ResponseWriter, r *http.Request) (*v2.BindResponse, error) {
	return nil, nil
}

func (b *Implementation) Unbind(w http.ResponseWriter, r *http.Request) (*v2.UnbindResponse, error) {
	return nil, nil
}
