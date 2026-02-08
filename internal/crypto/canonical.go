package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type CanonicalPayload struct {
	TenantID    string
	UserID      string
	DeviceID    string
	ChallengeID string
	Nonce       string
	Action      string
	IssuedAt    time.Time
	ExpiresAt   time.Time
	ContextHash string
}

func HashContext(raw json.RawMessage) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func CanonicalJSON(p CanonicalPayload) []byte {
	issued := p.IssuedAt.UTC().Format(time.RFC3339Nano)
	expires := p.ExpiresAt.UTC().Format(time.RFC3339Nano)
	var b strings.Builder
	b.Grow(256)
	b.WriteString("{")
	b.WriteString(fmt.Sprintf("\"tenant_id\":%s,", mustJSON(p.TenantID)))
	b.WriteString(fmt.Sprintf("\"user_id\":%s,", mustJSON(p.UserID)))
	b.WriteString(fmt.Sprintf("\"device_id\":%s,", mustJSON(p.DeviceID)))
	b.WriteString(fmt.Sprintf("\"challenge_id\":%s,", mustJSON(p.ChallengeID)))
	b.WriteString(fmt.Sprintf("\"nonce\":%s,", mustJSON(p.Nonce)))
	b.WriteString(fmt.Sprintf("\"action\":%s,", mustJSON(p.Action)))
	b.WriteString(fmt.Sprintf("\"issued_at\":%s,", mustJSON(issued)))
	b.WriteString(fmt.Sprintf("\"expires_at\":%s,", mustJSON(expires)))
	b.WriteString(fmt.Sprintf("\"context_hash\":%s", mustJSON(p.ContextHash)))
	b.WriteString("}")
	return []byte(b.String())
}

func mustJSON(v string) string {
	b, _ := json.Marshal(v)
	return string(b)
}
