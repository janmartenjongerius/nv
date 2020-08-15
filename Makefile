VERSION := $(shell git describe --tags)

plugins/encoding-xml.so: plugins/encoding-xml.go
	@go build -v -buildmode plugin -o plugins/encoding-xml.so plugins/encoding-xml.go

plugins: plugins/encoding-xml.so

bin/nv: main.go plugins.go cmd/*.go config/*.go neighbor/*.go
	@mkdir -p bin
	@go build -v -ldflags="-X 'main.version=$(VERSION)'" -o bin/nv main.go plugins.go

docs/cmd/%.md: nv
	@rm -f docs/cmd/*.md
	@./nv docs markdown -o docs/cmd

docs: docs/cmd/%.md

dist/nv_%.deb: bin/nv
	@echo $(VERSION)

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
