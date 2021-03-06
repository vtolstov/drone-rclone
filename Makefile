.PHONY: install build test html update docker_build docker_image docker_deploy clean version rclone

VERSION := $(shell git describe --tags --always || git rev-parse --short HEAD)
DEPLOY_ACCOUNT := "vtolstov"
DEPLOY_IMAGE := "drone-rclone"

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

rclone:
	@curl -s --output rclone.zip http://beta.rclone.org/rclone-beta-latest-linux-amd64.zip
#	@curl -s --output rclone.zip http://downloads.rclone.org/rclone-current-linux-amd64.zip
	@unzip -o -q -j rclone.zip '*/rclone'
	@rm rclone.zip

install:
	glide install

build:
	go build -ldflags="$(EXTLDFLAGS)-s -w -X main.Version=$(VERSION)"

test:
	go test -v -coverprofile=coverage.txt

html:
	go tool cover -html=coverage.txt

update:
	glide up

docker_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags="-X main.Version=$(VERSION)"

docker_image:
	docker build -t $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE) .

docker: docker_build docker_image

docker_deploy:
ifeq ($(tag),)
	@echo "Usage: make $@ tag=<tag>"
	@exit 1
endif
	docker tag $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):latest $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):$(tag)
	docker push $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):$(tag)

clean:
	rm -rf coverage.txt $(DEPLOY_IMAGE)

version:
	@echo $(VERSION)
