package httpapi

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"time"
)

func randomCode(length int) (string, error) {
	buf := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	return enc.EncodeToString(buf)[:length], nil
}

func hashPairingCode(code string) []byte {
	sum := sha256.Sum256([]byte(code))
	return sum[:]
}

func randomNonce() (string, error) {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func normalizeContext(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return json.RawMessage("{}"), nil
	}
	var tmp any
	if err := json.Unmarshal(raw, &tmp); err != nil {
		return nil, err
	}
	b, err := json.Marshal(tmp)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

func parseRFC3339(v string) (time.Time, error) {
	if v == "" {
		return time.Time{}, errors.New("time required")
	}
	t, err := time.Parse(time.RFC3339Nano, v)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
