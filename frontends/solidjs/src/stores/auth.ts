import { createSignal } from 'solid-js';
import type { User } from '../types/auth';

const [accessToken, setAccessToken] = createSignal<string | null>(null);
const [currentUser, setCurrentUser] = createSignal<User | null>(null);

export const authStore = {
  accessToken,
  currentUser,
  setSession(token: string, user: User) {
    setAccessToken(token);
    setCurrentUser(user);
  },
  clearSession() {
    setAccessToken(null);
    setCurrentUser(null);
  },
};