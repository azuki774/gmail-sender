package server

import (
	"context"
	"fmt"
	"gmail-sender/internal/model"
	"io"
	"net/http"

	"go.uber.org/zap"
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

func (s *Server) sendHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	title := q.Get("title")
	reqparam := model.MailContent{
		From:  from,
		To:    to,
		Title: title,
	}
	s.Logger.Info("parse parameter", zap.String("from", reqparam.From), zap.String("to", reqparam.To), zap.String("title", reqparam.Title))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v\n", err.Error())
		return
	}
	defer r.Body.Close()
	reqparam.Body = string(body)

	ctx := context.Background()

	err = s.Usecase.Send(ctx, reqparam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err.Error())
		return
	}
	fmt.Fprintf(w, "send\n")
}
