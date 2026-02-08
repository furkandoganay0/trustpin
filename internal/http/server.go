package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"trustpin/internal/audit"
	"trustpin/internal/cache"
	"trustpin/internal/config"
	"trustpin/internal/push"
	"trustpin/internal/store"
	"trustpin/internal/totp"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Config config.Config
	Store  *store.Postgres
	Cache  *cache.Redis
	Audit  audit.Logger
	Push   push.Provider
	TOTP   totp.Verifier
}

func New(cfg config.Config, store *store.Postgres, cache *cache.Redis, audit audit.Logger, pushProvider push.Provider, totpVerifier totp.Verifier) *Server {
	return &Server{Config: cfg, Store: store, Cache: cache, Audit: audit, Push: pushProvider, TOTP: totpVerifier}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(jsonOnly)
	r.Use(s.rateLimit)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/enrollments/init", s.handleEnrollmentInit)
		r.Post("/devices/activate", s.handleDeviceActivate)
		r.Post("/devices/{id}/revoke", s.handleDeviceRevoke)
		r.Post("/devices/{id}/suspend", s.handleDeviceSuspend)
		r.Post("/devices/{id}/resume", s.handleDeviceResume)

		r.Post("/auth/challenges/init", s.handleChallengeInit)
		r.Get("/auth/challenges/{id}", s.handleChallengeGet)
		r.Post("/auth/challenges/{id}/approve", s.handleChallengeApprove)
		r.Post("/auth/challenges/{id}/reject", s.handleChallengeReject)
		r.Get("/auth/challenges/{id}/status", s.handleChallengeStatus)

		r.Get("/audit/events", s.handleAuditEvents)
	})

	return r
}

func (s *Server) Serve(ctx context.Context) error {
	srv := &http.Server{
		Addr:              s.Config.HTTPAddr,
		Handler:           s.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()
	return srv.ListenAndServe()
}

func jsonOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.Config.RateLimitPerMinute <= 0 {
			next.ServeHTTP(w, r)
			return
		}
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			next.ServeHTTP(w, r)
			return
		}
		key := "rl:" + tenantID + ":" + time.Now().UTC().Format("200601021504")
		count, err := s.Cache.Incr(r.Context(), key, time.Minute)
		if err == nil && int(count) > s.Config.RateLimitPerMinute {
			writeError(w, http.StatusTooManyRequests, "rate_limit_exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}
