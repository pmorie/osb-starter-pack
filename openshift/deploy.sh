#!/bin/bash

IMAGE=$1
CA=`oc get secret -n kube-service-catalog -o go-template='{{ range .items }}{{ if eq .type "kubernetes.io/service-account-token" }}{{ index .data "service-ca.crt" }}{{end}}{{"\n"}}{{end}}' | tail -n 1`
oc process -f openshift/dataverse-broker.yaml -p IMAGE=$IMAGE -p BROKER_CA_CERT=$CA | oc apply -f -
