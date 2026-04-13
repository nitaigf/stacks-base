import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { User } from '../src/types/auth';

const storedUser: User = {
  id: 'user-1',
  name: 'Admin',
  email: 'admin@stacks-base.local',
  role: 'admin',
  status: 'active',
  createdAt: '2026-04-08T00:00:00.000Z',
  updatedAt: '2026-04-08T00:00:00.000Z',
};

describe('authStore', () => {
  beforeEach(() => {
    vi.resetModules();
    window.localStorage.clear();
  });

  it('rehydrates the session from localStorage', async () => {
    window.localStorage.setItem(
      'stacks-base.auth-session',
      JSON.stringify({
        accessToken: 'persisted-token',
        user: storedUser,
      }),
    );

    const { authStore } = await import('../src/stores/auth');

    expect(authStore.accessToken()).toBe('persisted-token');
    expect(authStore.currentUser()).toEqual(storedUser);
  });

  it('persists the session after login', async () => {
    const { AUTH_STORAGE_KEY, authStore } = await import('../src/stores/auth');

    authStore.setSession('fresh-token', storedUser);

    expect(window.localStorage.getItem(AUTH_STORAGE_KEY)).toBe(
      JSON.stringify({
        accessToken: 'fresh-token',
        user: storedUser,
      }),
    );
  });

  it('clears the in-memory and persisted session on logout', async () => {
    const { AUTH_STORAGE_KEY, authStore } = await import('../src/stores/auth');

    authStore.setSession('fresh-token', storedUser);
    authStore.clearSession();

    expect(authStore.accessToken()).toBeNull();
    expect(authStore.currentUser()).toBeNull();
    expect(window.localStorage.getItem(AUTH_STORAGE_KEY)).toBeNull();
  });
});
