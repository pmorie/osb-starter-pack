build:
	go build -i github.com/pmorie/osb-starter-pack/cmd/servicebroker

test:
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build --ldflags="-s" github.com/pmorie/osb-starter-pack/cmd/servicebroker

image: linux
	cp servicebroker image/
	docker build image/ -t osb-starter-pack/broker

clean:
	rm -f servicebroker

deploy-helm: image
	helm install charts/servicebroker \
	--name osb-starter-pack --namespace osb-starter-pack \
	--set imagePullPolicy=Never
