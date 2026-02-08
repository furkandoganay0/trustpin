package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"trustpin/internal/crypto"
)

type challengeGetResponse struct {
	ChallengeID string          `json:"challenge_id"`
	Action      string          `json:"action"`
	Context     json.RawMessage `json:"context"`
	Payload     canonicalPayload `json:"payload"`
	Canonical   string          `json:"canonical_payload"`
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

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	switch os.Args[1] {
	case "keygen":
		cmdKeygen()
	case "enroll":
		cmdEnroll()
	case "activate":
		cmdActivate()
	case "fetch":
		cmdFetch()
	case "approve", "reject":
		cmdDecision(os.Args[1])
	default:
		usage()
	}
}

func usage() {
	fmt.Println("trustpin-cli keygen|enroll|activate|fetch|approve|reject")
}

func cmdKeygen() {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	fmt.Println("PUBLIC_KEY_BASE64=" + base64.StdEncoding.EncodeToString(pub))
	fmt.Println("PRIVATE_KEY_BASE64=" + base64.StdEncoding.EncodeToString(priv))
}

func cmdEnroll() {
	fs := flag.NewFlagSet("enroll", flag.ExitOnError)
	base := fs.String("base", "http://localhost:8080", "base url")
	tenant := fs.String("tenant", "t1", "tenant id")
	user := fs.String("user", "u1", "user id")
	fs.Parse(os.Args[2:])

	body := map[string]string{"tenant_id": *tenant, "user_id": *user}
	b, _ := json.Marshal(body)
	resp := doReq("POST", *base+"/v1/enrollments/init", b, nil)
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

func cmdActivate() {
	fs := flag.NewFlagSet("activate", flag.ExitOnError)
	base := fs.String("base", "http://localhost:8080", "base url")
	pairing := fs.String("code", "", "pairing code")
	pub := fs.String("pub", "", "public key base64")
	label := fs.String("label", "device", "label")
	fs.Parse(os.Args[2:])
	body := map[string]string{"pairing_code": *pairing, "public_key": *pub, "label": *label}
	b, _ := json.Marshal(body)
	resp := doReq("POST", *base+"/v1/devices/activate", b, nil)
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

func cmdFetch() {
	fs := flag.NewFlagSet("fetch", flag.ExitOnError)
	base := fs.String("base", "http://localhost:8080", "base url")
	id := fs.String("id", "", "challenge id")
	fs.Parse(os.Args[2:])
	resp := doReq("GET", *base+"/v1/auth/challenges/"+*id, nil, nil)
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

func cmdDecision(action string) {
	fs := flag.NewFlagSet("decision", flag.ExitOnError)
	base := fs.String("base", "http://localhost:8080", "base url")
	id := fs.String("id", "", "challenge id")
	device := fs.String("device", "", "device id")
	priv := fs.String("priv", "", "private key base64")
	fs.Parse(os.Args[2:])

	resp := doReq("GET", *base+"/v1/auth/challenges/"+*id, nil, nil)
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		io.Copy(os.Stdout, resp.Body)
		return
	}
	var payloadResp challengeGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&payloadResp); err != nil {
		fmt.Println("decode failed:", err)
		return
	}
	if *device == "" {
		*device = payloadResp.Payload.DeviceID
	}
	if *priv == "" {
		fmt.Println("private key required")
		return
	}
	privBytes, err := base64.StdEncoding.DecodeString(*priv)
	if err != nil {
		fmt.Println("invalid private key")
		return
	}
	canonical := crypto.CanonicalJSON(crypto.CanonicalPayload{
		TenantID:    payloadResp.Payload.TenantID,
		UserID:      payloadResp.Payload.UserID,
		DeviceID:    payloadResp.Payload.DeviceID,
		ChallengeID: payloadResp.Payload.ChallengeID,
		Nonce:       payloadResp.Payload.Nonce,
		Action:      payloadResp.Payload.Action,
		IssuedAt:    mustParse(payloadResp.Payload.IssuedAt),
		ExpiresAt:   mustParse(payloadResp.Payload.ExpiresAt),
		ContextHash: payloadResp.Payload.ContextHash,
	})
	sig := ed25519.Sign(ed25519.PrivateKey(privBytes), canonical)
	body := map[string]any{
		"device_id": *device,
		"signature": base64.StdEncoding.EncodeToString(sig),
		"payload":   payloadResp.Payload,
	}
	b, _ := json.Marshal(body)
	resp2 := doReq("POST", *base+"/v1/auth/challenges/"+*id+"/"+action, b, nil)
	defer resp2.Body.Close()
	io.Copy(os.Stdout, resp2.Body)
}

func doReq(method, url string, body []byte, headers map[string]string) *http.Response {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, _ := http.NewRequest(method, url, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("request failed:", err)
		os.Exit(1)
	}
	return resp
}

func mustParse(s string) time.Time {
	t, _ := time.Parse(time.RFC3339Nano, strings.TrimSpace(s))
	return t
}
