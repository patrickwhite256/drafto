package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"

	"github.com/patrickwhite256/drafto/internal/datastore"
)

type contextKey string

const (
	sessionName                 = "_drafto_session"
	sessionKeyUserID            = "skUSERID"
	contextKeyUserID contextKey = "userid"
)

type Auth struct {
	datastore  *datastore.Datastore
	sessSecret string
	store      sessions.Store
}

type Config struct {
	Datastore     *datastore.Datastore
	DiscordKey    string
	DiscordSecret string
	SessionSecret string
	Host          string
}

func New(conf Config) *Auth {
	callbackURL := fmt.Sprintf("%s/auth/discord/callback", conf.Host)
	goth.UseProviders(discord.New(conf.DiscordKey, conf.DiscordSecret, callbackURL))
	store := sessions.NewCookieStore([]byte(conf.SessionSecret))
	return &Auth{
		datastore:  conf.Datastore,
		sessSecret: conf.SessionSecret,
		store:      store,
	}
}

func (a *Auth) AuthHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req = req.WithContext(context.WithValue(req.Context(), gothic.ProviderParamKey, "discord"))
		if user, err := gothic.CompleteUserAuth(rw, req); err == nil {
			a.loginUser(user, rw, req)
			http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
			return
		}
		gothic.BeginAuthHandler(rw, req)
	})
}

func (a *Auth) CallbackHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(rw, req)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("auth failed"))
			return
		}

		a.loginUser(user, rw, req)
	})
}

func (a *Auth) loginUser(gothUser goth.User, rw http.ResponseWriter, req *http.Request) {
	user, err := a.datastore.GetUserByDiscordID(req.Context(), gothUser.UserID)
	if err != nil {
		if err != datastore.NotFound {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err = a.datastore.NewUser(req.Context(), gothUser.UserID, gothUser.Name, gothUser.AvatarURL)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// ignore error - session will always be valid
	session, _ := a.store.Get(req, sessionName)
	session.Values[sessionKeyUserID] = user.ID
	if err = session.Save(req, rw); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
}

// UserID returns "" if no logged-in user is found
func UserID(ctx context.Context) string {
	id, _ := ctx.Value(contextKeyUserID).(string)
	return id
}

func (a *Auth) withUserID(req *http.Request) *http.Request {
	// ignore error - session will be valid
	session, _ := a.store.Get(req, sessionName)
	id, _ := session.Values[sessionKeyUserID].(string)
	ctx := context.WithValue(req.Context(), contextKeyUserID, id)
	return req.WithContext(ctx)
}

func (a *Auth) Middleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req = a.withUserID(req)
		inner.ServeHTTP(rw, req)
	})
}
