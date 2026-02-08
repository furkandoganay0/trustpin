package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	DB *sql.DB
}

func Open(ctx context.Context, url string) (*Postgres, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}

func (p *Postgres) EnsureTenant(ctx context.Context, tenantID string) error {
	_, err := p.DB.ExecContext(ctx, `
		insert into tenants (id, name) values ($1, $1)
		on conflict (id) do nothing
	`, tenantID)
	return err
}

func (p *Postgres) EnsureUser(ctx context.Context, tenantID, userID string) error {
	_, err := p.DB.ExecContext(ctx, `
		insert into users (id, tenant_id, external_ref) values ($1, $2, $1)
		on conflict (id, tenant_id) do nothing
	`, userID, tenantID)
	return err
}

func (p *Postgres) CreateEnrollment(ctx context.Context, tenantID, userID string, pairingHash []byte, expiresAt time.Time) (Enrollment, error) {
	id := uuid.NewString()
	row := p.DB.QueryRowContext(ctx, `
		insert into enrollments (id, tenant_id, user_id, state, pairing_code_hash, expires_at)
		values ($1, $2, $3, $4, $5, $6)
		returning id, tenant_id, user_id, state, pairing_code_hash, device_id, expires_at, created_at, updated_at
	`, id, tenantID, userID, EnrollmentPairing, pairingHash, expiresAt)
	var e Enrollment
	var deviceID sql.NullString
	if err := row.Scan(&e.ID, &e.TenantID, &e.UserID, &e.State, &e.PairingHash, &deviceID, &e.ExpiresAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return Enrollment{}, err
	}
	if deviceID.Valid {
		e.DeviceID = &deviceID.String
	}
	return e, nil
}

func (p *Postgres) GetEnrollmentByID(ctx context.Context, id string) (Enrollment, error) {
	row := p.DB.QueryRowContext(ctx, `
		select id, tenant_id, user_id, state, pairing_code_hash, device_id, expires_at, created_at, updated_at
		from enrollments where id = $1
	`, id)
	var e Enrollment
	var deviceID sql.NullString
	if err := row.Scan(&e.ID, &e.TenantID, &e.UserID, &e.State, &e.PairingHash, &deviceID, &e.ExpiresAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return Enrollment{}, err
	}
	if deviceID.Valid {
		e.DeviceID = &deviceID.String
	}
	return e, nil
}

func (p *Postgres) GetEnrollmentByPairingHash(ctx context.Context, hash []byte) (Enrollment, error) {
	row := p.DB.QueryRowContext(ctx, `
		select id, tenant_id, user_id, state, pairing_code_hash, device_id, expires_at, created_at, updated_at
		from enrollments where pairing_code_hash = $1
	`, hash)
	var e Enrollment
	var deviceID sql.NullString
	if err := row.Scan(&e.ID, &e.TenantID, &e.UserID, &e.State, &e.PairingHash, &deviceID, &e.ExpiresAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return Enrollment{}, err
	}
	if deviceID.Valid {
		e.DeviceID = &deviceID.String
	}
	return e, nil
}

func (p *Postgres) UpdateEnrollmentState(ctx context.Context, id string, state EnrollmentState, deviceID *string) error {
	_, err := p.DB.ExecContext(ctx, `
		update enrollments set state = $2, device_id = $3, updated_at = now() where id = $1
	`, id, state, deviceID)
	return err
}

func (p *Postgres) UpdateEnrollmentStateByDeviceID(ctx context.Context, deviceID string, state EnrollmentState) error {
	_, err := p.DB.ExecContext(ctx, `
		update enrollments set state = $2, updated_at = now() where device_id = $1
	`, deviceID, state)
	return err
}

func (p *Postgres) CreateDevice(ctx context.Context, tenantID, userID string, publicKey []byte, label string, state DeviceState) (Device, error) {
	id := uuid.NewString()
	row := p.DB.QueryRowContext(ctx, `
		insert into devices (id, tenant_id, user_id, public_key, state, label)
		values ($1, $2, $3, $4, $5, $6)
		returning id, tenant_id, user_id, public_key, state, label, created_at, updated_at
	`, id, tenantID, userID, publicKey, state, label)
	var d Device
	if err := row.Scan(&d.ID, &d.TenantID, &d.UserID, &d.PublicKey, &d.State, &d.Label, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return Device{}, err
	}
	return d, nil
}

func (p *Postgres) GetDevice(ctx context.Context, id string) (Device, error) {
	row := p.DB.QueryRowContext(ctx, `
		select id, tenant_id, user_id, public_key, state, label, created_at, updated_at
		from devices where id = $1
	`, id)
	var d Device
	if err := row.Scan(&d.ID, &d.TenantID, &d.UserID, &d.PublicKey, &d.State, &d.Label, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return Device{}, err
	}
	return d, nil
}

