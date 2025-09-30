package api_handlers

import (
	"log"
	"net/http"

	"github.com/jassi-singh/mini-forge/internal/services"
)

type ApiHandler struct {
	keyPool *services.KeyPool
}

func NewApiHandler(keyPool *services.KeyPool) *ApiHandler {
	return &ApiHandler{keyPool: keyPool}
}

func (h *ApiHandler) GetKey(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /get-key")
	key := h.keyPool.Get()
	log.Printf("Provided key: %s", key)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
}
