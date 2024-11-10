.PHONY: tools
tools:
	@echo "Building tools"
	@./.tools/tcc/download.sh

.PHONY: extension
extension:
	@echo "Building extension"
	@cd .extra/extension && ./build.sh
	@echo "Installing extension"
	@cd .extra/extension && ./install.sh

	