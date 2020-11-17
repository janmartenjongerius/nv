VERSION := $(shell git describe --tags | tr -d v)
GOARCH=amd64
LD_FLAGS=-X 'main.version=$(VERSION)'

go.sum: go.mod
	@go mod tidy

bin/nv: main.go */*.go go.sum
	@mkdir -p bin
	@go build -v \
		-ldflags="$(LD_FLAGS)" \
		-o bin/nv \
		main.go

docs/%: bin/nv
	@bin/nv doc $* -o docs/$*

dist/nv_%.deb: bin/nv docs/man
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
	@$(eval MAN_DIR="/tmp/build/nv_$(VERSION)_$*/usr/share/man/man1")
	@echo "Man pages: $(MAN_DIR)"
	@mkdir -p "$(MAN_DIR)"
	@cp docs/man/*.1 "$(MAN_DIR)"
	@for man in "$(MAN_DIR)/"*.1; do gzip -f "$$man"; done
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

install: dist/nv_amd64.deb
	@for pkg in $?; do sudo dpkg -i "$$pkg"; done

clean:
	@rm -f nv
	@rm -rf bin
	@rm -rf dist
	@rm -rf docs/man
	@rm -rf docs/markdown
	@rm -rf docs/rst
	@rm -rf docs/yaml
	@rm -f coverage.txt
	@rm -rf /tmp/build
	@rm -f coverage.html

test: coverage.txt
	@go vet ./...
	@go mod verify

coverage.txt: *.go */*.go
	@go test -covermode=atomic -count=1 -coverprofile=coverage.txt ./...

coverage.html: coverage.txt
	@go tool cover -html=coverage.txt -o coverage.html

format:
	@gofmt -s -w .
