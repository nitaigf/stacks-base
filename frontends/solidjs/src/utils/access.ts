import type { User } from '../types/auth';

export function isAuthenticated(user: User | null) {
  return Boolean(user);
}

export function canAccessAdmin(user: User | null) {
  return user?.role === 'admin';
}

export function canAccessPrivate(user: User | null) {
  return Boolean(user);
}
