package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

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
	mux.Handle("/", http.StripPrefix("/", http.HandlerFunc(staticHandler)))

	log.Println("starting server...")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")), mux)
}

func staticHandler(rw http.ResponseWriter, req *http.Request) {
	srvPath := req.URL.Path
	if srvPath == "" || strings.HasPrefix(srvPath, "player") {
		srvPath = "index.html"
	}

	ext := path.Ext(srvPath)
	ct := mime.TypeByExtension(ext)
	if ct != "" {
		rw.Header().Set("Content-Type", ct)
	}

	b, err := Asset(srvPath)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err := rw.Write(b); err != nil {
		log.Println("error writing response", err)
	}
}
