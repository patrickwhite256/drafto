package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/internal/packgen"
	"github.com/patrickwhite256/drafto/internal/twirpapi"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

var cardSets = []string{"ZNR"}

func main() {
	ds, err := datastore.New()
	if err != nil {
		log.Println(err)
		return
	}

	loader := &packgen.CardLoader{}

	for _, set := range cardSets {
		if err := loader.PreloadSet(set); err != nil {
			log.Println(err)
			return
		}
		log.Printf("preloaded %s\n", set)
	}

	handler := drafto.NewDraftoServer(&twirpapi.Server{
		Datastore:  ds,
		CardLoader: loader,
	})

	mux := http.NewServeMux()

	mux.Handle(handler.PathPrefix(), handler)

	log.Println("starting server...")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")), mux)
}
