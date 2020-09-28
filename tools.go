// +build tools

package main

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/thechriswalker/protoc-gen-twirp_js"
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
)
