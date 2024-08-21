package socket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Angelosewase/chatbuddiesgo/internal/auth"
	sockeio "github.com/googollee/go-socket.io"
)

// Server wraps the Socket.IO server
type Server struct {
	SocketServer *sockeio.Server
}

// NewServer initializes a new Socket.IO server instance
func (server *Server) NewServer() error {
	server.SocketServer = sockeio.NewServer(nil)



	var userId string

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Extract cookies here
		jwtCookie, err := r.Cookie("jwt")

		if err != nil {
			http.Error(w, "there is no authorisation info", http.StatusUnauthorized)
			return
		}

		id, err := auth.ValidateToken(jwtCookie.Value)

		if err != nil {
			return
		}
		userId = id

		// Forward the request to the Socket.IO server
		server.SocketServer.ServeHTTP(w, r)
	})

	server.SocketServer.OnConnect("/", func(c sockeio.Conn) error {

		if userId == "" {
			return fmt.Errorf("invalid user id: %v", userId)
		}

		c.Join(userId)

		fmt.Println("connected ", c.ID())
		return nil
	})

	server.SocketServer.OnDisconnect("/", func(c sockeio.Conn, id string) {
		c.Leave(id)
		fmt.Println("disconnected ", c.ID())
	})
	return nil
}

// Start begins serving the Socket.IO server
func (server *Server) Start() error {
	go server.SocketServer.Serve()
	defer server.SocketServer.Close()

	log.Println("Socket.IO server started")
	return nil
}

// Close stops the Socket.IO server
func (server *Server) Close() {
	server.SocketServer.Close()
}


//this function should have/take the socketio server instacnce and emit a socket event to given client using the user id 
func(s *Server) EmitSocket(receiverId string, evnt string, msg interface{})error  {
  s.SocketServer.BroadcastToRoom("",receiverId,evnt,msg)
	return nil 
}
