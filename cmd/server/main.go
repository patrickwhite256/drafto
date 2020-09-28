package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/internal/packgen"
	"github.com/patrickwhite256/drafto/internal/twirpapi"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func main() {
	ds, err := datastore.New()
	if err != nil {
		// log
		return
	}

	handler := drafto.NewDraftoServer(&twirpapi.Server{
		Datastore:  ds,
		CardLoader: &packgen.CardLoader{},
	})

	mux := http.NewServeMux()

	mux.Handle("api/", handler)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")), mux)
}
