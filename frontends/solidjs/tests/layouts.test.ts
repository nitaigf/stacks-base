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

const admin = { ...member, role: 'admin' as const };

describe('layout assignment', () => {
  it('public route uses public layout', () => {
    const route = resolveRoute('/', null);
    expect(route.kind).toBe('public');
  });

  it('auth routes use auth layout', () => {
    expect(resolveRoute('/auth/login', null).kind).toBe('auth');
    expect(resolveRoute('/auth/register', null).kind).toBe('auth');
  });

  it('private route uses private layout for logged-in users', () => {
    const route = resolveRoute('/app', member);
    expect(route.kind).toBe('private');
  });

  it('admin route uses admin layout for admin users', () => {
    const route = resolveRoute('/admin', admin);
    expect(route.kind).toBe('admin');
  });

  it('error routes use error layout', () => {
    expect(resolveRoute('/errors/403', null).kind).toBe('error');
    expect(resolveRoute('/errors/500', null).kind).toBe('error');
    expect(resolveRoute('/unknown', null).kind).toBe('error');
  });
});
