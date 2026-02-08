package httpapi

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"trustpin/internal/crypto"
	"trustpin/internal/store"
)

type challengeInitRequest struct {
	TenantID string          `json:"tenant_id"`
	UserID   string          `json:"user_id"`
	DeviceID string          `json:"device_id"`
	Action   string          `json:"action"`
	Context  json.RawMessage `json:"context"`
}

type challengeInitResponse struct {
	ChallengeID string `json:"challenge_id"`
	State       string `json:"state"`
	IssuedAt    string `json:"issued_at"`
	ExpiresAt   string `json:"expires_at"`
}

func (s *Server) handleChallengeInit(w http.ResponseWriter, r *http.Request) {
	var req challengeInitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.TenantID == "" || req.UserID == "" || req.DeviceID == "" || req.Action == "" {
		writeError(w, http.StatusBadRequest, "missing_fields")
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
	device, err := s.Store.GetDevice(r.Context(), req.DeviceID)
	if err != nil {
		writeError(w, http.StatusNotFound, "device_not_found")
		return
	}
	if device.TenantID != req.TenantID || device.UserID != req.UserID {
		writeError(w, http.StatusConflict, "device_mismatch")
		return
	}
	if device.State != store.DeviceActive {
		writeError(w, http.StatusConflict, "device_not_active")
		return
	}
	ctxNormalized, err := normalizeContext(req.Context)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_context")
		return
	}
	contextHash := crypto.HashContext(ctxNormalized)
	nonce, err := randomNonce()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "nonce_failed")
		return
	}
	issued := time.Now().UTC()
	expires := issued.Add(s.Config.ChallengeTTL)
	challenge, err := s.Store.CreateChallenge(r.Context(), store.Challenge{
		TenantID:   req.TenantID,
		UserID:     req.UserID,
		DeviceID:   req.DeviceID,
		Nonce:      nonce,
		Action:     req.Action,
		Context:    ctxNormalized,
		ContextHash: contextHash,
		State:      store.ChallengeCreated,
		IssuedAt:   issued,
		ExpiresAt:  expires,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "challenge_create_failed")
		return
	}
	_ = s.Cache.Set(r.Context(), "challenge:"+challenge.ID, string(challenge.State), s.Config.ChallengeTTL)
	if err := s.Push.Send(r.Context(), challenge.DeviceID, map[string]any{"challenge_id": challenge.ID}); err != nil {
		_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengeFailed)
		writeError(w, http.StatusServiceUnavailable, "push_failed")
		return
	}
	_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengePushSent)
	_ = s.Audit.Event(r.Context(), store.AuditEvent{
		TenantID:   req.TenantID,
		UserID:     &req.UserID,
		DeviceID:   &req.DeviceID,
		ChallengeID: &challenge.ID,
		ActorType:  "system",
		ActorID:    "challenge_init",
		EventType:  "challenge_created",
	}, map[string]any{"action": req.Action})

	writeJSON(w, http.StatusCreated, challengeInitResponse{
		ChallengeID: challenge.ID,
		State:       string(store.ChallengePushSent),
		IssuedAt:    issued.Format(time.RFC3339Nano),
		ExpiresAt:   expires.Format(time.RFC3339Nano),
	})
}

type challengeGetResponse struct {
	ChallengeID string          `json:"challenge_id"`
	Action      string          `json:"action"`
	Context     json.RawMessage `json:"context"`
	Payload     canonicalPayload `json:"payload"`
	Canonical   string          `json:"canonical_payload"`
}

func (s *Server) handleChallengeGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "challenge_id_required")
		return
	}
	challenge, err := s.Store.GetChallenge(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "challenge_not_found")
		return
	}
	if time.Now().UTC().After(challenge.ExpiresAt) {
		_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengeExpired)
		writeError(w, http.StatusGone, "challenge_expired")
		return
	}
	if challenge.State != store.ChallengeCreated && challenge.State != store.ChallengePushSent {
		writeError(w, http.StatusConflict, "invalid_challenge_state")
		return
	}
	payload := canonicalPayloadFromChallenge(challenge)
	canonical := crypto.CanonicalJSON(payload.ToCanonical())
	writeJSON(w, http.StatusOK, challengeGetResponse{
		ChallengeID: challenge.ID,
		Action:      challenge.Action,
		Context:     challenge.Context,
		Payload:     payload,
		Canonical:   base64.StdEncoding.EncodeToString(canonical),
	})
}

type canonicalPayload struct {
	TenantID    string `json:"tenant_id"`
	UserID      string `json:"user_id"`
	DeviceID    string `json:"device_id"`
	ChallengeID string `json:"challenge_id"`
	Nonce       string `json:"nonce"`
	Action      string `json:"action"`
	IssuedAt    string `json:"issued_at"`
	ExpiresAt   string `json:"expires_at"`
	ContextHash string `json:"context_hash"`
}

func canonicalPayloadFromChallenge(c store.Challenge) canonicalPayload {
	return canonicalPayload{
		TenantID:    c.TenantID,
		UserID:      c.UserID,
		DeviceID:    c.DeviceID,
		ChallengeID: c.ID,
		Nonce:       c.Nonce,
		Action:      c.Action,
		IssuedAt:    c.IssuedAt.UTC().Format(time.RFC3339Nano),
		ExpiresAt:   c.ExpiresAt.UTC().Format(time.RFC3339Nano),
		ContextHash: c.ContextHash,
	}
}

func (p canonicalPayload) ToCanonical() crypto.CanonicalPayload {
	issued, _ := parseRFC3339(p.IssuedAt)
	expires, _ := parseRFC3339(p.ExpiresAt)
	return crypto.CanonicalPayload{
		TenantID:    p.TenantID,
		UserID:      p.UserID,
		DeviceID:    p.DeviceID,
		ChallengeID: p.ChallengeID,
		Nonce:       p.Nonce,
		Action:      p.Action,
		IssuedAt:    issued,
		ExpiresAt:   expires,
		ContextHash: p.ContextHash,
	}
}

