// Package rest provides a REST server implementation.
package rest

import (
	"context"
	"fmt"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
)

// Server represents a REST server.
type Server struct {
	server http.Server
}

// New constructs a new REST server.
//
// host: The host to listen on.
// port: The port to listen on.
//
// Returns a new REST server and a Gin engine(http handler).
func New(host, port string) (*Server, *gin.Engine) {
	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))

	addr := fmt.Sprintf("%s:%s", host, port)

	return &Server{
		server: http.Server{
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Addr:         addr,
			Handler:      router,
		},
	}, router
}

// Start starts the REST server.
//
// Returns an error if the server fails to start.
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Stop stops the REST server.
//
// ctx: The context of the request.
//
// Returns an error if the server fails to stop.
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
