OUT := project-dinner
PKG := project-dinner/pkg
PKG_LIST := $(shell go list ${PKG}/...)

# This version-strategy uses a manual value to set the version string
VERSION := 0.0.12
BUILD_ENV := "production"
APP_NAME_STAGING := project-dinner-staging
APP_NAME_PRODUCTION := project-dinner-production

# Where to push the docker image.
REGISTRY ?= mbvofdocker

SRC_DIRS := cmd pkg # directories which hold app source (not vendored)

dev_mode := DEVELOPMENT_MODE=true
send_grid_usr := SEND_GRID_USER=c24d9e9cb588d3
send_grid_key := SEND_GRID_API_KEY=11b128ec80a258
host := HOST=smtp.mailtrap.io
mail_port := MAIL_PORT=25
is_staging := IS_STAGING=false
db_url := DATABASE_URL="postgres://postgres:postgres@localhost/bandlokaler_test?sslmode=disable"
#db_url := DATABASE_URL="postgres://xcizeumdmzsahd:bd31860486d2182540c497365e85644ed39844ca12c030a6a3d79120e668d083@ec2-52-31-233-101.eu-west-1.compute.amazonaws.com:5432/dd6gntnimjqmji"
# Used internally.  Users should pass GOOS and/or GOARCH.
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

BUILD_IMAGE ?= golang:1.14-alpine

BUILD_DIRS := bin/$(OS)_$(ARCH)     \
              .go/bin/$(OS)_$(ARCH) \
              .go/cache

vet:
	@go vet ${PKG_LIST}

format:
	@go fmt ./...

dev:
	air

run-dev:
	./start.sh ${dev_mode} ${send_grid_usr} ${send_grid_key} ${host} ${mail_port} ${is_staging} ${db_url}

# Work around for connecting to DB on host when running the app in docker
dev-build: vet format
	@docker build --rm --build-arg VERSION=${VERSION} -t ${REGISTRY}/${OUT}:${VERSION} .
	@docker run --rm -e DEVELOPMENT_MODE=true -e SEND_GRID_USER=00d529b7247ff7 -e SEND_GRID_API_KEY=84af903e641a55 -e HOST=smtp.mailtrap.io -e MAIL_PORT=25 -e IS_STAGING=true -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal/bandlokaler_test?sslmode=disable" -p 5000:5000 ${REGISTRY}/${OUT}:${VERSION}

test: vet $(BUILD_DIRS)
	@docker run                                                 	\
    	    -i                                                      \
    	    --rm                                                    \
    	    -u $$(id -u):$$(id -g)                                  \
    	    -v $$(pwd):/src                                         \
    	    -w /src                                                 \
    	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
    	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
    	    -v $$(pwd)/.go/cache:/.cache                            \
    	    $(BUILD_IMAGE)                                          \
            	    /bin/sh -c "                                    \
            	        ARCH=$(ARCH)                                \
            	        OS=$(OS)                                    \
            	        VERSION=$(VERSION)                          \
            	        ./build/test.sh $(SRC_DIRS)                 \
            	    "

push-docker-hub:
	@docker push ${REGISTRY}/${OUT}:${VERSION}

# Makes a build for production
container: vet format test
	@docker build --rm --build-arg VERSION=${VERSION} -t ${REGISTRY}/${OUT}:${VERSION} .

deploy-stag: container bin-clean
	@docker tag ${REGISTRY}/${OUT}:${VERSION} registry.heroku.com/${APP_NAME_STAGING}/web
	@docker push registry.heroku.com/${APP_NAME_STAGING}/web
	heroku container:release web -a ${APP_NAME_STAGING}
	docker rmi registry.heroku.com/${APP_NAME_STAGING}/web

deploy-prod: bin-clean
	@docker tag ${REGISTRY}/${OUT}:${VERSION} registry.heroku.com/${APP_NAME_PRODUCTION}/web
	@docker push registry.heroku.com/${APP_NAME_PRODUCTION}/web
	heroku container:release web -a ${APP_NAME_PRODUCTION}
	docker rmi registry.heroku.com/${APP_NAME_PRODUCTION}/web

$(BUILD_DIRS):
	@mkdir -p $@

bin-clean:
	rm -rf .go bin