func (p *Postgres) UpdateDeviceState(ctx context.Context, id string, state DeviceState) error {
	_, err := p.DB.ExecContext(ctx, `
		update devices set state = $2, updated_at = now() where id = $1
	`, id, state)
	return err
}

func (p *Postgres) CreateChallenge(ctx context.Context, c Challenge) (Challenge, error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	row := p.DB.QueryRowContext(ctx, `
		insert into challenges (id, tenant_id, user_id, device_id, nonce, action, context, context_hash, state, issued_at, expires_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		returning id, tenant_id, user_id, device_id, nonce, action, context, context_hash, state, issued_at, expires_at, created_at, updated_at
	`, c.ID, c.TenantID, c.UserID, c.DeviceID, c.Nonce, c.Action, c.Context, c.ContextHash, c.State, c.IssuedAt, c.ExpiresAt)
	var out Challenge
	if err := row.Scan(&out.ID, &out.TenantID, &out.UserID, &out.DeviceID, &out.Nonce, &out.Action, &out.Context, &out.ContextHash, &out.State, &out.IssuedAt, &out.ExpiresAt, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return Challenge{}, err
	}
	return out, nil
}

func (p *Postgres) GetChallenge(ctx context.Context, id string) (Challenge, error) {
	row := p.DB.QueryRowContext(ctx, `
		select id, tenant_id, user_id, device_id, nonce, action, context, context_hash, state, issued_at, expires_at, created_at, updated_at
		from challenges where id = $1
	`, id)
	var c Challenge
	if err := row.Scan(&c.ID, &c.TenantID, &c.UserID, &c.DeviceID, &c.Nonce, &c.Action, &c.Context, &c.ContextHash, &c.State, &c.IssuedAt, &c.ExpiresAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return Challenge{}, err
	}
	return c, nil
}

func (p *Postgres) UpdateChallengeState(ctx context.Context, id string, state ChallengeState) error {
	res, err := p.DB.ExecContext(ctx, `
		update challenges set state = $2, updated_at = now() where id = $1
	`, id, state)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("challenge not found")
	}
	return nil
}

func (p *Postgres) InsertAuditEvent(ctx context.Context, e AuditEvent) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	_, err := p.DB.ExecContext(ctx, `
		insert into audit_events (id, tenant_id, user_id, device_id, challenge_id, actor_type, actor_id, event_type, metadata)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, e.ID, e.TenantID, e.UserID, e.DeviceID, e.ChallengeID, e.ActorType, e.ActorID, e.EventType, e.Metadata)
	return err
}

func (p *Postgres) EnsurePolicy(ctx context.Context, tenantID, name string) error {
	_, err := p.DB.ExecContext(ctx, `
		insert into policies (id, tenant_id, name, totp_required)
		values ($1, $2, $3, false)
		on conflict (id) do nothing
	`, tenantID+":"+name, tenantID, name)
	return err
}

func (p *Postgres) GetPolicy(ctx context.Context, tenantID, name string) (Policy, error) {
	row := p.DB.QueryRowContext(ctx, `
		select id, tenant_id, name, totp_required, created_at, updated_at
		from policies where id = $1
	`, tenantID+":"+name)
	var pol Policy
	if err := row.Scan(&pol.ID, &pol.TenantID, &pol.Name, &pol.TOTPRequired, &pol.CreatedAt, &pol.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Policy{TenantID: tenantID, Name: name, TOTPRequired: false}, nil
		}
		return Policy{}, err
	}
	return pol, nil
}

func (p *Postgres) ListAuditEvents(ctx context.Context, tenantID string, limit int) ([]AuditEvent, error) {
	rows, err := p.DB.QueryContext(ctx, `
		select id, tenant_id, user_id, device_id, challenge_id, actor_type, actor_id, event_type, metadata, created_at
		from audit_events where tenant_id = $1 order by created_at desc limit $2
	`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AuditEvent
	for rows.Next() {
		var e AuditEvent
		var userID, deviceID, challengeID sql.NullString
		if err := rows.Scan(&e.ID, &e.TenantID, &userID, &deviceID, &challengeID, &e.ActorType, &e.ActorID, &e.EventType, &e.Metadata, &e.CreatedAt); err != nil {
			return nil, err
		}
		if userID.Valid {
			e.UserID = &userID.String
		}
		if deviceID.Valid {
			e.DeviceID = &deviceID.String
		}
		if challengeID.Valid {
			e.ChallengeID = &challengeID.String
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
