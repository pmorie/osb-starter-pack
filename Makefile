# If the USE_SUDO_FOR_DOCKER env var is set, prefix docker commands with 'sudo'
ifdef USE_SUDO_FOR_DOCKER
	SUDO_CMD = sudo
endif

IMAGE ?= quay.io/osb-starter-pack/servicebroker
TAG ?= $(shell git describe --tags --always)
PULL ?= IfNotPresent

build:
	go build -i github.com/pmorie/osb-starter-pack/cmd/servicebroker

test:
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build --ldflags="-s" github.com/pmorie/osb-starter-pack/cmd/servicebroker

image: linux
	cp servicebroker image/
	$(SUDO_CMD) docker build image/ -t "$(IMAGE):$(TAG)"

clean:
	rm -f servicebroker

push: image
	$(SUDO_CMD) docker push "$(IMAGE):$(TAG)"

deploy-helm: image
	helm install charts/servicebroker \
	--name broker-skeleton --namespace broker-skeleton \
	--set imagePullPolicy="$(PULL)"

deploy-openshift: image
	oc new-project osb-starter-pack
	oc process -f openshift/starter-pack.yaml | oc create -f -

create-ns:
	kubectl create ns test-ns

provision: create-ns
	kubectl apply -f manifests/service-instance.yaml 

bind:
	kubectl apply -f manifests/service-binding.yaml	

.PHONY: build test linux image clean push deploy-help deploy-openshift create-ns provision bind
