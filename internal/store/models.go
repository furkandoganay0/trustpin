package store

import "time"

type EnrollmentState string

type DeviceState string

type ChallengeState string

const (
	EnrollmentNotEnrolled  EnrollmentState = "NOT_ENROLLED"
	EnrollmentPairing      EnrollmentState = "PAIRING_PENDING"
	EnrollmentDeviceActive EnrollmentState = "DEVICE_ACTIVE"
	EnrollmentDeviceRev    EnrollmentState = "DEVICE_REVOKED"
)

const (
	DevicePending  DeviceState = "PENDING"
	DeviceActive   DeviceState = "ACTIVE"
	DeviceSuspended DeviceState = "SUSPENDED"
	DeviceRevoked  DeviceState = "REVOKED"
)

const (
	ChallengeCreated  ChallengeState = "CREATED"
	ChallengePushSent ChallengeState = "PUSH_SENT"
	ChallengeApproved ChallengeState = "APPROVED"
	ChallengeRejected ChallengeState = "REJECTED"
	ChallengeFailed   ChallengeState = "FAILED"
	ChallengeExpired  ChallengeState = "EXPIRED"
)

type Tenant struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type User struct {
	ID        string
	TenantID  string
	ExternalRef string
	CreatedAt time.Time
}

type Device struct {
	ID         string
	TenantID   string
	UserID     string
	PublicKey  []byte
	State      DeviceState
	Label      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Enrollment struct {
	ID           string
	TenantID     string
	UserID       string
	State        EnrollmentState
	PairingHash  []byte
	DeviceID     *string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Challenge struct {
	ID         string
	TenantID   string
	UserID     string
	DeviceID   string
	Nonce      string
	Action     string
	Context    []byte
	ContextHash string
	State      ChallengeState
	IssuedAt   time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Policy struct {
	ID           string
	TenantID     string
	Name         string
	TOTPRequired bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuditEvent struct {
	ID          string
	TenantID    string
	UserID      *string
	DeviceID    *string
	ChallengeID *string
	ActorType   string
	ActorID     string
	EventType   string
	Metadata    []byte
	CreatedAt   time.Time
}
