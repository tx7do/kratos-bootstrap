# Not Support Windows

.PHONY: help wire conf ent build api openapi init all

ifeq ($(OS),Windows_NT)
    IS_WINDOWS:=1
endif

CURRENT_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

SRCS_MK := $(foreach dir, app, $(wildcard $(dir)/*/*/Makefile))

# generate protobuf api go code
api:
	cd api
	buf generate

# show help
help:
	@echo ""
	@echo "Usage:"
	@echo " make [target]"
	@echo ""
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
