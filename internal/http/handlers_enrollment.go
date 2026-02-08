package httpapi

import (
	"net/http"
	"time"

	"trustpin/internal/store"
)

type enrollmentInitRequest struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

type enrollmentInitResponse struct {
	EnrollmentID string `json:"enrollment_id"`
	PairingCode  string `json:"pairing_code"`
	ExpiresAt    string `json:"expires_at"`
}

func (s *Server) handleEnrollmentInit(w http.ResponseWriter, r *http.Request) {
	var req enrollmentInitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.TenantID == "" || req.UserID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id_and_user_id_required")
		return
	}
	if err := s.Store.EnsureTenant(r.Context(), req.TenantID); err != nil {
		writeError(w, http.StatusInternalServerError, "tenant_upsert_failed")
		return
	}
	if err := s.Store.EnsureUser(r.Context(), req.TenantID, req.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "user_upsert_failed")
		return
	}
	if err := s.Store.EnsurePolicy(r.Context(), req.TenantID, s.Config.DefaultTenantPolicy); err != nil {
		writeError(w, http.StatusInternalServerError, "policy_upsert_failed")
		return
	}
	code, err := randomCode(8)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "pairing_code_failed")
		return
	}
	expires := time.Now().UTC().Add(s.Config.EnrollmentTTL)
	enrollment, err := s.Store.CreateEnrollment(r.Context(), req.TenantID, req.UserID, hashPairingCode(code), expires)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "enrollment_create_failed")
		return
	}
	_ = s.Cache.Set(r.Context(), "enroll:code:"+code, enrollment.ID, s.Config.EnrollmentTTL)
	_ = s.Cache.Set(r.Context(), "enroll:id:"+enrollment.ID, code, s.Config.EnrollmentTTL)

	_ = s.Audit.Event(r.Context(), store.AuditEvent{
		TenantID: req.TenantID,
		UserID:   &req.UserID,
		ActorType: "system",
		ActorID:  "enrollment_init",
		EventType: "enrollment_init",
	}, map[string]any{"enrollment_id": enrollment.ID})

	writeJSON(w, http.StatusCreated, enrollmentInitResponse{
		EnrollmentID: enrollment.ID,
		PairingCode:  code,
		ExpiresAt:    expires.Format(time.RFC3339Nano),
	})
}
