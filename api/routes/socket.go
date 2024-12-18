package routes

import (
	"github.com/Kleinish/dascr-board/ws"
	"github.com/go-chi/chi"
)

// SocketRoutes provide websocket routes
func SocketRoutes(h *ws.Hub) chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", ws.ServeWS(h))
	return r
}
