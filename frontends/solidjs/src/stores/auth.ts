import { createSignal } from 'solid-js';
import type { User } from '../types/auth';

export const AUTH_STORAGE_KEY = 'stacks-base.auth-session';

type PersistedSession = {
  accessToken: string;
  user: User;
};

const persistedSession = readPersistedSession();
const [accessToken, setAccessToken] = createSignal<string | null>(persistedSession?.accessToken ?? null);
const [currentUser, setCurrentUser] = createSignal<User | null>(persistedSession?.user ?? null);
let sessionValidationPromise: Promise<User | null> | null = null;

export const authStore = {
  accessToken,
  currentUser,
  hasSession() {
    return Boolean(accessToken() && currentUser());
  },
  setSession(token: string, user: User) {
    setAccessToken(token);
    setCurrentUser(user);
    persistSession(token, user);
  },
  setCurrentUser(user: User) {
    const token = accessToken();
    if (!token) {
      authStore.clearSession();
      return;
    }

    setCurrentUser(user);
    persistSession(token, user);
  },
  clearSession() {
    sessionValidationPromise = null;
    setAccessToken(null);
    setCurrentUser(null);
    clearPersistedSession();
  },
  async revalidateSession(loadUser: () => Promise<User>) {
    if (!authStore.hasSession()) {
      authStore.clearSession();
      return null;
    }

    if (!sessionValidationPromise) {
      sessionValidationPromise = (async () => {
        try {
          const user = await loadUser();
          authStore.setCurrentUser(user);
          return user;
        } finally {
          sessionValidationPromise = null;
        }
      })();
    }

    return sessionValidationPromise;
  },
};

function readPersistedSession(): PersistedSession | null {
  if (typeof window === 'undefined') {
    return null;
  }

  const raw = window.localStorage.getItem(AUTH_STORAGE_KEY);
  if (!raw) {
    return null;
  }

  try {
    const parsed = JSON.parse(raw) as Partial<PersistedSession>;
    if (!parsed || typeof parsed.accessToken !== 'string' || !isUser(parsed.user)) {
      clearPersistedSession();
      return null;
    }

    return {
      accessToken: parsed.accessToken,
      user: parsed.user,
    };
  } catch {
    clearPersistedSession();
    return null;
  }
}

function persistSession(token: string, user: User) {
  if (typeof window === 'undefined') {
    return;
  }

  window.localStorage.setItem(
    AUTH_STORAGE_KEY,
    JSON.stringify({
      accessToken: token,
      user,
    } satisfies PersistedSession),
  );
}

function clearPersistedSession() {
  if (typeof window === 'undefined') {
    return;
  }

  window.localStorage.removeItem(AUTH_STORAGE_KEY);
}

function isUser(value: unknown): value is User {
  if (!value || typeof value !== 'object') {
    return false;
  }

  const candidate = value as Partial<User>;
  return (
    typeof candidate.id === 'string' &&
    typeof candidate.name === 'string' &&
    typeof candidate.email === 'string' &&
    (candidate.role === 'admin' || candidate.role === 'member') &&
    (candidate.status === 'active' || candidate.status === 'inactive') &&
    typeof candidate.createdAt === 'string' &&
    typeof candidate.updatedAt === 'string'
  );
}
