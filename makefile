.PHONY: help run test tools deps update_deps vendor

all:
	@echo "*******************************"
	@echo "************ Amore ************"
	@echo "*******************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  run         - run in dev mode"
	@echo "  test        - run go tests"
	@echo "  tools       - go get's a bunch of tools for dev"
	@echo "  deps        - get apps deps"
	@echo "  update-deps - resets apps vendored deps"
	@echo "  vendor 		 - vendors the dependancies"

##
## Tools
##
tools:
	go get github.com/pressly/glock
	go get github.com/pressly/gv

##
## Development
##
run:
	cd example && GO15VENDOREXPERIMENT=1 go run main.go

test:
	go test $$(GO15VENDOREXPERIMENT=1 go list ./... | grep -v '/vendor/')

dist-test:
	@GO15VENDOREXPERIMENT=1 $(MAKE) test

##
## Dependency mgmt
##
deps:
	@glock sync github.com/tanema/amore

update-deps:
	@echo "Updating Glockfile with package versions from GOPATH..."
	@rm -rf ./vendor
	@glock save github.com/tanema/amore
	@$(MAKE) vendor

vendor:
	@echo "Syncing dependencies into vendor directory..."
	@rm -rf ./vendor
	@gv < Glockfile
