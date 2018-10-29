package main

import (
	"github.com/donfrigo/img-manipulation/server/handlers"
	"github.com/donfrigo/img-manipulation/server/middlewares"
	"github.com/donfrigo/img-manipulation/server/socket"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {

	go func() {
		// Socket.io router for websocket communication
		ws := http.NewServeMux()

		server, err := socketio.NewServer(nil)
		if err != nil {
			log.Fatal(err)
		}

		server.On("connection", func(so socketio.Socket) {

			socketId := so.Id()

			socket.Clients[socketId] = so

			so.On("disconnect", func() {
				delete(socket.Clients, socketId)
			})

		})
		server.On("error", func(so socketio.Socket, err error) {
			log.Println("error:", err)
		})

		// provide default cors to the mux
		handler := cors.AllowAll().Handler(ws)

		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowCredentials: true,
		})

		// decorate existing handler with cors functionality set in c
		handler = c.Handler(server)

		ws.Handle("/socket.io/", handler)
		ws.Handle("/assets", http.FileServer(http.Dir("./assets")))

		// start server
		log.Println("Socket.io serving at localhost:5000")
		log.Fatal(http.ListenAndServe(":5000", handler))
	}()

	// Gorilla router for file handling
	r := mux.NewRouter()

	// middleware
	r.Use(middlewares.LoggingMiddleware)

	// routes
	r.HandleFunc("/upload", handlers.UploadHandler).Methods("PUT")

	// cors
	handler := cors.AllowAll().Handler(r)

	// start server
	log.Println("Gorilla router serving at localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", handler))
}
