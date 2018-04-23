# `osb-broker-lib`

[![Build Status](https://travis-ci.org/pmorie/osb-broker-lib.svg?branch=master)](https://travis-ci.org/pmorie/osb-broker-lib "Travis")
[![Go Report Card](https://goreportcard.com/badge/github.com/pmorie/osb-broker-lib)](https://goreportcard.com/report/github.com/pmorie/osb-broker-lib)
[![Godoc documentation](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://godoc.org/github.com/pmorie/osb-broker-lib/pkg)

A go library for developing an [Open Service
Broker](https://github.com/openservicebrokerapi/servicebroker), using the
[`pmorie/go-open-service-broker-client`](https://github.com/pmorie/go-open-service-broker-client)
OSB client library types. This project was originally created as part of the
[OSB Starter Pack](https://github.com/pmorie/osb-starter-pack) project.

## Who should use this library?

This library is most useful if you want to build your own broker from scratch
and use it in a project just the way you want. If you're looking for an
opinionated quickstart to easily start iterating on a new broker you should
instead check out the [OSB Starter
Pack](https://github.com/pmorie/osb-starter-pack).

## Example: serving broker catalog

```go
import (
    osb "github.com/pmorie/go-open-service-broker-client/v2"
    broker "github.com/pmorie/osb-broker-lib/pkg/"

    "gopkg.in/yaml.v2"
)

type MyBroker struct {
    // internal state goes here
}

func (b *MyBroker) GetCatalog(ctx *broker.RequestContext) (*osb.CatalogResponse, error) {
    response := &osb.CatalogResponse{}

    data := `
---
services:
- name: example-service
  id: 4f6e6cf6-ffdd-425f-a2c7-3c9258ad246e
  description: The example service!
  bindable: true
  metadata:
    displayName: "Example service"
    imageUrl: https://avatars2.githubusercontent.com/u/19862012?s=200&v=4
  plans:
  - name: default
    id: 86064792-7ea2-467b-af93-ac9694d96d5c
    description: The default plan for the service
    free: true
`

    err := yaml.Unmarshal([]byte(data), &response)
    if err != nil {
        return nil, err
    }

    return response, nil
}
```

## Goals

- Provide a simple, composable way to implement the OSB API

## Current Status

Currently this library is used on the [OSB Starter
Pack](https://github.com/pmorie/osb-starter-pack) project.
