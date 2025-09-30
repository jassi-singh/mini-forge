package internal

import (
	"database/sql"
	"net/http"
)

type ApiHandler struct {
	db    *sql.DB
	cache CacheStore
}

func NewApiHandler(db *sql.DB, cache CacheStore) *ApiHandler {
	return &ApiHandler{db: db, cache: cache}
}

func (h *ApiHandler) GetKey(w http.ResponseWriter, r *http.Request) {
}
