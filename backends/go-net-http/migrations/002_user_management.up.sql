alter table users
  add column if not exists deleted_at timestamptz,
  add column if not exists deleted_by uuid references users(id) on delete set null,
  add column if not exists last_login_at timestamptz;

alter table audit_logs
  add column if not exists actor_name text,
  add column if not exists actor_email text,
  add column if not exists route text not null default 'legacy',
  add column if not exists method text not null default 'SYSTEM',
  add column if not exists ip_address text,
  add column if not exists user_agent text;

update users
set status = 'inactive'
where status = 'blocked';

alter table users drop constraint if exists users_email_key;
alter table users drop constraint if exists users_status_check;
alter table users
  add constraint users_status_check check (status in ('active', 'inactive'));

create unique index if not exists idx_users_email_unique on users (lower(email)) where deleted_at is null;
create index if not exists idx_users_status on users(status);
create index if not exists idx_users_role on users(role);
create index if not exists idx_users_deleted_at on users(deleted_at);

create table if not exists password_reset_tokens (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references users(id) on delete cascade,
  token_hash text not null unique,
  expires_at timestamptz not null,
  used_at timestamptz,
  created_at timestamptz not null default now()
);

create index if not exists idx_password_reset_tokens_user_id on password_reset_tokens(user_id);
create index if not exists idx_password_reset_tokens_expires_at on password_reset_tokens(expires_at);

create index if not exists idx_audit_logs_action on audit_logs(action);
create index if not exists idx_audit_logs_route on audit_logs(route);
