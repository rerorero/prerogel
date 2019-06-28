package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/command"
	"github.com/sirupsen/logrus"
)

// CtrlServer is server that provides management API
type CtrlServer struct {
	mux         *http.ServeMux
	server      *http.Server
	coordinator *actor.PID
	actorCtx    actor.SenderContext
	timeout     time.Duration
	logger      *logrus.Logger
}

// ErrRes is error response
type ErrRes struct {
	Message string `json:"message"`
}

const (
	// APIPathStats is path for stats
	APIPathStats = "/ctl/stats"
)

func newCtrlServer(coordinator *actor.PID, logger *logrus.Logger) *CtrlServer {
	s := &CtrlServer{
		mux:         http.NewServeMux(),
		coordinator: coordinator,
		actorCtx:    actor.EmptyRootContext,
		timeout:     30 * time.Second, // TODO: be configurable
		logger:      logger,
	}

	s.mux.Handle(APIPathStats, http.HandlerFunc(s.statsHandler))

	return s
}

func (s *CtrlServer) serve(ln net.Listener) error {
	server := &http.Server{
		Handler: s.mux,
	}
	s.server = server

	if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *CtrlServer) shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *CtrlServer) requestAndWait(w http.ResponseWriter, message interface{}) (interface{}, error) {
	fut := s.actorCtx.RequestFuture(s.coordinator, message, s.timeout)
	if err := fut.Wait(); err != nil {
		return nil, err
	}
	return fut.Result()
}

func (s *CtrlServer) respondError(w http.ResponseWriter, status int, err error) {
	s.logger.WithError(err).Error("internal server error")
	s.respond(w, status, &ErrRes{Message: err.Error()})
}

func (s *CtrlServer) respond(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		s.logger.WithError(err).Error("failed to encode response")
	}
}

func (s *CtrlServer) statsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.requestAndWait(w, &command.CoordinatorStats{})
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err)
		return
	}

	stat, ok := res.(*command.CoordinatorStatsAck)
	if !ok {
		s.respondError(w, http.StatusInternalServerError, errors.New(fmt.Sprintf("not stats ack: %#v", res)))
		return
	}

	s.respond(w, http.StatusOK, stat)
}
