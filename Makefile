.PHONY: generate dev build

generate:
	protoc --twirp_out=./rpc/drafto --go_out=./rpc/drafto --proto_path=./rpc/drafto --twirp_opt=paths=source_relative --go_opt=paths=source_relative ./rpc/drafto/service.proto
	protoc --js_out=import_style=commonjs,binary:./drafto-web/src  --proto_path=./rpc/drafto --twirp_js_out=./drafto-web/src ./rpc/drafto/service.proto

dev:
	HOST=http://localhost:8000 AWS_REGION=us-west-2 HTTP_PORT=8000 go run ./cmd/server/...

build:
	cd drafto-web && npm run-script build
	go-bindata -o ./cmd/server/bindata.go -pkg main -prefix "drafto-web/build" ./drafto-web/build/...
	GOOS=linux go build ./cmd/server/...
