VERSION := $(shell git describe --tags | tr -d v)

plugins/%.so: plugins/%.go
	@echo Building plugin: $*
	@GOARCH=amd64 go build -v -buildmode plugin -o plugins/$*.so plugins/$*.go

plugins: \
	plugins/encoding-xml.so

bin/nv: main.go plugins.go cmd/*.go config/*.go neighbor/*.go search/*.go
	@mkdir -p bin
	@go build -v -ldflags="-X 'main.version=$(VERSION)'" -o bin/nv main.go plugins.go

docs/%: bin/nv
	@bin/nv doc $* -o docs/$*

dist/nv_ext_%.deb: plugins/%.so
	@echo "Plugin: $*"
	@echo "Version: $(VERSION)"
	@echo "Architecture: amd64"
	@echo "Archive: $@"
	@echo "Prerequisites: $?"
	@$(eval BUILD_DIR="/tmp/build/nv_ext_$*_$(VERSION)_amd64/DEBIAN")
	@echo "Build dir: $(BUILD_DIR)"
	@rm -rf "$(BUILD_DIR)"
	@mkdir -p "$(BUILD_DIR)"
	@$(eval BUILD_LIB="/tmp/build/nv_ext_$*_$(VERSION)_amd64/usr/lib/nv")
	@echo "Build lib dir: $(BUILD_LIB)"
	@rm -rf "$(BUILD_LIB)"
	@mkdir -p "$(BUILD_LIB)"
	@cp "plugins/$*.so" "$(BUILD_LIB)/$*.so"
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
	@mv "/tmp/build/nv_ext_$*_$(VERSION)_amd64.deb" "dist/nv_ext_$*.deb"

dist/nv_%.deb: bin/nv #docs/man
	@echo "Version: $(VERSION)"
	@echo "Architecture: $*"
	@echo "Archive: $@"
	@echo "Prerequisites: $?"
	@$(eval BUILD_DIR="/tmp/build/nv_$(VERSION)_$*/DEBIAN")
	@echo "Build dir: $(BUILD_DIR)"
	@rm -rf "$(BUILD_DIR)"
	@mkdir -p "$(BUILD_DIR)"
	@$(eval BUILD_BIN="/tmp/build/nv_$(VERSION)_$*/usr/bin")
	@echo "Build bin dir: $(BUILD_BIN)"
	@rm -rf "$(BUILD_BIN)"
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

dist: \
	dist/nv_amd64.deb \
	dist/nv_ext_encoding-xml.deb

install: dist
	@for pkg in dist/*.deb; do sudo dpkg -i "$$pkg"; done

clean:
	@rm -f nv
	@rm -rf bin
	@rm -rf dist
	@rm -f plugins/*.so
	@rm -rf docs/cmd
	@rm -f coverage.txt
	@rm -rf /tmp/build

test:
	@go vet ./...
	@go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go mod verify

format:
	@gofmt -s -w .
