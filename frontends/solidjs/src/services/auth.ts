import { api } from './api';
import type { AuthEnvelope, UserEnvelope } from '../types/auth';
import type { LoginInput, RegisterInput } from '../schemas/auth';

// REF.AUTH-01|Register
export async function register(input: RegisterInput) {
  return api.post('api/v1/auth/register', { json: input }).json<AuthEnvelope>();
}

// REF.AUTH-02|Login
export async function login(input: LoginInput) {
  return api.post('api/v1/auth/login', { json: input }).json<AuthEnvelope>();
}

// REF.AUTH-04|Me
export async function me() {
  return api.get('api/v1/users/me').json<UserEnvelope>();
}

// REF.AUTH-03|Logout
export async function logout() {
  await api.post('api/v1/auth/logout');
}