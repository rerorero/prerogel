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
	// APIPathLoadVertex is path for loading vertex
	APIPathLoadVertex = "/ctl/load-vertex"
	// APIPathLoadPartitionVertices is path for loading vertices of each partition
	APIPathLoadPartitionVertices = "/ctl/load-partition"
	// APIPathStartSuperStep is path for starting super step
	APIPathStartSuperStep = "/ctl/start-superstep"
	// APIPathShowAggregatedValue is path for showing aggregator values
	APIPathShowAggregatedValue = "/ctl/agg"
	// APIPathShutdown is path for shutdown
	APIPathShutdown = "/ctl/shutdown"
	// APIPathGetVertexValue is path for getting vertex valu
	APIPathGetVertexValue = "/ctl/vertex/value"
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
	s.mux.Handle(APIPathLoadVertex, http.HandlerFunc(s.loadVertexHandler))
	s.mux.Handle(APIPathLoadPartitionVertices, http.HandlerFunc(s.loadParionVerticesHandler))
	s.mux.Handle(APIPathStartSuperStep, http.HandlerFunc(s.startSuperstepHandler))
	s.mux.Handle(APIPathShowAggregatedValue, http.HandlerFunc(s.showAggValueHandler))
	s.mux.Handle(APIPathShutdown, http.HandlerFunc(s.shutdownHandler))
	s.mux.Handle(APIPathGetVertexValue, http.HandlerFunc(s.getVertexValueHandler))

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

func (s *CtrlServer) redirectAndWait(w http.ResponseWriter, r *http.Request, message interface{}) (interface{}, error) {
	if err := json.NewDecoder(r.Body).Decode(message); err != nil {
		return nil, errors.Wrap(err, "failed to parse request body")
	}
	return s.requestAndWait(w, message)
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

func (s *CtrlServer) loadVertexHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.redirectAndWait(w, r, &command.LoadVertex{})
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err)
		return
	}

	ack, ok := res.(*command.LoadVertexAck)
	if !ok {
		s.respondError(w, http.StatusInternalServerError, errors.New(fmt.Sprintf("not load vertex ack: %#v", res)))
		return
	}

	if ack.Error != "" {
		s.respondError(w, http.StatusInternalServerError, errors.New(ack.Error))
		return
	}

	s.respond(w, http.StatusOK, ack)
}

func (s *CtrlServer) loadParionVerticesHandler(w http.ResponseWriter, r *http.Request) {
	s.actorCtx.Send(s.coordinator, &command.LoadPartitionVertices{})
	s.respond(w, http.StatusOK, nil)
}

func (s *CtrlServer) startSuperstepHandler(w http.ResponseWriter, r *http.Request) {
	s.actorCtx.Send(s.coordinator, &command.StartSuperStep{})
	s.respond(w, http.StatusOK, nil)
}

func (s *CtrlServer) showAggValueHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.requestAndWait(w, &command.ShowAggregatedValue{})
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err)
		return
	}

	ack, ok := res.(*command.ShowAggregatedValueAck)
	if !ok {
		s.respondError(w, http.StatusInternalServerError, errors.New(fmt.Sprintf("not aggregated value ack: %#v", res)))
		return
	}

	s.respond(w, http.StatusOK, ack)
}

func (s *CtrlServer) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.requestAndWait(w, &command.Shutdown{})
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err)
		return
	}

	ack, ok := res.(*command.ShutdownAck)
	if !ok {
		s.respondError(w, http.StatusInternalServerError, errors.New(fmt.Sprintf("not shutdown ack: %#v", res)))
		return
	}

	s.respond(w, http.StatusOK, ack)
}

func (s *CtrlServer) getVertexValueHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.redirectAndWait(w, r, &command.GetVertexValue{})
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err)
		return
	}

	ack, ok := res.(*command.GetVertexValueAck)
	if !ok {
		s.respondError(w, http.StatusInternalServerError, errors.New(fmt.Sprintf("not get vertex value ack: %#v", res)))
		return
	}

	s.respond(w, http.StatusOK, ack)
}
