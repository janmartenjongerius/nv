VERSION := $(shell git describe --tags)

plugins/%.so: plugins/%.go
	@echo Building plugin: $*
	@GOARCH=amd64 go build -v -buildmode plugin -o plugins/$*.so plugins/$*.go

plugins: \
	plugins/encoding_xml.so

bin/nv: main.go plugins.go cmd/*.go config/*.go neighbor/*.go search/*.go
	@mkdir -p bin
	@go build -v -ldflags="-X 'main.version=$(VERSION)'" -o bin/nv main.go plugins.go

docs/%: bin/nv
	@bin/nv doc $* -o docs/$*

dist/ext_%.deb: plugins/%.so
	@echo "Plugin: $*"
	@echo "Version: $(VERSION)"
	@echo "Architecture: amd64"
	@echo "Archive: $@"
	@echo "Prerequisites: $?"
	@$(eval BUILD_DIR="/tmp/build/nv_ext-$*_$(VERSION)_amd64/DEBIAN")
	@echo "Build dir: $(BUILD_DIR)"
	@rm -rf "$(BUILD_DIR)"
	@mkdir -p "$(BUILD_DIR)"
	@$(eval BUILD_LIB="$(BUILD_DIR)/../usr/lib/nv")
	@echo "Build lib dir: $(BUILD_LIB)"
	@mkdir -p "$(BUILD_LIB)"
	@cp "plugins/$*.so" "$(BUILD_LIB)/ext-$*.so"
	@touch "$(BUILD_DIR)/control"
	@>>"$(BUILD_DIR)/control" echo "Package: nv-ext-$*"
	@>>"$(BUILD_DIR)/control" echo "Version: $(VERSION)"
	@>>"$(BUILD_DIR)/control" echo "Architecture: amd64"
	@>>"$(BUILD_DIR)/control" echo "Essential: no"
	@>>"$(BUILD_DIR)/control" echo "Priority: optional"
	@>>"$(BUILD_DIR)/control" echo "Maintainer: Jan-Marten de Boer"
	@>>"$(BUILD_DIR)/control" echo "Description: Nv extension $*"
	@dpkg-deb --build "$$(dirname "$(BUILD_DIR)")/"
	@mkdir -p dist
	@mv "/tmp/build/nv_ext-$*_$(VERSION)_amd64.deb" "dist/nv_ext-$*.deb"

dist/nv_%.deb: bin/nv #docs/man
	@echo "Version: $(VERSION)"
	@echo "Architecture: $*"
	@echo "Archive: $@"
	@echo "Prerequisites: $?"
	@$(eval BUILD_DIR="/tmp/build/nv_$(VERSION)_$*/DEBIAN")
	@echo "Build dir: $(BUILD_DIR)"
	@rm -rf "$(BUILD_DIR)"
	@mkdir -p "$(BUILD_DIR)"
	@$(eval BUILD_BIN="$(BUILD_DIR)/../usr/bin")
	@echo "Build bin dir: $(BUILD_BIN)"
	@mkdir -p "$(BUILD_BIN)"
	@cp bin/nv "$(BUILD_BIN)/nv"
	@touch "$(BUILD_DIR)/control"
	@>>"$(BUILD_DIR)/control" echo "Package: nv"
	@>>"$(BUILD_DIR)/control" echo "Version: $(VERSION)"
	@>>"$(BUILD_DIR)/control" echo "Architecture: $*"
	@>>"$(BUILD_DIR)/control" echo "Essential: no"
	@>>"$(BUILD_DIR)/control" echo "Priority: optional"
	@>>"$(BUILD_DIR)/control" echo "Maintainer: Jan-Marten de Boer"
	@>>"$(BUILD_DIR)/control" echo "Description: Environment lookup"
	@dpkg-deb --build "$$(dirname "$(BUILD_DIR)")/"
	@mkdir -p dist
	@mv "/tmp/build/nv_$(VERSION)_$*.deb" "dist/nv_$*.deb"

clean:
	@rm -f nv
	@rm -rf bin
	@rm -rf dist
	@rm -f plugins/*.so
	@rm -rf docs/cmd
	@rm -f coverage.txt

test:
	@go vet ./...
	@go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go mod verify

format:
	@gofmt -s -w .
