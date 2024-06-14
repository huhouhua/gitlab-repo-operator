# Copyright 2024 The huhouhua Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http:www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build all by default, even if it's not first

GO := go
# include the common make file
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR) && pwd -P))
endif

.PHONY: copyright.verify
copyright.verify: tools.verify.licctl
	@echo "===========> Verifying the boilerplate headers for all files"
	@licctl --check -f $(ROOT_DIR)/boilerplate.txt $(ROOT_DIR) --skip-dirs=.idea

.PHONY: copyright.add
copyright.add: tools.verify.licctl
	@echo $(ROOT_DIR)
	@licctl -v -f $(ROOT_DIR)/boilerplate.txt $(ROOT_DIR) --skip-dirs=.idea


.PHONY: install.licctl
install.addlicense:
	@$(GO) install github.com/seacraft/licctl@latest

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: tidy
tidy:
	@$(GO) mod tidy

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"