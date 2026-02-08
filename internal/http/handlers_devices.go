package httpapi

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"trustpin/internal/store"
)

type deviceActivateRequest struct {
	PairingCode string `json:"pairing_code"`
	PublicKey   string `json:"public_key"`
	Label       string `json:"label"`
}

type deviceActivateResponse struct {
	DeviceID string `json:"device_id"`
	State    string `json:"state"`
}

func (s *Server) handleDeviceActivate(w http.ResponseWriter, r *http.Request) {
	var req deviceActivateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.PairingCode == "" || req.PublicKey == "" {
		writeError(w, http.StatusBadRequest, "pairing_code_and_public_key_required")
		return
	}
	pubKey, err := base64.StdEncoding.DecodeString(req.PublicKey)
	if err != nil || len(pubKey) == 0 {
		writeError(w, http.StatusBadRequest, "invalid_public_key")
		return
	}

	var enrollment store.Enrollment
	id, err := s.Cache.Get(r.Context(), "enroll:code:"+req.PairingCode)
	if err == nil && id != "" {
		enrollment, err = s.Store.GetEnrollmentByID(r.Context(), id)
	} else {
		enrollment, err = s.Store.GetEnrollmentByPairingHash(r.Context(), hashPairingCode(req.PairingCode))
	}
	if err != nil {
		writeError(w, http.StatusNotFound, "enrollment_not_found")
		return
	}
	if enrollment.ExpiresAt.Before(nowUTC()) {
		writeError(w, http.StatusGone, "pairing_expired")
		return
	}
	if enrollment.State != store.EnrollmentPairing {
		writeError(w, http.StatusConflict, "invalid_enrollment_state")
		return
	}
	device, err := s.Store.CreateDevice(r.Context(), enrollment.TenantID, enrollment.UserID, pubKey, req.Label, store.DevicePending)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "device_create_failed")
		return
	}
	if err := s.Store.UpdateDeviceState(r.Context(), device.ID, store.DeviceActive); err != nil {
		writeError(w, http.StatusInternalServerError, "device_activate_failed")
		return
	}
	device.State = store.DeviceActive
	if err := s.Store.UpdateEnrollmentState(r.Context(), enrollment.ID, store.EnrollmentDeviceActive, &device.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "enrollment_update_failed")
		return
	}

	_ = s.Audit.Event(r.Context(), store.AuditEvent{
		TenantID: enrollment.TenantID,
		UserID:   &enrollment.UserID,
		DeviceID: &device.ID,
		ActorType: "device",
		ActorID:  device.ID,
		EventType: "device_activated",
	}, map[string]any{"enrollment_id": enrollment.ID})

	writeJSON(w, http.StatusCreated, deviceActivateResponse{DeviceID: device.ID, State: string(device.State)})
}

func (s *Server) handleDeviceRevoke(w http.ResponseWriter, r *http.Request) {
	s.handleDeviceStateChange(w, r, store.DeviceRevoked)
}

func (s *Server) handleDeviceSuspend(w http.ResponseWriter, r *http.Request) {
	s.handleDeviceStateChange(w, r, store.DeviceSuspended)
}

func (s *Server) handleDeviceResume(w http.ResponseWriter, r *http.Request) {
	s.handleDeviceStateChange(w, r, store.DeviceActive)
}

func (s *Server) handleDeviceStateChange(w http.ResponseWriter, r *http.Request, target store.DeviceState) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "device_id_required")
		return
	}
	device, err := s.Store.GetDevice(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "device_not_found")
		return
	}
	if !validDeviceTransition(device.State, target) {
		writeError(w, http.StatusConflict, "invalid_device_state")
		return
	}
	if err := s.Store.UpdateDeviceState(r.Context(), id, target); err != nil {
		writeError(w, http.StatusInternalServerError, "device_update_failed")
		return
	}
	if target == store.DeviceRevoked {
		_ = s.Store.UpdateEnrollmentStateByDeviceID(r.Context(), id, store.EnrollmentDeviceRev)
	}
	_ = s.Audit.Event(r.Context(), store.AuditEvent{
		TenantID: device.TenantID,
		UserID:   &device.UserID,
		DeviceID: &device.ID,
		ActorType: "system",
		ActorID:  "device_state_change",
		EventType: "device_" + string(target),
	}, map[string]any{"from": device.State, "to": target})

	writeJSON(w, http.StatusOK, map[string]string{"device_id": id, "state": string(target)})
}

func validDeviceTransition(current, target store.DeviceState) bool {
	switch current {
	case store.DevicePending:
		return target == store.DeviceActive || target == store.DeviceRevoked
	case store.DeviceActive:
		return target == store.DeviceSuspended || target == store.DeviceRevoked
	case store.DeviceSuspended:
		return target == store.DeviceActive || target == store.DeviceRevoked
	case store.DeviceRevoked:
		return false
	default:
		return false
	}
}

func nowUTC() time.Time {
	return time.Now().UTC()
}