type challengeDecisionRequest struct {
	DeviceID  string           `json:"device_id"`
	Signature string           `json:"signature"`
	Payload   canonicalPayload `json:"payload"`
	TOTPCode  string           `json:"totp_code,omitempty"`
}

func (s *Server) handleChallengeApprove(w http.ResponseWriter, r *http.Request) {
	s.handleChallengeDecision(w, r, store.ChallengeApproved)
}

func (s *Server) handleChallengeReject(w http.ResponseWriter, r *http.Request) {
	s.handleChallengeDecision(w, r, store.ChallengeRejected)
}

func (s *Server) handleChallengeDecision(w http.ResponseWriter, r *http.Request, decision store.ChallengeState) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "challenge_id_required")
		return
	}
	var req challengeDecisionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.DeviceID == "" || req.Signature == "" {
		writeError(w, http.StatusBadRequest, "device_id_and_signature_required")
		return
	}
	challenge, err := s.Store.GetChallenge(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "challenge_not_found")
		return
	}
	if time.Now().UTC().After(challenge.ExpiresAt) {
		_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengeExpired)
		writeError(w, http.StatusGone, "challenge_expired")
		return
	}
	if challenge.State == store.ChallengeApproved || challenge.State == store.ChallengeRejected {
		writeJSON(w, http.StatusOK, map[string]string{"challenge_id": challenge.ID, "state": string(challenge.State)})
		return
	}
	if challenge.State != store.ChallengeCreated && challenge.State != store.ChallengePushSent {
		writeError(w, http.StatusConflict, "invalid_challenge_state")
		return
	}
	device, err := s.Store.GetDevice(r.Context(), req.DeviceID)
	if err != nil {
		writeError(w, http.StatusNotFound, "device_not_found")
		return
	}
	if device.State != store.DeviceActive {
		writeError(w, http.StatusConflict, "device_not_active")
		return
	}
	if device.ID != challenge.DeviceID {
		writeError(w, http.StatusConflict, "device_mismatch")
		return
	}
	policy, err := s.Store.GetPolicy(r.Context(), challenge.TenantID, s.Config.DefaultTenantPolicy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "policy_fetch_failed")
		return
	}
	if policy.TOTPRequired {
		if !s.Config.EnableTOTP {
			writeError(w, http.StatusPreconditionFailed, "totp_disabled")
			return
		}
		ok, err := s.TOTP.Verify(r.Context(), challenge.TenantID, challenge.UserID, req.TOTPCode)
		if err != nil || !ok {
			writeError(w, http.StatusPreconditionFailed, "totp_invalid")
			return
		}
	}

	if req.Payload.ChallengeID != challenge.ID || req.Payload.TenantID != challenge.TenantID ||
		req.Payload.UserID != challenge.UserID || req.Payload.DeviceID != challenge.DeviceID ||
		req.Payload.Nonce != challenge.Nonce || req.Payload.Action != challenge.Action ||
		req.Payload.ContextHash != challenge.ContextHash {
		_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengeFailed)
		writeError(w, http.StatusPreconditionFailed, "payload_mismatch")
		return
	}
	issued, err := parseRFC3339(req.Payload.IssuedAt)
	if err != nil || !issued.Equal(challenge.IssuedAt.UTC()) {
		_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengeFailed)
		writeError(w, http.StatusPreconditionFailed, "issued_at_mismatch")
		return
	}
	expires, err := parseRFC3339(req.Payload.ExpiresAt)
	if err != nil || !expires.Equal(challenge.ExpiresAt.UTC()) {
		_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengeFailed)
		writeError(w, http.StatusPreconditionFailed, "expires_at_mismatch")
		return
	}
	if ok, err := s.Cache.SetNX(r.Context(), "nonce:"+req.Payload.Nonce, "1", time.Until(challenge.ExpiresAt)); err != nil {
		writeError(w, http.StatusServiceUnavailable, "nonce_store_unavailable")
		return
	} else if !ok {
		writeError(w, http.StatusConflict, "nonce_reuse")
		return
	}
	canonical := crypto.CanonicalJSON(req.Payload.ToCanonical())
	sig, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_signature")
		return
	}
	if !ed25519.Verify(device.PublicKey, canonical, sig) {
		_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengeFailed)
		writeError(w, http.StatusPreconditionFailed, "signature_invalid")
		return
	}
	_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, decision)
	_ = s.Audit.Event(r.Context(), store.AuditEvent{
		TenantID:   challenge.TenantID,
		UserID:     &challenge.UserID,
		DeviceID:   &challenge.DeviceID,
		ChallengeID: &challenge.ID,
		ActorType:  "device",
		ActorID:    challenge.DeviceID,
		EventType:  "challenge_" + string(decision),
	}, map[string]any{"action": challenge.Action})

	writeJSON(w, http.StatusOK, map[string]string{"challenge_id": challenge.ID, "state": string(decision)})
}

func (s *Server) handleChallengeStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "challenge_id_required")
		return
	}
	challenge, err := s.Store.GetChallenge(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "challenge_not_found")
		return
	}
	if time.Now().UTC().After(challenge.ExpiresAt) && challenge.State != store.ChallengeApproved && challenge.State != store.ChallengeRejected {
		_ = s.Store.UpdateChallengeState(r.Context(), challenge.ID, store.ChallengeExpired)
		challenge.State = store.ChallengeExpired
	}
	writeJSON(w, http.StatusOK, map[string]string{"challenge_id": challenge.ID, "state": string(challenge.State)})
}
