export

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build_emulator
build_emulator:
	go build -o ./bin/emulator ./cmd/emulator/main.go

.PHONY: run_emulator
run_emulator:
	./bin/emulator


.PHONY: build_server
build_server:
	go build -o ./bin/server ./cmd/server_tcp/main.go

.PHONY: run_server
run_server:
	./bin/server

.PHONY: build_rmq
build_rmq:
	go build -o ./bin/proc_rmq ./cmd/processor_rmq/main.go

.PHONY: run_rmq
run_rmq:
	./bin/proc_rmq

.PHONY: build_api
build_api:
	go build -o ./bin/api ./cmd/server_api/main.go

.PHONY: run_api
run_api:
	./bin/api
