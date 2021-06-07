# OSB Starter Pack

[![Build Status](https://travis-ci.org/pmorie/osb-starter-pack.svg?branch=master)](https://travis-ci.org/pmorie/osb-starter-pack "Travis")

A go quickstart for creating service brokers that implement the [Open Service
Broker API](https://github.com/openservicebrokerapi/servicebroker) based on
[`osb-broker-lib`](https://github.com/pmorie/osb-broker-lib). Broker authors
implement an interface that uses the same types as the
[`go-open-service-broker-client`](https://github.com/pmorie/go-open-service-broker-client)
project.

## Who should use this project?

You should use this project if you're looking for a quick way to implement an
Open Service Broker and start iterating on it.

## Prerequisites

You'll need:

- [`go`](https://golang.org/dl/)
- A running [Kubernetes](https://github.com/kubernetes/kubernetes) (or [openshift](https://github.com/openshift/origin/)) cluster
- The [service-catalog](https://github.com/kubernetes-incubator/service-catalog)
  [installed](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install.md)
  in that cluster

If you're using [Helm](https://helm.sh) to deploy this project, you'll need to
have it [installed](https://docs.helm.sh/using_helm/#quickstart) in the cluster.
Make sure [RBAC is correctly configured](https://docs.helm.sh/using_helm/#rbac)
for helm.

## Getting started

You can `go get` this repo or `git clone` it to start poking around right away.

The project comes ready with a minimal example service that you can easily
deploy and begin iterating on.

### Get the project

```console
$ go get github.com/pmorie/osb-starter-pack/cmd/servicebroker
```

Or clone the repo:

```console
$ cd $GOPATH/src && mkdir -p github.com/pmorie && cd github.com/pmorie && git clone git://github.com/pmorie/osb-starter-pack
```

Change into the project directory:

```console
$ cd $GOPATH/src/github.com/pmorie/osb-starter-pack
```

### Deploy broker using Helm

Deploy with Helm and pass custom image and tag name.
Note: This also pushes the generated image with docker.

```console
$ GO111MODULE=off IMAGE=myimage TAG=latest make push deploy-helm
```

### Deploy broker using Openshift

Deploy to OpenShift cluster by passing a custom image and tag name.
Note: You must already be logged into an OpenShift cluster. 
This also pushes the generated image with docker.

```console
$ GO111MODULE=off IMAGE=myimage TAG=latest make push deploy-openshift
```

Running either of these flavors of deploy targets will build the broker binary,
build the image, deploy the broker into your Kubernetes, and add a
`ClusterServiceBroker` to the service-catalog.

## Adding your business logic

To implement your broker, you fill out just a few methods and types in
`pkg/broker` package:

- The `Options` type, which holds options for the broker
- The `AddFlags` function, which adds CLI flags for an Options
- The methods of the `BusinessLogic` type, which implements the broker's
  business logic
- The `NewBusinessLogic` function, which creates a BusinessLogic from the
  Options the program is run with

## Goals of this project

- Make it extremely easy to create a new broker
- Have a batteries-included experience that gives you the good stuff right out
  of the box, for example:
  - Checks on who can make calls to the broker using Kubernetes
    [subject-access-reviews](https://kubernetes.io/docs/admin/accessing-the-api/)
  - Easy on-ramp to instrumenting your broker with
    [Prometheus](https://prometheus.io/) metrics
