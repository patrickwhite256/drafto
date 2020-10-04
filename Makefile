.PHONY: generate dev build

# TODO: make these go to the proper places
generate:
	protoc --twirp_out=. --go_out=. ./rpc/drafto/service.proto
	protoc --js_out=import_style=commonjs,binary:. --proto_path=. --twirp_js_out=. ./rpc/drafto/service.proto

dev:
	AWS_REGION=us-west-2 HTTP_PORT=8000 go run ./cmd/server/...

build:
	cd drafto-web && npm run-script build
	go-bindata -o ./cmd/server/bindata.go -pkg main -prefix "drafto-web/build" ./drafto-web/build/...
	GOOS=linux go build ./cmd/server/...
