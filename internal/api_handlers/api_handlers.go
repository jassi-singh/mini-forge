package api_handlers

import (
	"log"
	"net/http"

	"github.com/jassi-singh/mini-forge/internal/keypool"
)

type ApiHandler struct {
	keyPool *utils.KeyPool
}

func NewApiHandler(keyPool *utils.KeyPool) *ApiHandler {
	return &ApiHandler{keyPool: keyPool}
}

func (h *ApiHandler) GetKey(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /get-key")
	key := h.keyPool.Get()
	log.Printf("Provided key: %s", key)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
}
