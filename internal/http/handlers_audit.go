package httpapi

import (
	"net/http"
	"strconv"
)

func (s *Server) handleAuditEvents(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_required")
		return
	}
	limit := 100
	if v := r.URL.Query().Get("limit"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 && i <= 1000 {
			limit = i
		}
	}
	events, err := s.Store.ListAuditEvents(r.Context(), tenantID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "audit_fetch_failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": events})
}
