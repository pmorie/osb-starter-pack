package server

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/golang/glog"

	"github.com/pmorie/go-open-service-broker-skeleton/pkg/rest"
)

// Server might be redundant :)
type Server struct {
	api rest.APISurface
}

// New creates a new server.
func New(api rest.APISurface) *Server {
	return &Server{
		api: api,
	}
}

// Run creates the HTTP handler and begins to listen on the specified address.
func (s *Server) Run(ctx context.Context, addr string) error {
	listenAndServe := func(srv *http.Server) error {
		return srv.ListenAndServe()
	}
	return s.run(ctx, addr, listenAndServe)
}

// RunTLS creates the HTTPS handler based on the certifications that were passed
// and begins to listen on the specified address.
func (s *Server) RunTLS(ctx context.Context, addr string, cert string, key string) error {
	var decodedCert, decodedKey []byte
	var tlsCert tls.Certificate
	var err error
	decodedCert, err = base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return err
	}
	decodedKey, err = base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	tlsCert, err = tls.X509KeyPair(decodedCert, decodedKey)
	if err != nil {
		return err
	}
	listenAndServe := func(srv *http.Server) error {
		srv.TLSConfig = new(tls.Config)
		srv.TLSConfig.Certificates = []tls.Certificate{tlsCert}
		return srv.ListenAndServeTLS("", "")
	}
	return s.run(ctx, addr, listenAndServe)
}

func (s *Server) run(ctx context.Context, addr string, listenAndServe func(srv *http.Server) error) error {
	glog.Infof("Starting server on %s\n", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.api.Router,
	}
	go func() {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if srv.Shutdown(c) != nil {
			srv.Close()
		}
	}()
	return listenAndServe(srv)
}
