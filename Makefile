.PHONY: tools
tools:
	@echo "Building tools"
	@./.tools/tcc/download.sh

.PHONY: extension-build
extension-build:
	@echo "Building extension"
	@cd .extra/extension && ./build.sh

.PHONY: extension-install
extension-install:
	@echo "Installing extension"
	@cd .extra/extension && ./install.sh

.PHONY: extension
extension: extension-build extension-install