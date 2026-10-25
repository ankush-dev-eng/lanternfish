package server

import (
	"context"
	"io/fs"
	"net/http"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/detection"
	"github.com/ankush-dev-eng/lanternfish/pkg/eventbus"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/simulator"
)

// Server coordinates HTTP endpoints, event streaming, and asset serving.
type Server struct {
	graph     *model.NetworkGraph
	decoys    *decoy.Registry
	bus       *eventbus.Bus
	detector  *detection.Engine
	simEngine *simulator.Engine
	staticFS  fs.FS
	mux       *http.ServeMux
}

// Config holds initialization parameters for the Server.
type Config struct {
	Graph     *model.NetworkGraph
	Decoys    *decoy.Registry
	Bus       *eventbus.Bus
	Detector  *detection.Engine
	SimEngine *simulator.Engine
	StaticFS  fs.FS
}

// New creates and registers an HTTP API and UI server.
func New(cfg Config) *Server {
	s := &Server{
		graph:     cfg.Graph,
		decoys:    cfg.Decoys,
		bus:       cfg.Bus,
		detector:  cfg.Detector,
		simEngine: cfg.SimEngine,
		staticFS:  cfg.StaticFS,
		mux:       http.NewServeMux(),
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	// API Endpoints
	s.mux.HandleFunc("/healthz", s.handleHealthz)
	s.mux.HandleFunc("/api/network", s.handleGetNetwork)
	s.mux.HandleFunc("/api/decoys", s.handleGetDecoys)
	s.mux.HandleFunc("/api/alerts", s.handleGetAlerts)
	s.mux.HandleFunc("/api/timeline", s.handleGetTimeline)
	s.mux.HandleFunc("/api/report.json", s.handleGetReportJSON)
	s.mux.HandleFunc("/api/report.md", s.handleGetReportMD)
	s.mux.HandleFunc("/api/stream", s.handleSSE)
	s.mux.HandleFunc("/api/simulate", s.handleSimulate)
	s.mux.HandleFunc("/api/reset", s.handleReset)
	s.mux.HandleFunc("/report", s.handleGetReportHTML)

	// Embedded Static UI files
	if s.staticFS != nil {
		fileServer := http.FileServer(http.FS(s.staticFS))
		s.mux.Handle("/", fileServer)
	}
}

// Handler returns the HTTP handler with global headers applied.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		s.mux.ServeHTTP(w, r)
	})
}

// ListenAndServe starts the HTTP listener on the given port.
func (s *Server) ListenAndServe(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.Handler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return srv.ListenAndServe()
}

// Shutdown gracefully terminates HTTP services.
func (s *Server) Shutdown(ctx context.Context) error {
	// Clean shutdown stub
	return nil
}
