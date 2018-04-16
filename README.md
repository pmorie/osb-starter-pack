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

### Using the Catalog

When logging in, if you are not automatically directed to the service catalog, you can do so manually by using the dropdown menu labelled "Add to Project" and selecting "Browse Catalog." There, you will see dataverse subtree icons among the list of services supported by the catalog.

### Utilizing a Service

To begin the process of provisioning and binding a dataverse subtree service, click on a dataverse subtree icon on the service catalog to generate a dialog window. The dialog window contains the following information in the order presented:

#### Information:

![Information](/screenshots/Information.png?raw=true "Information tab of a Dataverse Service")

Provides a description of the corresponding dataverse subtree, including plans, if more than one. This tab is purely educational, and has no bearing on the actual provisioning/binding phase of the service.

#### Configuration:

![Configuration](/screenshots/Configuration.png?raw=true "Configuration tab of a Dataverse Service")

Configure service to be provisioned/binded. Along with prompts to create a new project, you will be prompted to enter your API-token for this subtree (optional). The broker will check that your token has the necessary credentials to access that dataverse. During the Results tab, the provision step will fail if a provided token is invalid.

#### Service Binding

![Binding](/screenshots/Binding.png?raw=true "Binding tab of a Dataverse Service")

Allows for the option to bind the service and store the necessary information in a secret, or to create the binding at a later time inside a project.

#### Results

![Results](/screenshots/Results.png?raw=true "Results tab of a Dataverse Service")

At this page, the broker will attempt to provision and bind the service. Upon successful provision, the bind step will create a secret with the dataverse coordiates and your credentials. Use this secret with your created project to connect to the Dataverse Service.

## Consuming a generated secret

## Goals of this project

- Make it easy for clients to interact with Dataverse
- Access datasets for use in containerized applications
