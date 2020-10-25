package main

import (
	"fmt"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/patrickwhite256/drafto/internal/auth"
	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/internal/discord"
	"github.com/patrickwhite256/drafto/internal/notifications"
	"github.com/patrickwhite256/drafto/internal/packgen"
	"github.com/patrickwhite256/drafto/internal/socket"
	"github.com/patrickwhite256/drafto/internal/twirpapi"
	"github.com/patrickwhite256/drafto/rpc/drafto"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ds, err := datastore.New()
	if err != nil {
		log.Println(err)
		return
	}

	discord, err := discord.New(os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Println(err)
		return
	}

	socketServer := socket.NewServer()

	notifications := notifications.New(discord, ds, socketServer)

	loader := &packgen.CardLoader{}

	err = loader.Preload()
	if err != nil {
		log.Printf("error preloading card data: %v", err)
		return
	}

	auther := auth.New(auth.Config{
		Datastore:     ds,
		DiscordKey:    os.Getenv("DISCORD_KEY"),
		DiscordSecret: os.Getenv("DISCORD_SECRET"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		Host:          os.Getenv("HOST"),
	})

	handler := drafto.NewDraftoServer(&twirpapi.Server{
		Datastore:     ds,
		CardLoader:    loader,
		Notifications: notifications,
	})

	mux := http.NewServeMux()

	mux.Handle(handler.PathPrefix(), auther.Middleware(handler))
	mux.Handle("/auth", auther.AuthHandler())
	mux.Handle("/auth/discord/callback", auther.CallbackHandler())
	mux.Handle("/ws/", socketServer.Handler())
	mux.Handle("/", http.StripPrefix("/", http.HandlerFunc(staticHandler)))

	log.Println("starting server...")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")), mux)
}

func staticHandler(rw http.ResponseWriter, req *http.Request) {
	srvPath := req.URL.Path
	// TODO: make the default index?
	if srvPath == "" || strings.HasPrefix(srvPath, "seat") || strings.HasPrefix(srvPath, "table") {
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
