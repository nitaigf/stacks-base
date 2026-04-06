import { api } from './api';
import type { UserEnvelope } from '../types/auth';

// REF.ADMIN-01|ListUsers
export async function listUsers() {
  return api.get('api/v1/users').json<UserEnvelope>();
}

// REF.ADMIN-02|UpdateUserStatus
export async function updateUserStatus(userId: string, status: 'active' | 'blocked') {
  return api.patch(`api/v1/users/${userId}/status`, { 
    json: { status } 
  }).json<UserEnvelope>();
}

// REF.ADMIN-03|UpdateUserRole
export async function updateUserRole(userId: string, role: 'admin' | 'member') {
  return api.patch(`api/v1/users/${userId}/role`, { 
    json: { role } 
  }).json<UserEnvelope>();
}
