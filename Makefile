GOOS = linux
GOARCH = amd64
PROJECT_NAME = ext-auth-svc
MAIN_FILE = cmd/server/main.go
BUILD_PATH = build
SHELL := /bin/bash


init-build-dirs:
	@rm -rf $(BUILD_PATH)
	@mkdir $(BUILD_PATH)
	$(info $(BUILD_PATH) was created)

build-api: init-build-dirs
	$(info Will build API for $(GOOS) and $(GOARCH))
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(ARM_VERSION) go build -o $(BUILD_PATH)/$(PROJECT_NAME) $(MAIN_FILE)
	@echo "API build completed"

clean:
	@rm -rf $(BUILD_PATH)

run:
	go run $(MAIN_FILE)

deploy-ext-authz: build-api
	scp $(BUILD_PATH)/$(PROJECT_NAME) skrao@api.skrao.net:~/Dev/apis/extauth/
	$(call clean)

deploy-envoy:
	scp ./envoy-config/envoy.yml skrao@api.skrao.net:~/Dev/apis/envoy-proxy/

lint:
	$(HOME)/go/bin/golangci-lint run
