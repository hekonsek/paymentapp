PACKAGES := github.com/hekonsek/paymentapp/api github.com/hekonsek/paymentapp/payments \
github.com/hekonsek/paymentapp/main github.com/hekonsek/paymentapp/cmd

all: format silent-test validate-contract build

format:
	GO111MODULE=on go fmt $(PACKAGES) github.com/hekonsek/paymentapp/contract

build:
	GO111MODULE=on go build -o paymentapp main/main.go

test: build
	GO111MODULE=on go test -v $(PACKAGES)

silent-test:
	GO111MODULE=on go test $(PACKAGES)

validate-contract:
	GO111MODULE=on go test github.com/hekonsek/paymentapp/contract

docker-build: build
	docker build -t hekonsek/paymentapp .

docker-push: docker-build
	docker push hekonsek/paymentapp

lint: format
	~/go/bin/golint $(PACKAGES) github.com/hekonsek/paymentapp/contract