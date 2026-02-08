create table if not exists tenants (
  id text primary key,
  name text not null,
  created_at timestamptz not null default now()
);

create table if not exists users (
  id text not null,
  tenant_id text not null references tenants(id),
  external_ref text not null,
  created_at timestamptz not null default now(),
  primary key (id, tenant_id)
);

create table if not exists devices (
  id text primary key,
  tenant_id text not null references tenants(id),
  user_id text not null,
  public_key bytea not null,
  state text not null,
  label text not null default '',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
create index if not exists idx_devices_tenant_user on devices (tenant_id, user_id);

create table if not exists enrollments (
  id text primary key,
  tenant_id text not null references tenants(id),
  user_id text not null,
  state text not null,
  pairing_code_hash bytea not null,
  device_id text null references devices(id),
  expires_at timestamptz not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
create index if not exists idx_enrollments_pairing_hash on enrollments (pairing_code_hash);
create index if not exists idx_enrollments_tenant_user on enrollments (tenant_id, user_id);

create table if not exists challenges (
  id text primary key,
  tenant_id text not null references tenants(id),
  user_id text not null,
  device_id text not null references devices(id),
  nonce text not null,
  action text not null,
  context jsonb not null,
  context_hash text not null,
  state text not null,
  issued_at timestamptz not null,
  expires_at timestamptz not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
create index if not exists idx_challenges_tenant_user on challenges (tenant_id, user_id);
create index if not exists idx_challenges_device on challenges (device_id);
create index if not exists idx_challenges_expires on challenges (expires_at);

create table if not exists policies (
  id text primary key,
  tenant_id text not null references tenants(id),
  name text not null,
  totp_required boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists audit_events (
  id text primary key,
  tenant_id text not null references tenants(id),
  user_id text null,
  device_id text null,
  challenge_id text null,
  actor_type text not null,
  actor_id text not null,
  event_type text not null,
  metadata jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);
create index if not exists idx_audit_events_tenant_created on audit_events (tenant_id, created_at desc);
