package server

import (
	"context"
	"fmt"
	"net/http"
)

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

func (s *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := s.Usecase.RefineNewToken(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err.Error())
		return
	}
	fmt.Fprintf(w, "access token refreshed\n")
}
