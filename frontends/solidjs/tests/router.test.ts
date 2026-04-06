import { describe, expect, it } from 'vitest';
import { resolveRoute } from '../src/utils/router';

const member = {
  id: 'user-1',
  name: 'Member',
  email: 'member@example.com',
  role: 'member' as const,
  status: 'active' as const,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
};

const admin = {
  ...member,
  role: 'admin' as const,
};

describe('route resolution', () => {
  it('maps unknown routes to 404', () => {
    expect(resolveRoute('/missing', null)).toEqual({ kind: 'error', path: '/errors/404', statusCode: 404 });
  });

  it('redirects anonymous users from private routes to login', () => {
    expect(resolveRoute('/app', null)).toEqual({ kind: 'auth', path: '/auth/login', mode: 'login' });
  });

  it('blocks non-admin users from admin route', () => {
    expect(resolveRoute('/admin', member)).toEqual({ kind: 'error', path: '/errors/403', statusCode: 403 });
  });

  it('allows admins into admin route', () => {
    expect(resolveRoute('/admin', admin)).toEqual({ kind: 'admin', path: '/admin' });
  });
});