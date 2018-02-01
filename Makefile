build:
	go build -i github.com/pmorie/osb-starter-pack/cmd/servicebroker

test:
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build --ldflags="-s" github.com/pmorie/osb-starter-pack/cmd/servicebroker

image: linux
	cp servicebroker image/
	sudo docker build image/ -t quay.io/osb-starter-pack/servicebroker

clean:
	rm -f servicebroker

push: image
	docker push quay.io/osb-starter-pack/servicebroker:latest

deploy-helm: image
	helm install charts/servicebroker \
	--name broker-skeleton --namespace broker-skeleton \
	--set imagePullPolicy=Never

deploy-openshift: image
	oc new-project osb-starter-pack
	oc process -f openshift/starter-pack.yaml | oc create -f -
