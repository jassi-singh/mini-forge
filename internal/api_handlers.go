package internal

import (
	"log"
	"net/http"
)

type ApiHandler struct {
	keyPool *KeyPool
}

func NewApiHandler(keyPool *KeyPool) *ApiHandler {
	return &ApiHandler{keyPool: keyPool}
}

func (h *ApiHandler) GetKey(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /get-key")
	key := h.keyPool.Get()
	log.Printf("Provided key: %s", key)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
}
