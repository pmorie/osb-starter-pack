# osb-starter-pack

[![Build Status](https://travis-ci.org/pmorie/osb-starter-pack.svg?branch=master)](https://travis-ci.org/pmorie/osb-starter-pack "Travis")

A go quickstart for creating service brokers that implement the
[Open Service Broker API](https://github.com/openservicebrokerapi/servicebroker).
Broker authors implement an interface that uses the same types as the
[go-open-service-broker-client](https://github.com/pmorie/go-open-service-broker-client)
project.

## Who should use this project?

You should use this project if you're looking for a quick way to implement an
Open Service Broker.

## Prerequisites

You'll need:

- A running [kubernetes](https://github.com/kubernetes/kubernetes) cluster
- The [service-catalog](https://github.com/kubernetes-incubator/service-catalog)
  installed in that cluster

If you're using helm to deploy this broker, you'll need to have helm installed
in the cluster.

## Getting started

You can just clone this project and deploy it to start poking around right away:

```console
# clone the repo
$ git clone git://github.com/pmorie/go-open-service-broker-client

# deploy with helm
$ make deploy-helm
```

## Adding your business logic

To implement your broker, you fill out just a few methods and types in
`pkg/user` package:

- The `Options` type, which holds options for the broker
- The `AddFlags` function, which adds CLI flags for an Options
- The methods of the `BusinessLogic` type, which implements the broker's
  business logic
- The `NewBusinessLogic` function, which creates a BusinessLogic from the
  Options the program is run with

