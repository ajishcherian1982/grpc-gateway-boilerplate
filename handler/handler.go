package handler

import (
	"github.com/ajishcherian1982/grpc-gateway-boilerplate/store"
	"github.com/rs/zerolog"
)

// Handler definition
type Handler struct {
	logger *zerolog.Logger
	us     *store.UserStore
	as     *store.ArticleStore
}

// New returns a new handler with logger and database
func New(l *zerolog.Logger, us *store.UserStore, as *store.ArticleStore) *Handler {
	return &Handler{logger: l, us: us, as: as}
}
