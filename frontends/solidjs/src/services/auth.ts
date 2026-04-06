import { api } from './api';
import type { AuthEnvelope, UserEnvelope } from '../types/auth';
import type { LoginInput, RegisterInput } from '../schemas/auth';

export async function register(input: RegisterInput) {
  return api.post('api/v1/auth/register', { json: input }).json<AuthEnvelope>();
}

export async function login(input: LoginInput) {
  return api.post('api/v1/auth/login', { json: input }).json<AuthEnvelope>();
}

export async function me() {
  return api.get('api/v1/users/me').json<UserEnvelope>();
}

export async function logout() {
  await api.post('api/v1/auth/logout');
}