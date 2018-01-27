build:
	go build -i github.com/pmorie/go-open-service-broker-skeleton/cmd/servicebroker

test:
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build --ldflags="-s" github.com/pmorie/go-open-service-broker-skeleton/cmd/servicebroker

image: linux
	cp servicebroker image/
	docker build image/ -t osb-skeleton/servicebroker

clean:
	rm -f servicebroker

deploy-helm: image
	helm install charts/servicebroker \
	--name broker-skeleton --namespace broker-skeleton \
	--set imagePullPolicy=Never
