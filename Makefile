
build:
	@cd ./cmd/go-auth; \
	go build -o ../bin/go-auth

run: build
	@cd ./cmd/bin; \
	./go-auth