import { describe, expect, it } from 'vitest';
import { canAccessAdmin, canAccessPrivate, isAuthenticated } from '../src/utils/access';

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

describe('route access helpers', () => {
  it('knows when a user is authenticated', () => {
    expect(isAuthenticated(null)).toBe(false);
    expect(isAuthenticated(member)).toBe(true);
  });

  it('allows private access only for authenticated users', () => {
    expect(canAccessPrivate(null)).toBe(false);
    expect(canAccessPrivate(member)).toBe(true);
  });

  it('allows admin access only for admins', () => {
    expect(canAccessAdmin(null)).toBe(false);
    expect(canAccessAdmin(member)).toBe(false);
    expect(canAccessAdmin(admin)).toBe(true);
  });
});
