import { api } from './api';
import type {
  UserEnvelope,
  UserListFilters,
  UserRole,
  UserStatus,
  UsersEnvelope,
} from '../types/auth';

type UserPayload = {
  name: string;
  email: string;
  role: UserRole;
};

type CreateUserPayload = UserPayload & {
  password: string;
  status: UserStatus;
};

// REF.ADMIN-01|ListUsers
export async function listUsers(filters: UserListFilters = {}) {
  return api.get('api/v1/users', { searchParams: buildUserFilters(filters) }).json<UsersEnvelope>();
}

// REF.ADMIN-02|ShowUser
export async function getUser(userId: string) {
  return api.get(`api/v1/users/${userId}`).json<UserEnvelope>();
}

// REF.ADMIN-03|CreateUser
export async function createUser(payload: CreateUserPayload) {
  return api.post('api/v1/users', { json: payload }).json<UserEnvelope>();
}

// REF.ADMIN-04|UpdateUser
export async function updateUser(userId: string, payload: UserPayload) {
  return api.patch(`api/v1/users/${userId}`, { json: payload }).json<UserEnvelope>();
}

// REF.ADMIN-05|DeactivateUser
export async function deactivateUser(userId: string) {
  return api.post(`api/v1/users/${userId}/deactivate`).json<UserEnvelope>();
}

// REF.ADMIN-06|ReactivateUser
export async function reactivateUser(userId: string) {
  return api.post(`api/v1/users/${userId}/reactivate`).json<UserEnvelope>();
}

// REF.ADMIN-07|SoftDeleteUser
export async function softDeleteUser(userId: string) {
  return api.post(`api/v1/users/${userId}/soft-delete`).json<UserEnvelope>();
}

// REF.ADMIN-08|RestoreUser
export async function restoreUser(userId: string) {
  return api.post(`api/v1/users/${userId}/restore`).json<UserEnvelope>();
}

// REF.ADMIN-09|HardDeleteUser
export async function hardDeleteUser(userId: string) {
  await api.delete(`api/v1/users/${userId}`);
}

export async function exportUsersCsv(filters: UserListFilters = {}) {
  return api.get('api/v1/users/export.csv', { searchParams: buildUserFilters(filters) }).blob();
}

export async function exportUsersXlsx(filters: UserListFilters = {}) {
  return api.get('api/v1/users/export.xlsx', { searchParams: buildUserFilters(filters) }).blob();
}

export async function printUsers(filters: UserListFilters = {}) {
  return api.get('api/v1/users/print', { searchParams: buildUserFilters(filters) }).blob();
}

function buildUserFilters(filters: UserListFilters) {
  const searchParams = new URLSearchParams();
  if (filters.page) {
    searchParams.set('page', String(filters.page));
  }
  if (filters.perPage) {
    searchParams.set('perPage', String(filters.perPage));
  }
  if (filters.query) {
    searchParams.set('query', filters.query);
  }
  if (filters.role) {
    searchParams.set('role', filters.role);
  }
  if (filters.status) {
    searchParams.set('status', filters.status);
  }
  if (filters.includeDeleted) {
    searchParams.set('includeDeleted', 'true');
  }
  return searchParams;
}
