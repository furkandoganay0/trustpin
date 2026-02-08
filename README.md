# TRUSTPIN

Production-grade, on-prem MFA backend with asymmetric cryptography, strict state machines, and append-only audit logging.

## Architecture
- API service (Go) with stateless HTTP handlers
- PostgreSQL for source-of-truth state and audit log
- Redis for TTL caches (enrollment sessions, challenge cache, nonce uniqueness, rate limiting)
- Push abstraction with a mock provider

## State Machines

Enrollment:
- NOT_ENROLLED → PAIRING_PENDING → DEVICE_ACTIVE → DEVICE_REVOKED

Device Lifecycle:
- PENDING → ACTIVE → SUSPENDED → REVOKED

Challenge Lifecycle:
- CREATED → PUSH_SENT → APPROVED / REJECTED / FAILED / EXPIRED

Invalid transitions return HTTP 409.

## Security Model
- Device generates Ed25519 keypair; private key never leaves device.
- Server stores only public key.
- Challenges include nonce + TTL; nonce uniqueness enforced with Redis SETNX.
- Canonical payload is deterministic JSON with ordered fields.
- Server compromise cannot approve without device signature.
- All actions emit audit events (append-only).

## Canonical Signing Payload
Fields and order:

```
{
  tenant_id,
  user_id,
  device_id,
  challenge_id,
  nonce,
  action,
  issued_at,
  expires_at,
  context_hash
}
```

- `issued_at` / `expires_at` are RFC3339Nano UTC.
- `context_hash` is SHA-256 hex of normalized JSON context.
- Signature is Ed25519 over the canonical JSON string.

## Running Locally

```
docker compose up --build
```

API listens on `http://localhost:8080`.
Swagger UI is available at `http://localhost:8080/swagger/`.

## Example Approval Flow (cURL)

1) Init enrollment
```
curl -s -X POST http://localhost:8080/v1/enrollments/init \
  -H Content-Type:
