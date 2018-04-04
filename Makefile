# If the USE_SUDO_FOR_DOCKER env var is set, prefix docker commands with 'sudo'
ifdef USE_SUDO_FOR_DOCKER
	SUDO_CMD = sudo
endif

IMAGE ?= quay.io/brutto/dataverse-broker
TAG ?= $(shell git describe --tags --always)
PULL ?= IfNotPresent

build:
	go build -i github.com/SamiSousa/dataverse-broker/cmd/dataverse-broker

test: ## Runs the tests
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

linux: ## Builds a Linux executable
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -o dataverse-broker-linux --ldflags="-s" github.com/SamiSousa/dataverse-broker/cmd/dataverse-broker

image: linux ## Builds a Linux based image
	cp dataverse-broker-linux image/dataverse-broker
	$(SUDO_CMD) docker build image/ -t "$(IMAGE):$(TAG)"

clean: ## Cleans up build artifacts
	rm -f dataverse-broker
	rm -f dataverse-broker-linux
	rm -f image/dataverse-broker

push: image ## Pushes the image to dockerhub, REQUIRES SPECIAL PERMISSION
	$(SUDO_CMD) docker push "$(IMAGE):$(TAG)"

deploy-helm: image ## Deploys image with helm
	helm upgrade --install broker-skeleton --namespace broker-skeleton \
	charts/servicebroker \
	--set image="$(IMAGE):$(TAG)",imagePullPolicy="$(PULL)"

deploy-openshift: image ## Deploys image to openshift
	oc get project dataverse-broker || oc new-project dataverse-broker
	oc create configmap whitelist-config --from-file=manifests/subtree-whitelist/
	oc process -f openshift/dataverse-broker.yaml -p IMAGE=$(IMAGE):$(TAG) | oc apply -f -

create-ns: ## Cleans up the namespaces
	kubectl create ns test-ns

provision: create-ns ## Provisions a service instance
	kubectl apply -f manifests/service-instance.yaml

bind: ## Creates a binding
	kubectl apply -f manifests/service-binding.yaml

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: build test linux image clean push deploy-helm deploy-openshift create-ns provision bind help
