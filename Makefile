.DEFAULT_GOAL := help

BUNDLE ?= bundle
RUBY ?= ruby

.PHONY: help setup install test spec build console version clean ci

help:
	@printf "Available targets:\n"
	@printf "  make setup    Install gem dependencies\n"
	@printf "  make install  Alias for setup\n"
	@printf "  make test     Run the default test suite\n"
	@printf "  make spec     Run RSpec directly\n"
	@printf "  make build    Build the gem package\n"
	@printf "  make version  Smoke-check the CLI version command\n"
	@printf "  make console  Start an IRB console with the gem loaded\n"
	@printf "  make clean    Remove generated gem packages\n"
	@printf "  make ci       Run the local CI checks\n"

setup:
	$(BUNDLE) install

install: setup

test:
	$(BUNDLE) exec rake

spec:
	$(BUNDLE) exec rspec

build:
	$(BUNDLE) exec rake build

version:
	$(BUNDLE) exec $(RUBY) exe/opensource --version

console:
	$(BUNDLE) exec bin/console

clean:
	$(BUNDLE) exec rake clean

ci: test build version
