# Dataverse Broker

[![Build Status](https://travis-ci.org/dataverse-broker/dataverse-broker.svg?branch=master)](https://travis-ci.org/dataverse-broker/dataverse-broker "Travis")


A go service broker for [Dataverse](https://dataverse.org) that implements the
[Open Service Broker API](https://github.com/openservicebrokerapi/servicebroker).

This project is an implementation of [`osb-starter-pack`](https://github.com/pmorie/osb-starter-pack).

## Who should use this project?

You should use this project if you're interfacing a containerized application in Kubernetes that will utilize data stored on Dataverse.

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
$ go get github.com/dataverse-broker/dataverse-broker/cmd/dataverse-broker
```

Or clone the repo:

```console
$ cd $GOPATH/src && mkdir -p github.com/dataverse-broker && cd github.com/dataverse-broker && git clone git://github.com/dataverse-broker/dataverse-broker
```

Change into the project directory:

```console
$ cd $GOPATH/src/github.com/dataverse-broker/dataverse-broker
```

### Deploy broker using Helm

```console
$ make deploy-helm
```

### Deploy broker using Openshift

```console
$ make push deploy-openshift
```

Running either of these flavors of deploy targets will build the dataverse-broker binary,
build the image, deploy the broker into your Kubernetes, and add a
`ClusterServiceBroker` to the service-catalog.

## Using a Dataverse Service

### Provisioning a Service

To provision a service, click on a dataverse subtree icon on the service catalog. You will be prompted to enter your API-token for this subtree (optional). The broker will check that your token has the credentials necessary to access that dataverse. The provision step will fail if your token is invalid.

### Service Binding 

Upon successful provision, the bind step will create a secret with the dataverse coordiates and your credentials. Use this secret in your application.

## Goals of this project

- Make it easy for clients to interact with Dataverse
- Access datasets for use in containerized applications
