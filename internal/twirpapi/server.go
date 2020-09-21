package twirpapi

import (
	"github.com/patrickwhite256/drafto/internal/datastore"
	"github.com/patrickwhite256/drafto/internal/packgen"
)

type Server struct {
	Datastore *datastore.Datastore
	Packgen   *packgen.Generator
}
